package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Nerzal/gocloak/v13"
	"github.com/champion19/api-flighthours/core/interactor"
	"github.com/champion19/api-flighthours/core/interactor/dto"
	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/input"
	"github.com/champion19/api-flighthours/core/ports/output"
	"github.com/champion19/api-flighthours/middleware"
	"github.com/champion19/api-flighthours/platform/cache/messaging"
	cachetypes "github.com/champion19/api-flighthours/platform/cache/types"
	platformLogger "github.com/champion19/api-flighthours/platform/logger"
	"github.com/champion19/api-flighthours/tools/idencoder"
	"github.com/gin-gonic/gin"
)

type noopLogger struct{}

func (noopLogger) Info(string, ...any)                      {}
func (noopLogger) Error(string, ...any)                     {}
func (noopLogger) Debug(string, ...any)                     {}
func (noopLogger) Warn(string, ...any)                      {}
func (noopLogger) Success(string, ...any)                   {}
func (noopLogger) Fatal(string, ...any)                     {}
func (noopLogger) Panic(string, ...any)                     {}
func (noopLogger) WithTraceID(string) platformLogger.Logger { return noopLogger{} }

type fakeTx struct{}

func (fakeTx) Commit() error   { return nil }
func (fakeTx) Rollback() error { return nil }

type fakeService struct{}

var _ input.Service = (*fakeService)(nil)

func (f *fakeService) BeginTx(context.Context) (output.Tx, error) { return fakeTx{}, nil }
func (f *fakeService) RegisterEmployee(ctx context.Context, employee domain.Employee) (*dto.RegisterEmployee, error) {
	return &dto.RegisterEmployee{Employee: employee, Message: "ok"}, nil
}
func (f *fakeService) GetEmployeeByEmail(context.Context, string) (*domain.Employee, error) {
	return nil, errors.New("not implemented")
}
func (f *fakeService) GetEmployeeByID(context.Context, string) (*domain.Employee, error) {
	return nil, errors.New("not implemented")
}
func (f *fakeService) LocateEmployee(context.Context, string) (*dto.RegisterEmployee, error) {
	return nil, errors.New("not implemented")
}
func (f *fakeService) CheckAndCleanInconsistentState(context.Context, string) error { return nil }
func (f *fakeService) SaveEmployeeToDB(context.Context, output.Tx, domain.Employee) error {
	return nil
}
func (f *fakeService) UpdateEmployeeKeycloakID(context.Context, output.Tx, string, string) error {
	return nil
}
func (f *fakeService) CreateUserInKeycloak(context.Context, *domain.Employee) (string, error) {
	return "kc1", nil
}
func (f *fakeService) SetUserPassword(context.Context, string, string) error { return nil }
func (f *fakeService) AssignUserRole(context.Context, string, string) error  { return nil }
func (f *fakeService) SendVerificationEmail(context.Context, string) error   { return nil }
func (f *fakeService) SendPasswordResetEmail(context.Context, string) error  { return nil }
func (f *fakeService) Login(context.Context, string, string) (*gocloak.JWT, error) {
	return &gocloak.JWT{}, nil
}
func (f *fakeService) VerifyEmailByToken(context.Context, string) (string, error) {
	return "", nil
}
func (f *fakeService) GetUserByEmail(context.Context, string) (*gocloak.User, error) {
	return nil, errors.New("not implemented")
}
func (f *fakeService) RollbackEmployee(context.Context, string) error     { return nil }
func (f *fakeService) RollbackKeycloakUser(context.Context, string) error { return nil }
func (f *fakeService) UpdatePassword(context.Context, string, string) (string, error) {
	return "", nil
}
func (f *fakeService) UpdateEmployee(context.Context, output.Tx, domain.Employee) error {
	return nil
}
func (f *fakeService) ChangePassword(context.Context, string, string, string) (string, error) {
	return "", nil
}
func (f *fakeService) DeleteEmployee(context.Context, string, string) error {
	return nil
}
func (f *fakeService) GetEmployeeByKeycloakID(context.Context, string) (*domain.Employee, error) {
	return nil, errors.New("not implemented")
}

// GetEmployeesByRole mock for HU47 - Virtual Entity pattern
func (f *fakeService) GetEmployeesByRole(context.Context, string) ([]domain.Employee, error) {
	return nil, errors.New("not implemented")
}

