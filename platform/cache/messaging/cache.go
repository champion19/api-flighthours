package messaging

import (
	"context"
	"net/http"
	"sync"
	"time"

	cachetypes "github.com/champion19/flighthours-api/platform/cache/types"
	"github.com/champion19/flighthours-api/platform/logger"
)

type MessageType = cachetypes.MessageType
type CachedMessage = cachetypes.CachedMessage
type MessageResponse = cachetypes.MessageResponse
type MessageCacheRepository = cachetypes.MessageCacheRepository

const (
	TypeError   = cachetypes.TypeError
	TypeSuccess = cachetypes.TypeSuccess
	TypeWarning = cachetypes.TypeWarning
	TypeInfo    = cachetypes.TypeInfo
	TypeDebug   = cachetypes.TypeDebug
)

type MessageCache struct {
	repo            MessageCacheRepository
	messages        map[string]*CachedMessage
	mu              sync.RWMutex
	refreshInterval time.Duration
	stopRefresh     chan bool
}

func NewMessageCache(repo MessageCacheRepository, refreshInterval time.Duration) *MessageCache {
	return &MessageCache{
		repo:            repo,
		messages:        make(map[string]*CachedMessage),
		refreshInterval: refreshInterval,
		stopRefresh:     make(chan bool),
	}
}

var log logger.Logger = logger.NewSlogLogger()

func (c *MessageCache) LoadMessages(ctx context.Context) error {
	messages, err := c.repo.GetAllActiveForCache(ctx)
	if err != nil {
		log.Error(logger.LogMsgCacheLoadError, "error", err.Error())
		return err
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	c.messages = make(map[string]*CachedMessage)
	for i := range messages {
		c.messages[messages[i].Code] = &messages[i]
	}
	log.Info(logger.LogMsgCacheLoaded, "count", len(messages))
	return nil
}

func (c *MessageCache) ReloadMessages(ctx context.Context) error {
	return c.LoadMessages(ctx)
}

func (c *MessageCache) StartAutoRefresh(ctx context.Context) {
	if c.refreshInterval <= 0 {
		log.Info(logger.LogMsgCacheRefreshDisabled)
		return
	}

	log.Info(logger.LogMsgCacheRefreshStart, "interval", c.refreshInterval.String())

	go func() {
		ticker := time.NewTicker(c.refreshInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				log.Debug(logger.LogMsgCacheRefreshing)
				if err := c.ReloadMessages(ctx); err != nil {
					log.Error(logger.LogMsgCacheRefreshError, "error", err.Error())
				} else {
					log.Debug(logger.LogMsgCacheRefreshOK, "count", c.MessageCount())
				}
			case <-c.stopRefresh:
				log.Info(logger.LogMsgCacheRefreshStop)
				return
			}
		}
	}()
}

func (c *MessageCache) StopAutoRefresh() {
	if c.refreshInterval > 0 {
		close(c.stopRefresh)
	}
}

func (c *MessageCache) GetMessage(code string) *CachedMessage {

	c.mu.RLock()
	msg, found := c.messages[code]
	c.mu.RUnlock()

	if found {
		return msg
	}

	log.Debug(logger.LogMsgNotInCache, "code", code)
	dbMsg, err := c.repo.GetByCodeForCache(context.Background(), code)
	if err != nil {
		log.Warn(logger.LogMsgNotInDB, "code", code, "error", err)

		if code == "GEN_MSG_INACTIVE_ERR_00002" {
			return nil
		}
		return c.GetMessage("GEN_MSG_INACTIVE_ERR_00002")
	}

	if dbMsg != nil {

		c.mu.Lock()
		c.messages[code] = dbMsg
		c.mu.Unlock()

		log.Debug(logger.LogMsgCachedFromDB, "code", code)
		return dbMsg
	}

	inactiveMsg, err := c.repo.GetByCodeWithStatusForCache(context.Background(), code)
	if err != nil {
		log.Warn(logger.LogMsgNotInDB, "code", code, "error", err)
		if code == "GEN_MSG_INACTIVE_ERR_00002" {
			return nil
		}
		return c.GetMessage("GEN_MSG_INACTIVE_ERR_00002")
	}

	if inactiveMsg != nil && !inactiveMsg.Active {
		log.Warn(logger.LogMsgInactive, "code", code)
		return c.GetMessage("GEN_MSG_INACTIVE_ERR_00002")
	}
	log.Warn(logger.LogMsgNotInDB, "code", code)
	if code == "GEN_MSG_INACTIVE_ERR_00002" {
		return nil
	}
	return c.GetMessage("GEN_MSG_INACTIVE_ERR_00002")
}

