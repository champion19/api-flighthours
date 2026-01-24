package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/champion19/flighthours-api/platform/cache/messaging"
	cachetypes "github.com/champion19/flighthours-api/platform/cache/types"
	"github.com/champion19/flighthours-api/platform/logger"
	"github.com/gin-gonic/gin"
)

type noopLogger struct{}

func (noopLogger) Info(string, ...any)    {}
func (noopLogger) Error(string, ...any)   {}
func (noopLogger) Debug(string, ...any)   {}
func (noopLogger) Warn(string, ...any)    {}
func (noopLogger) Success(string, ...any) {}
func (noopLogger) Fatal(string, ...any)   {}
func (noopLogger) Panic(string, ...any)   {}
func (noopLogger) WithTraceID(string) logger.Logger {
	return noopLogger{}
}

type fakeMessageCacheRepo struct {
	messages []cachetypes.CachedMessage
}

func (r fakeMessageCacheRepo) GetAllActiveForCache(context.Context) ([]cachetypes.CachedMessage, error) {
	return r.messages, nil
}
func (r fakeMessageCacheRepo) GetByCodeForCache(context.Context, string) (*cachetypes.CachedMessage, error) {
	return nil, nil
}
func (r fakeMessageCacheRepo) GetByCodeWithStatusForCache(context.Context, string) (*cachetypes.CachedMessage, error) {
	return nil, nil
}

func newTestCache(t *testing.T) *messaging.MessageCache {
	t.Helper()

	repo := fakeMessageCacheRepo{messages: []cachetypes.CachedMessage{
		{Code: domain.MsgUserDuplicate, Type: cachetypes.TypeError, Content: "duplicate"},
		{Code: domain.MsgValFieldRequired, Type: cachetypes.TypeError, Content: "missing ${0}"},
		{Code: domain.MsgPersonNotFound, Type: cachetypes.TypeError, Content: "person not found"},
		{Code: domain.MsgMessageNotFound, Type: cachetypes.TypeError, Content: "message not found"},
		{Code: domain.MsgMessageCodeDuplicate, Type: cachetypes.TypeError, Content: "message code duplicate"},
		{Code: domain.MsgValIDInvalid, Type: cachetypes.TypeError, Content: "invalid ID"},
	}}
	c := messaging.NewMessageCache(repo, 0)
	if err := c.LoadMessages(context.Background()); err != nil {
		t.Fatalf("load cache: %v", err)
	}
	return c
}

func TestErrorHandler_MappedError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cache := newTestCache(t)
	errHandler := NewErrorHandler(cache)

	r := gin.New()
	r.Use(errHandler.Handle())
	r.GET("/", func(c *gin.Context) {
		c.Error(domain.ErrDuplicateUser)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Fatalf("expected %d got %d body=%s", http.StatusConflict, w.Code, w.Body.String())
	}

	var resp ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if resp.Success {
		t.Fatalf("expected success=false")
	}
	if resp.Code != domain.MsgUserDuplicate {
		t.Fatalf("expected code %q got %q", domain.MsgUserDuplicate, resp.Code)
	}
}

func TestErrorHandler_ValidationFieldsParam(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cache := newTestCache(t)
	errHandler := NewErrorHandler(cache)

	r := gin.New()
	r.Use(errHandler.Handle())
	r.GET("/", func(c *gin.Context) {
		c.Set("validation_fields", []string{"a", "b"})
		c.Error(domain.ErrSchemaFieldRequired)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected %d got %d body=%s", http.StatusBadRequest, w.Code, w.Body.String())
	}

	var resp ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if resp.Message != "missing a, b" {
		t.Fatalf("expected substituted message, got %q", resp.Message)
	}
}

func TestErrorHandler_UnmappedError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cache := newTestCache(t)
	errHandler := NewErrorHandler(cache)

	r := gin.New()
	r.Use(errHandler.Handle())
	r.GET("/", func(c *gin.Context) {
		c.Error(errors.New("unmapped custom error"))
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected %d got %d", http.StatusInternalServerError, w.Code)
	}

	var resp ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if resp.Success {
		t.Fatal("expected success=false")
	}
	if resp.Code != domain.MsgServerError {
		t.Fatalf("expected code %q got %q", domain.MsgServerError, resp.Code)
	}
}

func TestErrorHandler_NoErrors(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cache := newTestCache(t)
	errHandler := NewErrorHandler(cache)

	r := gin.New()
	r.Use(errHandler.Handle())
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected %d got %d", http.StatusOK, w.Code)
	}
}

func TestErrorHandler_SingleFieldValidation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cache := newTestCache(t)
	errHandler := NewErrorHandler(cache)

	r := gin.New()
	r.Use(errHandler.Handle())
	r.GET("/", func(c *gin.Context) {
		c.Set("validation_fields", []string{"email"})
		c.Error(domain.ErrSchemaFieldRequired)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected %d got %d", http.StatusBadRequest, w.Code)
	}

	var resp ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if resp.Message != "missing email" {
		t.Fatalf("expected 'missing email', got %q", resp.Message)
	}
}

func TestErrorHandler_PersonNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cache := newTestCache(t)
	errHandler := NewErrorHandler(cache)

	r := gin.New()
	r.Use(errHandler.Handle())
	r.GET("/", func(c *gin.Context) {
		c.Error(domain.ErrPersonNotFound)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	r.ServeHTTP(w, req)

	var resp ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if resp.Success {
		t.Fatal("expected success=false")
	}
	if resp.Code != domain.MsgPersonNotFound {
		t.Fatalf("expected code %s, got %s", domain.MsgPersonNotFound, resp.Code)
	}
}

func TestErrorHandler_MessageErrors(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cache := newTestCache(t)
	errHandler := NewErrorHandler(cache)

	tests := []struct {
		name         string
		err          error
		expectedCode string
	}{
		{"message not found", domain.ErrMessageNotFound, domain.MsgMessageNotFound},
		{"message code duplicate", domain.ErrMessageCodeDuplicate, domain.MsgMessageCodeDuplicate},
		{"invalid ID", domain.ErrInvalidID, domain.MsgValIDInvalid},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			r.Use(errHandler.Handle())
			r.GET("/", func(c *gin.Context) {
				c.Error(tt.err)
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			r.ServeHTTP(w, req)

			var resp ErrorResponse
			if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
				t.Fatalf("invalid json: %v", err)
			}
			if resp.Success {
				t.Fatal("expected success=false")
			}
			if resp.Code != tt.expectedCode {
				t.Fatalf("expected code %s, got %s", tt.expectedCode, resp.Code)
			}
		})
	}
}

func TestErrorHandler_ValidationFieldsNonSlice(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cache := newTestCache(t)
	errHandler := NewErrorHandler(cache)

	r := gin.New()
	r.Use(errHandler.Handle())
	r.GET("/", func(c *gin.Context) {
		c.Set("validation_fields", "not a slice")
		c.Error(domain.ErrSchemaFieldRequired)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected %d got %d", http.StatusBadRequest, w.Code)
	}
}
