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

// GetMessage retrieves a message by its code from cache
// If not found in cache, falls back to DB and caches it
func (c *MessageCache) GetMessage(code string) *CachedMessage {
	// Try cache first (read lock)
	c.mu.RLock()
	msg, found := c.messages[code]
	c.mu.RUnlock()

	if found {
		return msg
	}

	// Not in cache, try DB
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

		// Cache it for future use (write lock)
		c.mu.Lock()
		c.messages[code] = dbMsg
		c.mu.Unlock()

		log.Debug(logger.LogMsgCachedFromDB, "code", code)
		return dbMsg
	}

	// Not found in active messages, check if it exists but is inactive
	inactiveMsg, err := c.repo.GetByCodeWithStatusForCache(context.Background(), code)
	if err != nil {
		log.Warn(logger.LogMsgNotInDB, "code", code, "error", err)
		// Avoid infinite recursion
		if code == "GEN_MSG_INACTIVE_ERR_00002" {
			return nil
		}
		return c.GetMessage("GEN_MSG_INACTIVE_ERR_00002")
	}

	if inactiveMsg != nil && !inactiveMsg.Active {
		// Message exists but is inactive - return specific error message
		log.Warn(logger.LogMsgInactive, "code", code)
		return c.GetMessage("GEN_MSG_INACTIVE_ERR_00002")
	}

	// Message truly doesn't exist (not even in DB)
	log.Warn(logger.LogMsgNotInDB, "code", code)
	// Avoid infinite recursion
	if code == "GEN_MSG_INACTIVE_ERR_00002" {
		return nil
	}
	return c.GetMessage("GEN_MSG_INACTIVE_ERR_00002")
}

