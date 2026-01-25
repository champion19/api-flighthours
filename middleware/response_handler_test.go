package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/champion19/flighthours-api/platform/cache/messaging"
	cachetypes "github.com/champion19/flighthours-api/platform/cache/types"
	"github.com/gin-gonic/gin"
)

type fakeMessageCacheRepoForResponse struct {
	messages []cachetypes.CachedMessage
}

func (r fakeMessageCacheRepoForResponse) GetAllActiveForCache(context.Context) ([]cachetypes.CachedMessage, error) {
	return r.messages, nil
}
func (r fakeMessageCacheRepoForResponse) GetByCodeForCache(context.Context, string) (*cachetypes.CachedMessage, error) {
	return nil, nil
}
func (r fakeMessageCacheRepoForResponse) GetByCodeWithStatusForCache(context.Context, string) (*cachetypes.CachedMessage, error) {
	return nil, nil
}

func setupResponseHandler(t *testing.T) *ResponseHandler {
	t.Helper()
	repo := fakeMessageCacheRepoForResponse{messages: []cachetypes.CachedMessage{
		{Code: domain.MsgUserRegistered, Type: cachetypes.TypeSuccess, Content: "User registered successfully"},
		{Code: domain.MsgUserDuplicate, Type: cachetypes.TypeError, Content: "User already exists"},
	}}
	cache := messaging.NewMessageCache(repo, 0)
	if err := cache.LoadMessages(context.Background()); err != nil {
		t.Fatalf("failed to load cache: %v", err)
	}
	return NewResponseHandler(cache)
}

func TestResponseHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := setupResponseHandler(t)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	handler.Success(c, domain.MsgUserRegistered)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, w.Code)
	}
}

func TestResponseHandler_SuccessWithData(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := setupResponseHandler(t)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	data := map[string]string{"id": "123"}
	handler.SuccessWithData(c, domain.MsgUserRegistered, data)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, w.Code)
	}
}

func TestResponseHandler_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := setupResponseHandler(t)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	handler.Error(c, domain.MsgUserDuplicate)

	if w.Code != http.StatusConflict {
		t.Fatalf("expected status %d, got %d", http.StatusConflict, w.Code)
	}
}

func TestResponseHandler_ErrorWithData(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := setupResponseHandler(t)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	data := map[string]string{"field": "email"}
	handler.ErrorWithData(c, domain.MsgUserDuplicate, data)

	if w.Code != http.StatusConflict {
		t.Fatalf("expected status %d, got %d", http.StatusConflict, w.Code)
	}
}

func TestResponseHandler_Warning(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := setupResponseHandler(t)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	handler.Warning(c, domain.MsgUserRegistered)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestResponseHandler_WarningWithData(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := setupResponseHandler(t)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	data := map[string]string{"warning": "deprecated"}
	handler.WarningWithData(c, domain.MsgUserRegistered, data)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestResponseHandler_Info(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := setupResponseHandler(t)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	handler.Info(c, domain.MsgUserRegistered)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestResponseHandler_InfoWithData(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := setupResponseHandler(t)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	data := map[string]string{"info": "processing"}
	handler.InfoWithData(c, domain.MsgUserRegistered, data)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestResponseHandler_DataOnly(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := setupResponseHandler(t)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	data := map[string]string{"result": "data"}
	handler.DataOnly(c, data)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}