type fakeServiceErr struct {
	err error
}

var _ input.Service = (*fakeServiceErr)(nil)

func (f *fakeServiceErr) BeginTx(context.Context) (output.Tx, error) { return fakeTx{}, nil }
func (f *fakeServiceErr) RegisterEmployee(context.Context, domain.Employee) (*dto.RegisterEmployee, error) {
	return nil, f.err
}
func (f *fakeServiceErr) GetEmployeeByEmail(context.Context, string) (*domain.Employee, error) {
	return nil, errors.New("not implemented")
}
func (f *fakeServiceErr) GetEmployeeByID(context.Context, string) (*domain.Employee, error) {
	return nil, errors.New("not implemented")
}
func (f *fakeServiceErr) LocateEmployee(context.Context, string) (*dto.RegisterEmployee, error) {
	return nil, errors.New("not implemented")
}
func (f *fakeServiceErr) CheckAndCleanInconsistentState(context.Context, string) error { return nil }
func (f *fakeServiceErr) SaveEmployeeToDB(context.Context, output.Tx, domain.Employee) error {
	return nil
}
func (f *fakeServiceErr) UpdateEmployeeKeycloakID(context.Context, output.Tx, string, string) error {
	return nil
}
func (f *fakeServiceErr) CreateUserInKeycloak(context.Context, *domain.Employee) (string, error) {
	return "", errors.New("not implemented")
}
func (f *fakeServiceErr) SetUserPassword(context.Context, string, string) error { return nil }
func (f *fakeServiceErr) AssignUserRole(context.Context, string, string) error  { return nil }
func (f *fakeServiceErr) SendVerificationEmail(context.Context, string) error   { return nil }
func (f *fakeServiceErr) SendPasswordResetEmail(context.Context, string) error  { return nil }
func (f *fakeServiceErr) Login(context.Context, string, string) (*gocloak.JWT, error) {
	return &gocloak.JWT{}, nil
}
func (f *fakeServiceErr) VerifyEmailByToken(context.Context, string) (string, error) {
	return "", nil
}
func (f *fakeServiceErr) GetUserByEmail(context.Context, string) (*gocloak.User, error) {
	return nil, errors.New("not implemented")
}
func (f *fakeServiceErr) RollbackEmployee(context.Context, string) error     { return nil }
func (f *fakeServiceErr) RollbackKeycloakUser(context.Context, string) error { return nil }
func (f *fakeServiceErr) UpdatePassword(context.Context, string, string) (string, error) {
	return "", nil
}
func (f *fakeServiceErr) UpdateEmployee(context.Context, output.Tx, domain.Employee) error {
	return nil
}
func (f *fakeServiceErr) ChangePassword(context.Context, string, string, string) (string, error) {
	return "", nil
}
func (f *fakeServiceErr) DeleteEmployee(context.Context, string, string) error {
	return nil
}
func (f *fakeServiceErr) GetEmployeeByKeycloakID(context.Context, string) (*domain.Employee, error) {
	return nil, errors.New("not implemented")
}

func (f *fakeServiceErr) GetEmployeesByRole(context.Context, string) ([]domain.Employee, error) {
	return nil, errors.New("not implemented")
}

type fakeMessageCacheRepo struct {
	messages  []cachetypes.CachedMessage
	reloadErr error
}

func (r fakeMessageCacheRepo) GetAllActiveForCache(context.Context) ([]cachetypes.CachedMessage, error) {
	if r.reloadErr != nil {
		return nil, r.reloadErr
	}
	return r.messages, nil
}
func (r fakeMessageCacheRepo) GetByCodeForCache(context.Context, string) (*cachetypes.CachedMessage, error) {
	return nil, errors.New("not implemented")
}
func (r fakeMessageCacheRepo) GetByCodeWithStatusForCache(context.Context, string) (*cachetypes.CachedMessage, error) {
	return nil, errors.New("not implemented")
}