// GetMessageResponse retrieves formatted message response
func (c *MessageCache) GetMessageResponse(code string, params ...string) *MessageResponse {
	msg := c.GetMessage(code)
	if msg == nil {
		return nil
	}

	// Replace placeholders in content
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

// replaceAll is a simple helper for placeholder replacement
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

// messageCodeToHTTPStatus maps message codes to HTTP status codes
// Organized by modules: Users (MOD_U_*), Validation (MOD_V_*), Keycloak (MOD_KC_*),
// Infrastructure (MOD_INFRA_*), General (GEN_*), Messages (MOD_M_*)
var messageCodeToHTTPStatus = map[string]int{
	// ========================================
	// Users Module (MOD_U_*)
	// ========================================
	"MOD_U_REG_EXI_00001":        http.StatusCreated,      // 201 - Usuario registrado exitosamente
	"MOD_U_UPD_EXI_00002":        http.StatusOK,           // 200 - Usuario actualizado exitosamente
	"MOD_U_GET_EXI_00005":        http.StatusOK,           // 200 - Usuario encontrado
	"MOD_U_DUP_ERR_00001":        http.StatusConflict,     // 409 - Usuario duplicado
	"MOD_U_EMAIL_NF_ERR_00005":   http.StatusNotFound,     // 404 - Email no encontrado
	"MOD_U_GET_ERR_00003":        http.StatusNotFound,     // 404 - Usuario no encontrado
	"MOD_U_TOKEN_NF_ERR_00007":   http.StatusNotFound,     // 404 - Token no encontrado
	"MOD_U_EMAIL_NV_ERR_00006":   http.StatusForbidden,    // 403 - Email no verificado
	"MOD_U_TOKEN_EXP_ERR_00008":  http.StatusUnauthorized, // 401 - Token expirado
	"MOD_U_TOKEN_USED_ERR_00009": http.StatusUnauthorized, // 401 - Token ya usado
	// Update errors
	"MOD_U_UPD_ERR_00013":      http.StatusInternalServerError, // 500 - Error actualizando usuario
	"MOD_U_KC_UPD_ERR_00014":   http.StatusServiceUnavailable,  // 503 - Error sincronizando con Keycloak
	"MOD_U_ROLE_UPD_ERR_00015": http.StatusServiceUnavailable,  // 503 - Error actualizando rol en Keycloak
	// Delete
	"MOD_U_DEL_EXI_00003": http.StatusOK,                  // 200 - Usuario eliminado exitosamente
	"MOD_U_DEL_ERR_00012": http.StatusInternalServerError, // 500 - Error eliminando usuario

	// ========================================
	// Person Module (MOD_P_*)
	// ========================================
	"MOD_P_NOT_FOUND_ERR_00001": http.StatusNotFound, // 404 - Persona no encontrada

	// ========================================
	// Validation Module (MOD_V_*)
	// ========================================
	"MOD_V_VAL_ERR_00001":  http.StatusBadRequest, // 400 - Formato inválido
	"MOD_V_VAL_ERR_00002":  http.StatusBadRequest, // 400 - Request inválido
	"MOD_V_VAL_ERR_00006":  http.StatusBadRequest, // 400 - Validación fallida
	"MOD_V_VAL_ERR_00008":  http.StatusBadRequest, // 400 - Formato de campo
	"MOD_V_VAL_ERR_00009":  http.StatusBadRequest, // 400 - Campo requerido
	"MOD_V_VAL_ERR_00010":  http.StatusBadRequest, // 400 - Tipo de campo
	"MOD_V_VAL_ERR_00011":  http.StatusBadRequest, // 400 - Múltiples errores
	"MOD_V_JSON_ERR_00012": http.StatusBadRequest, // 400 - JSON inválido
	"MOD_V_ID_ERR_00013":   http.StatusBadRequest, // 400 - ID inválido
	// Data validation errors from DB constraints
	"MOD_V_FK_ERR_00014":   http.StatusUnprocessableEntity, // 422 - Invalid foreign key (e.g., airline doesn't exist)
	"MOD_V_LEN_ERR_00015":  http.StatusUnprocessableEntity, // 422 - Data too long for column
	"MOD_V_DATA_ERR_00016": http.StatusUnprocessableEntity, // 422 - Invalid data
	// Date validation errors
	"MOD_V_DATE_ERR_00017":  http.StatusBadRequest, // 400 - Start date after end date
	"MOD_V_DATE_ERR_00018":  http.StatusBadRequest, // 400 - Invalid date format
	"MOD_V_EMPTY_ERR_00019": http.StatusBadRequest, // 400 - Field contains only whitespace

	// ========================================
	// Keycloak Module (MOD_KC_*) - Email Verification & Auth
	// ========================================
	"MOD_KC_EMAIL_VERIFIED_EXI_00001":          http.StatusOK,                  // 200 - Email verificado exitosamente
	"MOD_KC_INVALID_TOKEN_ERR_00001":           http.StatusBadRequest,          // 400 - Token inválido/malformado
	"MOD_KC_EMAIL_VERIFY_ERROR_ERR_00001":      http.StatusInternalServerError, // 500 - Error de verificación (falla en Keycloak)
	"MOD_KC_USER_NOT_FOUND_ERR_00001":          http.StatusNotFound,            // 404 - Usuario no encontrado
	"MOD_KC_EMAIL_ALREADY_VERIFIED_WARN_00001": http.StatusOK,                  // 200 - Email ya verificado (warning)
	"MOD_KC_VERIF_EMAIL_SENT_EXI_00001":        http.StatusOK,                  // 200 - Email de verificación enviado
	"MOD_KC_VERIF_EMAIL_ERROR_ERR_00001":       http.StatusServiceUnavailable,  // 503 - Error enviando email
	"MOD_KC_VERIF_EMAIL_RESENT_EXI_00001":      http.StatusOK,                  // 200 - Email de verificación reenviado
	"MOD_KC_PWD_RESET_SENT_EXI_00001":          http.StatusOK,                  // 200 - Email de reset enviado
	"MOD_KC_PWD_RESET_ERROR_ERR_00001":         http.StatusServiceUnavailable,  // 503 - Error enviando reset
	// Login with email verification
	"MOD_KC_LOGIN_EMAIL_NOT_VERIFIED_ERR_00001": http.StatusUnauthorized, // 401 - Email no verificado, no puede hacer login
	"MOD_KC_LOGIN_SUCCESS_EXI_00001":            http.StatusOK,           // 200 - Login exitoso
	// Password update (via email token - forgot password flow)
	"MOD_KC_PWD_UPDATED_EXI_00001":              http.StatusOK,                  // 200 - Contraseña actualizada
	"MOD_KC_PWD_UPDATE_ERROR_ERR_00001":         http.StatusInternalServerError, // 500 - Error actualizando contraseña
	"MOD_KC_PWD_MISMATCH_ERR_00001":             http.StatusBadRequest,          // 400 - Contraseñas no coinciden
	"MOD_KC_PWD_UPDATE_TOKEN_INVALID_ERR_00001": http.StatusUnauthorized,        // 401 - Token de actualización inválido
	// Change password (authenticated user knows current password)
	"MOD_KC_PWD_CHANGED_EXI_00001":         http.StatusOK,                  // 200 - Contraseña cambiada exitosamente
	"MOD_KC_PWD_CHANGE_ERROR_ERR_00001":    http.StatusInternalServerError, // 500 - Error cambiando contraseña
	"MOD_KC_PWD_CURRENT_INVALID_ERR_00001": http.StatusUnauthorized,        // 401 - Contraseña actual incorrecta
	"MOD_KC_PWD_CHANGE_MISMATCH_ERR_00001": http.StatusBadRequest,          // 400 - Nuevas contraseñas no coinciden

	// ========================================
	// Authentication Module (MOD_AUTH_*)
	// ========================================
	"MOD_AUTH_LOGIN_SUCCESS_EXI_00001": http.StatusOK, // 200 - Login exitoso

	// ========================================
	// Infrastructure Module (MOD_INFRA_*)
	// ========================================
	"MOD_INFRA_KC_UNAVAIL_ERR_00004":      http.StatusLocked,              // 423 - Keycloak no disponible
	"MOD_INFRA_DB_UNAVAIL_ERR_00005":      http.StatusLocked,              // 423 - Base de datos no disponible
	"MOD_INFRA_DEP_FAIL_ERR_00006":        http.StatusLocked,              // 423 - Falla de dependencia
	"MOD_INFRA_KC_CLEANUP_ERR_00003":      http.StatusLocked,              // 423 - Error limpieza Keycloak
	"MOD_INFRA_KC_CREATE_ERR_00002":       http.StatusLocked,              // 423 - Error creación en Keycloak
	"MOD_INFRA_INCOMPLETE_REG_ERR_00009":  http.StatusConflict,            // 409 - Registro incompleto
	"MOD_INFRA_KC_INCONSISTENT_ERR_00001": http.StatusInternalServerError, // 500 - Estado inconsistente

	// ========================================
	// General Module (GEN_*)
	// ========================================
	"GEN_SRV_ERR_00001":          http.StatusInternalServerError, // 500 - Error interno del servidor
	"GEN_AUTH_ERR_00002":         http.StatusUnauthorized,        // 401 - No autorizado
	"GEN_FORBIDDEN_ERR_00003":    http.StatusForbidden,           // 403 - Acceso denegado
	"GEN_MSG_INACTIVE_ERR_00002": http.StatusServiceUnavailable,  // 503 - Mensaje no disponible

	// ========================================
	// Messages Module (MOD_M_*)
	// ========================================
	"MOD_M_CREATE_EXI_00001":    http.StatusCreated,    // 201 - Mensaje creado exitosamente
	"MOD_M_UPDATE_ERR_00010":    http.StatusBadRequest, // 400 - Error actualizando mensaje
	"MOD_M_NOT_FOUND_ERR_00001": http.StatusNotFound,   // 404 - Mensaje no encontrado

	// ========================================
	// Airline Module (MOD_AIR_*)
	// ========================================
	"MOD_AIR_GET_EXI_00001":        http.StatusOK,                  // 200 - Aerolínea obtenida exitosamente
	"MOD_AIR_ACTIVATE_EXI_00002":   http.StatusOK,                  // 200 - Aerolínea activada exitosamente
	"MOD_AIR_DEACTIVATE_EXI_00003": http.StatusOK,                  // 200 - Aerolínea desactivada exitosamente
	"MOD_AIR_NOT_FOUND_ERR_00001":  http.StatusNotFound,            // 404 - Aerolínea no encontrada
	"MOD_AIR_ACTIVATE_ERR_00002":   http.StatusUnprocessableEntity, // 422 - Error activando aerolínea (controlable)
	"MOD_AIR_DEACTIVATE_ERR_00003": http.StatusUnprocessableEntity, // 422 - Error desactivando aerolínea (controlable)

	// ========================================
	// Airport Module (MOD_APT_*)
	// ========================================
	"MOD_APT_GET_EXI_00001":        http.StatusOK,                  // 200 - Airport retrieved successfully
	"MOD_APT_ACTIVATE_EXI_00002":   http.StatusOK,                  // 200 - Airport activated successfully
	"MOD_APT_DEACTIVATE_EXI_00003": http.StatusOK,                  // 200 - Airport deactivated successfully
	"MOD_APT_NOT_FOUND_ERR_00001":  http.StatusNotFound,            // 404 - Airport not found
	"MOD_APT_ACTIVATE_ERR_00002":   http.StatusUnprocessableEntity, // 422 - Error activating airport
	"MOD_APT_DEACTIVATE_ERR_00003": http.StatusUnprocessableEntity, // 422 - Error deactivating airport

	// ========================================
	// City Module (CIU_*) - Ciudad (HU13 - Virtual Entity pattern)
	// ========================================
	"CIU_CON_EXI_01301": http.StatusOK,                  // 200 - City airports retrieved successfully
	"CIU_CON_ERR_01302": http.StatusNotFound,            // 404 - City not found (no airports in this city)
	"CIU_CON_ERR_01303": http.StatusInternalServerError, // 500 - Technical error querying city

	// ========================================
	// Country Module (PAI_*) - País (HU38 - Virtual Entity pattern)
	// ========================================
	"PAI_CON_EXI_03801": http.StatusOK,                  // 200 - Country airports retrieved successfully
	"PAI_CON_ERR_03802": http.StatusNotFound,            // 404 - Country not found (no airports in this country)
	"PAI_CON_ERR_03803": http.StatusInternalServerError, // 500 - Technical error querying country

	// ========================================
	// Airport Type Module (TAE_*) - Tipo Aeropuerto
	// ========================================
	"TAE_CON_EXI_04601": http.StatusOK,                  // 200 - Airport type retrieved successfully (airports of this type)
	"TAE_CON_ERR_04602": http.StatusNotFound,            // 404 - Airport type not found (no airports of this type)
	"TAE_CON_ERR_04603": http.StatusInternalServerError, // 500 - Technical error querying airport type

	// ========================================
	// Crew Member Type Module (TIN_*) - Tipo Integrante
	// ========================================
	"TIN_CON_EXI_04701": http.StatusOK,                  // 200 - Crew member type retrieved successfully (employees of this role)
	"TIN_CON_ERR_04702": http.StatusNotFound,            // 404 - Crew member type not found (no employees of this role)
	"TIN_CON_ERR_04703": http.StatusInternalServerError, // 500 - Technical error querying crew member type

	// ========================================
	// DailyLogbook Module (BIT_*) - Bitácora Diaria
	// ========================================
	// Consultar (HU7)
	"BIT_CON_EXI_01901": http.StatusOK,                  // 200 - Bitácora consultada exitosamente
	"BIT_CON_ERR_01902": http.StatusBadRequest,          // 400 - Bitácora no seleccionada
	"BIT_CON_ERR_01903": http.StatusNotFound,            // 404 - Bitácora no encontrada
	"BIT_CON_ERR_01904": http.StatusInternalServerError, // 500 - Error técnico al consultar

	// Agregar (HU8)
	"BIT_AGR_EXI_01801": http.StatusCreated,             // 201 - Bitácora creada exitosamente
	"BIT_AGR_ERR_01802": http.StatusBadRequest,          // 400 - Campos requeridos incompletos
	"BIT_AGR_ERR_01803": http.StatusBadRequest,          // 400 - Formato inválido
	"BIT_AGR_ERR_01804": http.StatusInternalServerError, // 500 - Error técnico al crear

	// Editar (HU9)
	"BIT_EDI_EXI_01701": http.StatusOK,                  // 200 - Bitácora actualizada exitosamente
	"BIT_EDI_ERR_01702": http.StatusBadRequest,          // 400 - Bitácora no seleccionada
	"BIT_EDI_ERR_01703": http.StatusBadRequest,          // 400 - Datos inválidos
	"BIT_EDI_ERR_01704": http.StatusInternalServerError, // 500 - Error técnico al editar

	// Eliminar (HU10)
	"BIT_DEL_EXI_01601": http.StatusOK,                  // 200 - Bitácora eliminada exitosamente
	"BIT_DEL_ERR_01602": http.StatusBadRequest,          // 400 - Bitácora no seleccionada
	"BIT_DEL_ERR_01603": http.StatusNotFound,            // 404 - Bitácora no existe o ya eliminada
	"BIT_DEL_ERR_01604": http.StatusInternalServerError, // 500 - Error técnico al eliminar

	// Activar (HU11)
	"BIT_ACT_EXI_01501": http.StatusOK,                  // 200 - Bitácora activada exitosamente
	"BIT_ACT_ERR_01502": http.StatusBadRequest,          // 400 - Bitácora no seleccionada
	"BIT_ACT_ERR_01503": http.StatusConflict,            // 409 - Bitácora ya está activa
	"BIT_ACT_ERR_01504": http.StatusInternalServerError, // 500 - Error técnico al activar

	// Inactivar (HU12)
	"BIT_INA_EXI_01401": http.StatusOK,                  // 200 - Bitácora inactivada exitosamente
	"BIT_INA_ERR_01402": http.StatusBadRequest,          // 400 - Bitácora no seleccionada
	"BIT_INA_ERR_01403": http.StatusConflict,            // 409 - Bitácora ya está inactiva
	"BIT_INA_ERR_01404": http.StatusInternalServerError, // 500 - Error técnico al inactivar

	// Listar
	"BIT_LIST_EXI_01001": http.StatusOK,                  // 200 - Lista obtenida exitosamente
	"BIT_LIST_ERR_01002": http.StatusInternalServerError, // 500 - Error al listar

	// Autorización
	"BIT_AUTH_ERR_00001": http.StatusForbidden, // 403 - No autorizado para esta bitácora

	// ========================================
	// Aircraft Registration Module (MAT_*) - Matrícula
	// ========================================
	// Consultar (HU33)
	"MAT_CON_EXI_03301": http.StatusOK,                  // 200 - Matrícula consultada exitosamente
	"MAT_CON_ERR_03302": http.StatusNotFound,            // 404 - Matrícula no encontrada
	"MAT_CON_ERR_03303": http.StatusInternalServerError, // 500 - Error técnico al consultar

	// Agregar (HU34)
	"MAT_AGR_EXI_03401": http.StatusCreated,    // 201 - Matrícula creada exitosamente
	"MAT_AGR_ERR_03402": http.StatusBadRequest, // 400 - Error al crear matrícula
	"MAT_AGR_ERR_03403": http.StatusConflict,   // 409 - Matrícula duplicada

	// Editar (HU35)
	"MAT_EDI_EXI_03501": http.StatusOK,         // 200 - Matrícula actualizada exitosamente
	"MAT_EDI_ERR_03502": http.StatusBadRequest, // 400 - Error al actualizar matrícula

	// Listar
	"MAT_LIST_EXI_03001": http.StatusOK,                  // 200 - Lista obtenida exitosamente
	"MAT_LIST_ERR_03002": http.StatusInternalServerError, // 500 - Error al listar

	// Validaciones
	"MAT_VAL_ERR_03601": http.StatusBadRequest, // 400 - Modelo de aeronave inválido
	"MAT_VAL_ERR_03602": http.StatusBadRequest, // 400 - Aerolínea inválida

	// ========================================
	// Aircraft Model Module (MOD_AM_*) - Modelo de Aeronave
	// ========================================
	// Consultar (HU36)
	"MOD_AM_CON_EXI_03601": http.StatusOK,                  // 200 - Modelo de aeronave consultado
	"MOD_AM_CON_ERR_03602": http.StatusNotFound,            // 404 - Modelo de aeronave no encontrado
	"MOD_AM_CON_ERR_03603": http.StatusInternalServerError, // 500 - Error técnico al consultar

	// Listar tipos (HU43)
	"MOD_AM_LIST_EXI_04301": http.StatusOK,                  // 200 - Lista de modelos/tipos obtenida
	"MOD_AM_LIST_ERR_04302": http.StatusInternalServerError, // 500 - Error al listar modelos/tipos

	// Inactivar (HU41)
	"MOD_AM_INA_EXI_04101": http.StatusOK,                  // 200 - Modelo de aeronave inactivado
	"MOD_AM_INA_ERR_04102": http.StatusInternalServerError, // 500 - Error técnico al inactivar

	// Activar (HU42)
	"MOD_AM_ACT_EXI_04201": http.StatusOK,                  // 200 - Modelo de aeronave activado
	"MOD_AM_ACT_ERR_04202": http.StatusInternalServerError, // 500 - Error técnico al activar

	// ========================================
	// Aircraft Family Module (FAM_*) - Familia de Aeronave (HU32)
	// ========================================
	"FAM_CON_EXI_03201": http.StatusOK,                  // 200 - Familia de aeronaves consultada exitosamente
	"FAM_CON_ERR_03202": http.StatusNotFound,            // 404 - Familia de aeronaves no encontrada
	"FAM_CON_ERR_03203": http.StatusInternalServerError, // 500 - Error técnico al consultar familia

	// ========================================
	// Engine Module (MOT_*) - Motor (HU37)
	// ========================================
	"MOT_CON_EXI_03701": http.StatusOK,                  // 200 - Motor consultado exitosamente
	"MOT_CON_ERR_03702": http.StatusNotFound,            // 404 - Motor no encontrado
	"MOT_CON_ERR_03703": http.StatusInternalServerError, // 500 - Error técnico al consultar

	"MOT_LIST_EXI_03704": http.StatusOK,                  // 200 - Lista de motores obtenida
	"MOT_LIST_ERR_03705": http.StatusInternalServerError, // 500 - Error al listar motores

	// ========================================
	// Manufacturer Module (FAB_*) - Fabricante (HU31)
	// ========================================
	"FAB_CON_EXI_03101": http.StatusOK,                  // 200 - Fabricante consultado exitosamente
	"FAB_CON_ERR_03102": http.StatusNotFound,            // 404 - Fabricante no encontrado
	"FAB_CON_ERR_03103": http.StatusInternalServerError, // 500 - Error técnico al consultar

	"FAB_LIST_EXI_03104": http.StatusOK,                  // 200 - Lista de fabricantes obtenida
	"FAB_LIST_ERR_03105": http.StatusInternalServerError, // 500 - Error al listar fabricantes

	// ========================================
	// Route Module (RUT_*) - Ruta
	// ========================================
	// Consultar (HU39)
	"RUT_CON_EXI_03901": http.StatusOK,                  // 200 - Ruta consultada exitosamente
	"RUT_CON_ERR_03902": http.StatusNotFound,            // 404 - Ruta no encontrada
	"RUT_CON_ERR_03903": http.StatusInternalServerError, // 500 - Error técnico al consultar

	// Listar
	"RUT_LIST_EXI_03001": http.StatusOK,                  // 200 - Lista de rutas obtenida
	"RUT_LIST_ERR_03002": http.StatusInternalServerError, // 500 - Error al listar rutas

	// ========================================
	// Airline Route Module (RUT_AIR_*) - Ruta Aerolínea
	// ========================================
	// Consultar (HU40)
	"RUT_AIR_CON_EXI_04001": http.StatusOK,                  // 200 - Ruta aerolínea consultada exitosamente
	"RUT_AIR_CON_ERR_04002": http.StatusNotFound,            // 404 - Ruta aerolínea no encontrada
	"RUT_AIR_CON_ERR_04003": http.StatusInternalServerError, // 500 - Error técnico al consultar

	// Desactivar (HU41)
	"RUT_AIR_INA_EXI_04101": http.StatusOK,                  // 200 - Ruta aerolínea desactivada
	"RUT_AIR_INA_ERR_04102": http.StatusInternalServerError, // 500 - Error técnico al desactivar

	// Activar (HU42)
	"RUT_AIR_ACT_EXI_04201": http.StatusOK,                  // 200 - Ruta aerolínea activada
	"RUT_AIR_ACT_ERR_04202": http.StatusInternalServerError, // 500 - Error técnico al activar

	// Listar
	"RUT_AIR_LIST_EXI_04001": http.StatusOK,                  // 200 - Lista de rutas aerolínea obtenida
	"RUT_AIR_LIST_ERR_04002": http.StatusInternalServerError, // 500 - Error al listar rutas aerolínea

	// Validaciones
	"RUT_AIR_VAL_ERR_04301": http.StatusBadRequest, // 400 - Ruta inválida
	"RUT_AIR_VAL_ERR_04302": http.StatusBadRequest, // 400 - Aerolínea inválida

	// ========================================
	// Flight Module (VUE_*) - Vuelo
	// ========================================
	// Consultar (HU48)
	"VUE_CON_EXI_04801": http.StatusOK,                  // 200 - Vuelo consultado exitosamente
	"VUE_CON_ERR_04802": http.StatusBadRequest,          // 400 - Vuelo no seleccionado
	"VUE_CON_ERR_04803": http.StatusNotFound,            // 404 - Vuelo no encontrado
	"VUE_CON_ERR_04804": http.StatusInternalServerError, // 500 - Error técnico al consultar

	// Editar (HU49)
	"VUE_EDI_EXI_04901": http.StatusOK,                  // 200 - Vuelo actualizado exitosamente
	"VUE_EDI_ERR_04902": http.StatusBadRequest,          // 400 - Vuelo no seleccionado
	"VUE_EDI_ERR_04903": http.StatusBadRequest,          // 400 - Datos inválidos
	"VUE_EDI_ERR_04904": http.StatusInternalServerError, // 500 - Error técnico al editar

	// Registrar (HU50)
	"VUE_REG_EXI_05001": http.StatusCreated,             // 201 - Vuelo registrado exitosamente
	"VUE_REG_ERR_05002": http.StatusBadRequest,          // 400 - Campos requeridos incompletos
	"VUE_REG_ERR_05003": http.StatusBadRequest,          // 400 - Formato inválido
	"VUE_REG_ERR_05004": http.StatusInternalServerError, // 500 - Error técnico al registrar

	// Listar
	"VUE_LIST_EXI_04800": http.StatusOK,                  // 200 - Lista de vuelos obtenida
	"VUE_LIST_ERR_04802": http.StatusInternalServerError, // 500 - Error al listar vuelos

	// Validaciones
	"VUE_VAL_ERR_04805": http.StatusBadRequest, // 400 - Ruta de aerolínea inválida
	"VUE_VAL_ERR_04806": http.StatusBadRequest, // 400 - Bitácora inválida
	"VUE_VAL_ERR_04807": http.StatusBadRequest, // 400 - Matrícula de aeronave inválida
	"VUE_VAL_ERR_04808": http.StatusBadRequest, // 400 - Secuencia de tiempos inválida

	// Eliminar (HU18)
	"VUE_DEL_EXI_01801": http.StatusOK,                  // 200 - Vuelo eliminado exitosamente
	"VUE_DEL_ERR_01802": http.StatusBadRequest,          // 400 - Vuelo no seleccionado
	"VUE_DEL_ERR_01803": http.StatusNotFound,            // 404 - Vuelo no existe o ya eliminado
	"VUE_DEL_ERR_01804": http.StatusInternalServerError, // 500 - Error técnico al eliminar

	// Autorización
	"VUE_AUTH_ERR_00001": http.StatusForbidden, // 403 - No autorizado para este vuelo

	// ========================================
	// Airline Employee Module (EMP_AIR_*) - Empleado Aerolínea
	// ========================================
	// Consultar (HU26)
	"EMP_AIR_CON_EXI_02601": http.StatusOK,                  // 200 - Airline employee retrieved successfully
	"EMP_AIR_CON_ERR_02602": http.StatusNotFound,            // 404 - Airline employee not found
	"EMP_AIR_CON_ERR_02603": http.StatusInternalServerError, // 500 - Technical error querying

	// Editar (HU27)
	"EMP_AIR_EDI_EXI_02701": http.StatusOK,                  // 200 - Airline employee updated successfully
	"EMP_AIR_EDI_ERR_02702": http.StatusInternalServerError, // 500 - Technical error updating

	// Agregar (HU28)
	"EMP_AIR_AGR_EXI_02801": http.StatusCreated,             // 201 - Airline employee created successfully
	"EMP_AIR_AGR_ERR_02802": http.StatusInternalServerError, // 500 - Technical error creating
	"EMP_AIR_AGR_ERR_02803": http.StatusConflict,            // 409 - Duplicate airline employee

	// Activar (HU29)
	"EMP_AIR_ACT_EXI_02901": http.StatusOK,                  // 200 - Airline employee activated successfully
	"EMP_AIR_ACT_ERR_02902": http.StatusInternalServerError, // 500 - Technical error activating

	// Inactivar (HU30)
	"EMP_AIR_INA_EXI_03001": http.StatusOK,                  // 200 - Airline employee deactivated successfully
	"EMP_AIR_INA_ERR_03002": http.StatusInternalServerError, // 500 - Technical error deactivating

	// Listar
	"EMP_AIR_LIST_EXI_02601": http.StatusOK,                  // 200 - Airline employee list retrieved successfully
	"EMP_AIR_LIST_ERR_02602": http.StatusInternalServerError, // 500 - Error listing airline employees

	// Validaciones
	"EMP_AIR_VAL_ERR_02604": http.StatusBadRequest, // 400 - Invalid airline
}

// GetHTTPStatus returns the HTTP status for a message code
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

// MessageCount returns the number of loaded messages in cache
func (c *MessageCache) MessageCount() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.messages)
}
