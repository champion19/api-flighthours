package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/champion19/flighthours-api/core/interactor"
	"github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/champion19/flighthours-api/core/ports/input"
	"github.com/champion19/flighthours-api/core/ports/output"
	"github.com/champion19/flighthours-api/middleware"
	"github.com/champion19/flighthours-api/platform/cache/messaging"
	cachetypes "github.com/champion19/flighthours-api/platform/cache/types"
	"github.com/champion19/flighthours-api/tools/idencoder"
	"github.com/gin-gonic/gin"
)

type fakeMessageService struct {
	validateErr error
	getByIDRes  *domain.Message
	getByIDErr  error
	listRes     []domain.Message
	listErr     error
}

func (f *fakeMessageService) BeginTx(context.Context) (output.Tx, error) { return fakeTx{}, nil }
func (f *fakeMessageService) ValidateMessage(context.Context, domain.Message) error {
	return f.validateErr
}
func (f *fakeMessageService) GetMessageByID(ctx context.Context, id string) (*domain.Message, error) {
	return f.getByIDRes, f.getByIDErr
}
func (f *fakeMessageService) GetMessageByCode(context.Context, string) (*domain.Message, error) {
	return nil, errors.New("not implemented")
}
func (f *fakeMessageService) ListMessages(ctx context.Context, filters map[string]interface{}) ([]domain.Message, error) {
	return f.listRes, f.listErr
}
func (f *fakeMessageService) ListActiveMessages(context.Context) ([]domain.Message, error) {
	return f.listRes, f.listErr
}
func (f *fakeMessageService) SaveMessageToDB(context.Context, output.Tx, domain.Message) error {
	return nil
}
func (f *fakeMessageService) UpdateMessageInDB(context.Context, output.Tx, domain.Message) error {
	return nil
}
func (f *fakeMessageService) DeleteMessageFromDB(context.Context, output.Tx, string) error {
	return nil
}

func newMessageRouter(msgSvc input.MessageService) *gin.Engine {
	// Inline cache creation without testing.T dependency
	repo := fakeMessageCacheRepo{messages: []cachetypes.CachedMessage{
		{Code: domain.MsgUserRegistered, Type: cachetypes.TypeSuccess, Content: "user registered"},
		{Code: domain.MsgValJSONInvalid, Type: cachetypes.TypeError, Content: "invalid json"},
		{Code: domain.MsgUserDuplicate, Type: cachetypes.TypeError, Content: "duplicate"},
		{Code: domain.MsgIncompleteRegistration, Type: cachetypes.TypeError, Content: "incomplete"},
		{Code: domain.MsgMessageNotFound, Type: cachetypes.TypeError, Content: "message not found"},
		{Code: domain.MsgValIDInvalid, Type: cachetypes.TypeError, Content: "invalid ID"},
	}}
	cache := messaging.NewMessageCache(repo, 0)
	_ = cache.LoadMessages(context.Background())

	resp := middleware.NewResponseHandler(cache)
	errHandler := middleware.NewErrorHandler(cache)

	enc, _ := idencoder.NewHashidsEncoder(idencoder.Config{Secret: "test-secret", MinLength: 10}, noopLogger{})

	msgInter := interactor.NewMessageInteractor(msgSvc, noopLogger{})
	h := New(nil, nil, enc, resp, msgInter, cache)

	r := gin.New()
	r.Use(middleware.RequestID())
	r.Use(errHandler.Handle())

	r.POST("/messages", h.CreateMessage())
	r.GET("/messages/:id", h.GetMessageByID())
	r.GET("/messages", h.ListMessages())
	r.PUT("/messages/:id", h.UpdateMessage())
	r.DELETE("/messages/:id", h.DeleteMessage())
	r.POST("/messages/cache/reload", h.ReloadMessageCache())

	return r
}