func newTestMessageCache(t *testing.T) *messaging.MessageCache {
	t.Helper()

	repo := fakeMessageCacheRepo{messages: []cachetypes.CachedMessage{
		{Code: domain.MsgUserRegistered, Type: cachetypes.TypeSuccess, Content: "user registered"},
		{Code: domain.MsgValJSONInvalid, Type: cachetypes.TypeError, Content: "invalid json"},
		{Code: domain.MsgUserDuplicate, Type: cachetypes.TypeError, Content: "duplicate"},
		{Code: domain.MsgIncompleteRegistration, Type: cachetypes.TypeError, Content: "incomplete"},
	}}
	cache := messaging.NewMessageCache(repo, 0)
	if err := cache.LoadMessages(context.Background()); err != nil {
		t.Fatalf("failed to load message cache: %v", err)
	}
	return cache
}

func TestHTTP_RegisterEmployee(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cache := newTestMessageCache(t)
	resp := middleware.NewResponseHandler(cache)
	errHandler := middleware.NewErrorHandler(cache)

	enc, err := idencoder.NewHashidsEncoder(idencoder.Config{Secret: "test-secret", MinLength: 10}, noopLogger{})
	if err != nil {
		t.Fatalf("failed to create encoder: %v", err)
	}

	newRouter := func(svc input.Service) *gin.Engine {
		inter := interactor.NewInteractor(svc)
		h := New(nil, inter, enc, resp, nil, nil, nil, nil, nil, nil, nil)

		r := gin.New()
		r.Use(middleware.RequestID())
		r.Use(errHandler.Handle())
		r.POST("/flighthours/api/v1/register", h.RegisterEmployee())
		return r
	}

	t.Run("success", func(t *testing.T) {
		r := newRouter(&fakeService{})
		body := map[string]any{
			"name":                 "John",
			"airline":              "ACME",
			"email":                "john@example.com",
			"password":             "pw",
			"identificationNumber": "123",
			"bp":                   "bp",
			"start_date":           "2025-01-01",
			"end_date":             "",
			"active":               true,
			"role":                 "user",
		}
		b, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/flighthours/api/v1/register", bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		if w.Code != http.StatusCreated {
			t.Fatalf("expected status %d, got %d. body=%s", http.StatusCreated, w.Code, w.Body.String())
		}

		var out middleware.APIResponse
		if err := json.Unmarshal(w.Body.Bytes(), &out); err != nil {
			t.Fatalf("invalid json response: %v; body=%s", err, w.Body.String())
		}
		if !out.Success {
			t.Fatalf("expected success=true, got false")
		}
		if out.Code != domain.MsgUserRegistered {
			t.Fatalf("expected code %q, got %q", domain.MsgUserRegistered, out.Code)
		}
		if w.Header().Get(middleware.RequestIDHeader) == "" {
			t.Fatalf("expected %s response header to be set", middleware.RequestIDHeader)
		}
	})

	t.Run("invalid json => 400", func(t *testing.T) {
		r := newRouter(&fakeService{})
		req := httptest.NewRequest(http.MethodPost, "/flighthours/api/v1/register", bytes.NewBufferString("{"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected status %d, got %d. body=%s", http.StatusBadRequest, w.Code, w.Body.String())
		}
	})

	t.Run("duplicate user => 409", func(t *testing.T) {
		r := newRouter(&fakeServiceErr{err: domain.ErrDuplicateUser})
		body := map[string]any{
			"name":                 "John",
			"airline":              "ACME",
			"email":                "john@example.com",
			"password":             "pw",
			"identificationNumber": "123",
			"bp":                   "bp",
			"start_date":           "2025-01-01",
			"end_date":             "",
			"active":               true,
			"role":                 "user",
		}
		b, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/flighthours/api/v1/register", bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		if w.Code != http.StatusConflict {
			t.Fatalf("expected status %d, got %d. body=%s", http.StatusConflict, w.Code, w.Body.String())
		}
	})

	t.Run("incomplete registration => 409", func(t *testing.T) {
		r := newRouter(&fakeServiceErr{err: domain.ErrIncompleteRegistration})
		body := map[string]any{
			"name":                 "John",
			"airline":              "ACME",
			"email":                "john@example.com",
			"password":             "pw",
			"identificationNumber": "123",
			"bp":                   "bp",
			"start_date":           "2025-01-01",
			"end_date":             "",
			"active":               true,
			"role":                 "user",
		}
		b, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/flighthours/api/v1/register", bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		if w.Code != http.StatusConflict {
			t.Fatalf("expected status %d, got %d. body=%s", http.StatusConflict, w.Code, w.Body.String())
		}
	})
}