func (c *MessageCache) GetMessageResponse(code string, params ...string) *MessageResponse {
	msg := c.GetMessage(code)
	if msg == nil {
		return nil
	}

	content := msg.Content
	for i, param := range params {
		placeholder := "${" + string(rune('0'+i)) + "}"
		content = replaceAll(content, placeholder, param)
	}

	return &MessageResponse{
		Code:    msg.Code,
		Type:    msg.Type,
		Title:   msg.Title,
		Content: content,
	}
}

func replaceAll(s, old, new string) string {
	result := ""
	for i := 0; i < len(s); {
		if i+len(old) <= len(s) && s[i:i+len(old)] == old {
			result += new
			i += len(old)
		} else {
			result += string(s[i])
			i++
		}
	}
	return result
}

var messageCodeToHTTPStatus = map[string]int{
	"MOD_U_REG_EXI_00001":        http.StatusCreated,
	"MOD_U_UPD_EXI_00002":        http.StatusOK,
	"MOD_U_GET_EXI_00005":        http.StatusOK,
	"MOD_U_DUP_ERR_00001":        http.StatusConflict,
	"MOD_U_DUP_IDNUM_ERR_00013":  http.StatusConflict,
	"MOD_U_EMAIL_NF_ERR_00005":   http.StatusNotFound,
	"MOD_U_GET_ERR_00003":        http.StatusNotFound,
	"MOD_U_TOKEN_NF_ERR_00007":   http.StatusNotFound,
	"MOD_U_EMAIL_NV_ERR_00006":   http.StatusForbidden,
	"MOD_U_TOKEN_EXP_ERR_00008":  http.StatusUnauthorized,
	"MOD_U_TOKEN_USED_ERR_00009": http.StatusUnauthorized,
	"MOD_U_UPD_ERR_00013":      http.StatusInternalServerError,
	"MOD_U_KC_UPD_ERR_00014":   http.StatusServiceUnavailable,
	"MOD_U_ROLE_UPD_ERR_00015": http.StatusServiceUnavailable,
	"MOD_U_DEL_EXI_00003": http.StatusOK,
	"MOD_U_DEL_ERR_00012": http.StatusInternalServerError,

	"MOD_P_NOT_FOUND_ERR_00001": http.StatusNotFound,

	"MOD_V_VAL_ERR_00001":  http.StatusBadRequest,
	"MOD_V_VAL_ERR_00002":  http.StatusBadRequest,
	"MOD_V_VAL_ERR_00006":  http.StatusBadRequest,
	"MOD_V_VAL_ERR_00008":  http.StatusBadRequest,
	"MOD_V_VAL_ERR_00009":  http.StatusBadRequest,
	"MOD_V_VAL_ERR_00010":  http.StatusBadRequest,
	"MOD_V_VAL_ERR_00011":  http.StatusBadRequest,
	"MOD_V_JSON_ERR_00012": http.StatusBadRequest,
	"MOD_V_ID_ERR_00013":   http.StatusBadRequest,
	"MOD_V_FK_ERR_00014":   http.StatusUnprocessableEntity,
	"MOD_V_LEN_ERR_00015":  http.StatusUnprocessableEntity,
	"MOD_V_DATA_ERR_00016": http.StatusUnprocessableEntity,
	"MOD_V_DATE_ERR_00017":  http.StatusBadRequest,
	"MOD_V_DATE_ERR_00018":  http.StatusBadRequest,
	"MOD_V_EMPTY_ERR_00019": http.StatusBadRequest,

	"MOD_KC_EMAIL_VERIFIED_EXI_00001":          http.StatusOK,
	"MOD_KC_INVALID_TOKEN_ERR_00001":           http.StatusBadRequest,
	"MOD_KC_EMAIL_VERIFY_ERROR_ERR_00001":      http.StatusInternalServerError,
	"MOD_KC_USER_NOT_FOUND_ERR_00001":          http.StatusNotFound,
	"MOD_KC_EMAIL_ALREADY_VERIFIED_WARN_00001": http.StatusOK,
	"MOD_KC_VERIF_EMAIL_SENT_EXI_00001":        http.StatusOK,
	"MOD_KC_VERIF_EMAIL_ERROR_ERR_00001":       http.StatusServiceUnavailable,
	"MOD_KC_VERIF_EMAIL_RESENT_EXI_00001":      http.StatusOK,
	"MOD_KC_PWD_RESET_SENT_EXI_00001":          http.StatusOK,
	"MOD_KC_PWD_RESET_ERROR_ERR_00001":         http.StatusServiceUnavailable,

	"MOD_KC_LOGIN_EMAIL_NOT_VERIFIED_ERR_00001": http.StatusUnauthorized,
	"MOD_KC_LOGIN_SUCCESS_EXI_00001":            http.StatusOK,

	"MOD_KC_PWD_UPDATED_EXI_00001":              http.StatusOK,
	"MOD_KC_PWD_UPDATE_ERROR_ERR_00001":         http.StatusInternalServerError,
	"MOD_KC_PWD_MISMATCH_ERR_00001":             http.StatusBadRequest,
	"MOD_KC_PWD_UPDATE_TOKEN_INVALID_ERR_00001": http.StatusUnauthorized,

	"MOD_KC_PWD_CHANGED_EXI_00001":         http.StatusOK,
	"MOD_KC_PWD_CHANGE_ERROR_ERR_00001":    http.StatusInternalServerError,
	"MOD_KC_PWD_CURRENT_INVALID_ERR_00001": http.StatusUnauthorized,
	"MOD_KC_PWD_CHANGE_MISMATCH_ERR_00001": http.StatusBadRequest,

	"MOD_AUTH_LOGIN_SUCCESS_EXI_00001": http.StatusOK,


	"MOD_INFRA_KC_UNAVAIL_ERR_00004":      http.StatusLocked,
	"MOD_INFRA_DB_UNAVAIL_ERR_00005":      http.StatusLocked,
	"MOD_INFRA_DEP_FAIL_ERR_00006":        http.StatusLocked,
	"MOD_INFRA_KC_CLEANUP_ERR_00003":      http.StatusLocked,
	"MOD_INFRA_KC_CREATE_ERR_00002":       http.StatusLocked,
	"MOD_INFRA_INCOMPLETE_REG_ERR_00009":  http.StatusConflict,
	"MOD_INFRA_KC_INCONSISTENT_ERR_00001": http.StatusInternalServerError,


	"GEN_SRV_ERR_00001":          http.StatusInternalServerError,
	"GEN_AUTH_ERR_00002":         http.StatusUnauthorized,
	"GEN_FORBIDDEN_ERR_00003":    http.StatusForbidden,
	"GEN_MSG_INACTIVE_ERR_00002": http.StatusServiceUnavailable,

	"MOD_M_CREATE_EXI_00001":    http.StatusCreated,
	"MOD_M_UPDATE_ERR_00010":    http.StatusBadRequest,
	"MOD_M_NOT_FOUND_ERR_00001": http.StatusNotFound,

	"MOD_AIR_GET_EXI_00001":        http.StatusOK,
	"MOD_AIR_ACTIVATE_EXI_00002":   http.StatusOK,
	"MOD_AIR_DEACTIVATE_EXI_00003": http.StatusOK,
	"MOD_AIR_NOT_FOUND_ERR_00001":  http.StatusNotFound,
	"MOD_AIR_ACTIVATE_ERR_00002":   http.StatusUnprocessableEntity,
	"MOD_AIR_DEACTIVATE_ERR_00003": http.StatusUnprocessableEntity,

	"MOD_APT_GET_EXI_00001":        http.StatusOK,
	"MOD_APT_ACTIVATE_EXI_00002":   http.StatusOK,
	"MOD_APT_DEACTIVATE_EXI_00003": http.StatusOK,
	"MOD_APT_NOT_FOUND_ERR_00001":  http.StatusNotFound,
	"MOD_APT_ACTIVATE_ERR_00002":   http.StatusUnprocessableEntity,
	"MOD_APT_DEACTIVATE_ERR_00003": http.StatusUnprocessableEntity,

	"CIU_CON_EXI_01301": http.StatusOK,
	"CIU_CON_ERR_01302": http.StatusNotFound,
	"CIU_CON_ERR_01303": http.StatusInternalServerError,

	"PAI_CON_EXI_03801": http.StatusOK,
	"PAI_CON_ERR_03802": http.StatusNotFound,
	"PAI_CON_ERR_03803": http.StatusInternalServerError,

	"TAE_CON_EXI_04601": http.StatusOK,
	"TAE_CON_ERR_04602": http.StatusNotFound,
	"TAE_CON_ERR_04603": http.StatusInternalServerError,

	"TIN_CON_EXI_04701": http.StatusOK,
	"TIN_CON_ERR_04702": http.StatusNotFound,
	"TIN_CON_ERR_04703": http.StatusInternalServerError,

	"BIT_CON_EXI_01901": http.StatusOK,
	"BIT_CON_ERR_01902": http.StatusBadRequest,
	"BIT_CON_ERR_01903": http.StatusNotFound,
	"BIT_CON_ERR_01904": http.StatusInternalServerError,

	"BIT_AGR_EXI_01801": http.StatusCreated,
	"BIT_AGR_ERR_01802": http.StatusBadRequest,
	"BIT_AGR_ERR_01803": http.StatusBadRequest,
	"BIT_AGR_ERR_01804": http.StatusInternalServerError,

	"BIT_EDI_EXI_01701": http.StatusOK,
	"BIT_EDI_ERR_01702": http.StatusBadRequest,
	"BIT_EDI_ERR_01703": http.StatusBadRequest,
	"BIT_EDI_ERR_01704": http.StatusInternalServerError,


	"BIT_DEL_EXI_01601": http.StatusOK,
	"BIT_DEL_ERR_01602": http.StatusBadRequest,
	"BIT_DEL_ERR_01603": http.StatusNotFound,
	"BIT_DEL_ERR_01604": http.StatusInternalServerError,

	"BIT_ACT_EXI_01501": http.StatusOK,
	"BIT_ACT_ERR_01502": http.StatusBadRequest,
	"BIT_ACT_ERR_01503": http.StatusConflict,
	"BIT_ACT_ERR_01504": http.StatusInternalServerError,

	"BIT_INA_EXI_01401": http.StatusOK,
	"BIT_INA_ERR_01402": http.StatusBadRequest,
	"BIT_INA_ERR_01403": http.StatusConflict,
	"BIT_INA_ERR_01404": http.StatusInternalServerError,

	"BIT_LIST_EXI_01001": http.StatusOK,
	"BIT_LIST_ERR_01002": http.StatusInternalServerError,

	"BIT_AUTH_ERR_00001": http.StatusForbidden,



	"MAT_CON_EXI_03301": http.StatusOK,
	"MAT_CON_ERR_03302": http.StatusNotFound,
	"MAT_CON_ERR_03303": http.StatusInternalServerError,

	"MAT_AGR_EXI_03401": http.StatusCreated,
	"MAT_AGR_ERR_03402": http.StatusBadRequest,
	"MAT_AGR_ERR_03403": http.StatusConflict,

	"MAT_EDI_EXI_03501": http.StatusOK,
	"MAT_EDI_ERR_03502": http.StatusBadRequest,

	"MAT_LIST_EXI_03001": http.StatusOK,
	"MAT_LIST_ERR_03002": http.StatusInternalServerError,

	"MAT_VAL_ERR_03601": http.StatusBadRequest,
	"MAT_VAL_ERR_03602": http.StatusBadRequest,

	"MOD_AM_CON_EXI_03601": http.StatusOK,
	"MOD_AM_CON_ERR_03602": http.StatusNotFound,
	"MOD_AM_CON_ERR_03603": http.StatusInternalServerError,

	"MOD_AM_LIST_EXI_04301": http.StatusOK,
	"MOD_AM_LIST_ERR_04302": http.StatusInternalServerError,

	"MOD_AM_INA_EXI_04101": http.StatusOK,
	"MOD_AM_INA_ERR_04102": http.StatusInternalServerError,

	"MOD_AM_ACT_EXI_04201": http.StatusOK,
	"MOD_AM_ACT_ERR_04202": http.StatusInternalServerError,

	"FAM_CON_EXI_03201": http.StatusOK,
	"FAM_CON_ERR_03202": http.StatusNotFound,
	"FAM_CON_ERR_03203": http.StatusInternalServerError,

	"MOT_CON_EXI_03701": http.StatusOK,
	"MOT_CON_ERR_03702": http.StatusNotFound,
	"MOT_CON_ERR_03703": http.StatusInternalServerError,

	"MOT_LIST_EXI_03704": http.StatusOK,
	"MOT_LIST_ERR_03705": http.StatusInternalServerError,

	"FAB_CON_EXI_03101": http.StatusOK,
	"FAB_CON_ERR_03102": http.StatusNotFound,
	"FAB_CON_ERR_03103": http.StatusInternalServerError,

	"FAB_LIST_EXI_03104": http.StatusOK,
	"FAB_LIST_ERR_03105": http.StatusInternalServerError,


	"RUT_CON_EXI_03901": http.StatusOK,
	"RUT_CON_ERR_03902": http.StatusNotFound,
	"RUT_CON_ERR_03903": http.StatusInternalServerError,

	"RUT_LIST_EXI_03001": http.StatusOK,
	"RUT_LIST_ERR_03002": http.StatusInternalServerError,

	"RUT_AIR_CON_EXI_04001": http.StatusOK,
	"RUT_AIR_CON_ERR_04002": http.StatusNotFound,
	"RUT_AIR_CON_ERR_04003": http.StatusInternalServerError,

	"RUT_AIR_INA_EXI_04101": http.StatusOK,
	"RUT_AIR_INA_ERR_04102": http.StatusInternalServerError,


	"RUT_AIR_ACT_EXI_04201": http.StatusOK,
	"RUT_AIR_ACT_ERR_04202": http.StatusInternalServerError,


	"RUT_AIR_LIST_EXI_04001": http.StatusOK,
	"RUT_AIR_LIST_ERR_04002": http.StatusInternalServerError,


	"RUT_AIR_VAL_ERR_04301": http.StatusBadRequest,
	"RUT_AIR_VAL_ERR_04302": http.StatusBadRequest,


	"VUE_CON_EXI_04801": http.StatusOK,
	"VUE_CON_ERR_04802": http.StatusBadRequest,
	"VUE_CON_ERR_04803": http.StatusNotFound,
	"VUE_CON_ERR_04804": http.StatusInternalServerError,

	"VUE_EDI_EXI_04901": http.StatusOK,
	"VUE_EDI_ERR_04902": http.StatusBadRequest,
	"VUE_EDI_ERR_04903": http.StatusBadRequest,
	"VUE_EDI_ERR_04904": http.StatusInternalServerError,

	"VUE_REG_EXI_05001": http.StatusCreated,
	"VUE_REG_ERR_05002": http.StatusBadRequest,
	"VUE_REG_ERR_05003": http.StatusBadRequest,
	"VUE_REG_ERR_05004": http.StatusInternalServerError,

	"VUE_LIST_EXI_04800": http.StatusOK,
	"VUE_LIST_ERR_04802": http.StatusInternalServerError,

	"VUE_VAL_ERR_04805": http.StatusBadRequest,
	"VUE_VAL_ERR_04806": http.StatusBadRequest,
	"VUE_VAL_ERR_04807": http.StatusBadRequest,
	"VUE_VAL_ERR_04808": http.StatusBadRequest,

	"VUE_DEL_EXI_01801": http.StatusOK,
	"VUE_DEL_ERR_01802": http.StatusBadRequest,
	"VUE_DEL_ERR_01803": http.StatusNotFound,
	"VUE_DEL_ERR_01804": http.StatusInternalServerError,

	"VUE_AUTH_ERR_00001": http.StatusForbidden,


	"EMP_AIR_CON_EXI_02601": http.StatusOK,
	"EMP_AIR_CON_ERR_02602": http.StatusNotFound,
	"EMP_AIR_CON_ERR_02603": http.StatusInternalServerError,


	"EMP_AIR_EDI_EXI_02701": http.StatusOK,
	"EMP_AIR_EDI_ERR_02702": http.StatusInternalServerError,

	"EMP_AIR_AGR_EXI_02801": http.StatusCreated,
	"EMP_AIR_AGR_ERR_02802": http.StatusInternalServerError,
	"EMP_AIR_AGR_ERR_02803": http.StatusConflict,

	"EMP_AIR_ACT_EXI_02901": http.StatusOK,
	"EMP_AIR_ACT_ERR_02902": http.StatusInternalServerError,

	"EMP_AIR_INA_EXI_03001": http.StatusOK,
	"EMP_AIR_INA_ERR_03002": http.StatusInternalServerError,

	"EMP_AIR_LIST_EXI_02601": http.StatusOK,
	"EMP_AIR_LIST_ERR_02602": http.StatusInternalServerError,

	"EMP_AIR_VAL_ERR_02604": http.StatusBadRequest,
}

func (c *MessageCache) GetHTTPStatus(code string) int {
	if status, ok := messageCodeToHTTPStatus[code]; ok {
		return status
	}

	msg := c.GetMessage(code)
	if msg == nil {
		return http.StatusInternalServerError
	}

	switch msg.Type {
	case TypeSuccess:
		return http.StatusOK
	case TypeError:
		return http.StatusInternalServerError
	case TypeWarning:
		return http.StatusOK
	case TypeInfo:
		return http.StatusOK
	case TypeDebug:
		return http.StatusOK
	default:
		return http.StatusOK
	}
}

func (c *MessageCache) MessageCount() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.messages)
}