func TestHTTP_CreateMessage(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		svc := &fakeMessageService{}
		r := newMessageRouter(svc)

		body := map[string]any{
			"code":     "TEST_001",
			"type":     "success",
			"category": "test",
			"module":   "test",
			"title":    "Test",
			"content":  "Test message",
			"active":   true,
		}
		b, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/messages", bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expected status %d, got %d. body=%s", http.StatusOK, w.Code, w.Body.String())
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		svc := &fakeMessageService{}
		r := newMessageRouter(svc)

		req := httptest.NewRequest(http.MethodPost, "/messages", bytes.NewBufferString("{"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
}

func TestHTTP_GetMessageByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create encoder with valid UUID
	enc, _ := idencoder.NewHashidsEncoder(idencoder.Config{Secret: "test-secret", MinLength: 10}, noopLogger{})
	testUUID := "550e8400-e29b-41d4-a716-446655440000"
	validID, _ := enc.Encode(testUUID)

	t.Run("success", func(t *testing.T) {
		msg := &domain.Message{
			ID:       testUUID,
			Code:     "TEST_001",
			Type:     "success",
			Category: "test",
			Module:   "test",
			Title:    "Test",
			Content:  "Test message",
			Active:   true,
		}
		svc := &fakeMessageService{getByIDRes: msg}
		r := newMessageRouter(svc)

		req := httptest.NewRequest(http.MethodGet, "/messages/"+validID, nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expected status %d, got %d. body=%s", http.StatusOK, w.Code, w.Body.String())
		}
	})

	t.Run("not found", func(t *testing.T) {
		svc := &fakeMessageService{getByIDErr: domain.ErrMessageNotFound}
		r := newMessageRouter(svc)

		req := httptest.NewRequest(http.MethodGet, "/messages/"+validID, nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		if w.Code != http.StatusNotFound {
			t.Fatalf("expected status %d, got %d", http.StatusNotFound, w.Code)
		}
	})
}

func TestHTTP_ListMessages(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Use valid UUIDs for message IDs
	uuid1 := "550e8400-e29b-41d4-a716-446655440001"
	uuid2 := "550e8400-e29b-41d4-a716-446655440002"

	t.Run("success", func(t *testing.T) {
		msgs := []domain.Message{
			{ID: uuid1, Code: "TEST_001", Type: "success", Title: "Test 1", Content: "Content 1", Active: true},
			{ID: uuid2, Code: "TEST_002", Type: "error", Title: "Test 2", Content: "Content 2", Active: true},
		}
		svc := &fakeMessageService{listRes: msgs}
		r := newMessageRouter(svc)

		req := httptest.NewRequest(http.MethodGet, "/messages", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expected status %d, got %d. body=%s", http.StatusOK, w.Code, w.Body.String())
		}
	})

	t.Run("with filters", func(t *testing.T) {
		msgs := []domain.Message{
			{ID: uuid1, Code: "TEST_001", Type: "success", Module: "users", Active: true},
		}
		svc := &fakeMessageService{listRes: msgs}
		r := newMessageRouter(svc)

		req := httptest.NewRequest(http.MethodGet, "/messages?module=users&type=success", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
		}
	})
}

func TestHTTP_UpdateMessage(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create encoder with valid UUID
	enc, _ := idencoder.NewHashidsEncoder(idencoder.Config{Secret: "test-secret", MinLength: 10}, noopLogger{})
	testUUID := "550e8400-e29b-41d4-a716-446655440000"
	validID, _ := enc.Encode(testUUID)

	t.Run("success", func(t *testing.T) {
		svc := &fakeMessageService{}
		r := newMessageRouter(svc)

		body := map[string]any{
			"code":     "TEST_001",
			"type":     "success",
			"category": "test",
			"module":   "test",
			"title":    "Updated",
			"content":  "Updated message",
			"active":   true,
		}
		b, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPut, "/messages/"+validID, bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expected status %d, got %d. body=%s", http.StatusOK, w.Code, w.Body.String())
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		svc := &fakeMessageService{}
		r := newMessageRouter(svc)

		req := httptest.NewRequest(http.MethodPut, "/messages/"+validID, bytes.NewBufferString("{"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
}

func TestHTTP_DeleteMessage(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create encoder with valid UUID
	enc, _ := idencoder.NewHashidsEncoder(idencoder.Config{Secret: "test-secret", MinLength: 10}, noopLogger{})
	testUUID := "550e8400-e29b-41d4-a716-446655440000"
	validID, _ := enc.Encode(testUUID)

	t.Run("success", func(t *testing.T) {
		svc := &fakeMessageService{}
		r := newMessageRouter(svc)

		req := httptest.NewRequest(http.MethodDelete, "/messages/"+validID, nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expected status %d, got %d. body=%s", http.StatusOK, w.Code, w.Body.String())
		}
	})
}

func TestHTTP_ReloadMessageCache(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		svc := &fakeMessageService{}
		r := newMessageRouter(svc)

		req := httptest.NewRequest(http.MethodPost, "/messages/cache/reload", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expected status %d, got %d. body=%s", http.StatusOK, w.Code, w.Body.String())
		}
	})
}
