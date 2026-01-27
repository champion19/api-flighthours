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
func (f *fakeService) ChangePassword(context.Context, string, string, string) (string, error) {
	return "", nil
}
func (f *fakeService) DeleteEmployee(context.Context, string, string) error {
	return nil
}
func (f *fakeService) GetEmployeeByKeycloakID(context.Context, string) (*domain.Employee, error) {
	return nil, errors.New("not implemented")
}

func (f *fakeService) GetEmployeesByRole(context.Context, string) ([]domain.Employee, error) {
	return nil, errors.New("not implemented")
}

func (f *fakeService) UpdateEmployee(context.Context, output.Tx, domain.Employee) error {
	return nil
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

func (f *fakeServiceErr) UpdateEmployee(context.Context, output.Tx, domain.Employee) error {
	return nil
}

type fakeMessageCacheRepo struct {
	messages []cachetypes.CachedMessage
}

func (r fakeMessageCacheRepo) GetAllActiveForCache(context.Context) ([]cachetypes.CachedMessage, error) {
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
		{Code: domain.MsgUserUpdated, Type: cachetypes.TypeSuccess, Content: "user updated"},
		{Code: domain.MsgValStartDateAfterEndDate, Type: cachetypes.TypeError, Content: "start date after end date"},
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
		inter := interactor.NewInteractor(svc, noopLogger{})
		h := New(nil, inter, enc, resp, nil, nil, nil)

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

func TestHTTP_UpdateEmployee(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cache := newTestMessageCache(t)
	resp := middleware.NewResponseHandler(cache)
	errHandler := middleware.NewErrorHandler(cache)

	enc, err := idencoder.NewHashidsEncoder(idencoder.Config{Secret: "test-secret", MinLength: 10}, noopLogger{})
	if err != nil {
		t.Fatalf("failed to create encoder: %v", err)
	}

	existingEmployee := &domain.Employee{
		ID:                   "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
		Name:                 "John Doe",
		Email:                "john@example.com",
		Airline:              "ACME",
		Bp:                   "bp123",
		IdentificationNumber: "ID123",
		Role:                 "user",
		Active:               true,
	}

	newRouter := func(svc input.Service) *gin.Engine {
		inter := interactor.NewInteractor(svc, noopLogger{})
		h := New(nil, inter, enc, resp, nil, nil, nil)

		r := gin.New()
		r.Use(middleware.RequestID())
		r.Use(errHandler.Handle())
		r.Use(func(c *gin.Context) {
			c.Set("authenticated_user", existingEmployee)
			c.Next()
		})
		r.PUT("/flighthours/api/v1/employees/me", h.UpdateEmployee())
		return r
	}

	t.Run("success", func(t *testing.T) {
		r := newRouter(&fakeService{})
		body := map[string]any{
			"name":                 "John Updated",
			"identificationNumber": "ID456",
			"start_date":           "2025-01-01",
			"end_date":             "2025-12-31",
			"active":               true,
			"role":                 "admin",
		}
		b, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPut, "/flighthours/api/v1/employees/me", bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expected status %d, got %d. body=%s", http.StatusOK, w.Code, w.Body.String())
		}
	})

	t.Run("invalid json => 400", func(t *testing.T) {
		r := newRouter(&fakeService{})
		req := httptest.NewRequest(http.MethodPut, "/flighthours/api/v1/employees/me", bytes.NewBufferString("{"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected status %d, got %d. body=%s", http.StatusBadRequest, w.Code, w.Body.String())
		}
	})

	t.Run("start_date after end_date => 400", func(t *testing.T) {
		r := newRouter(&fakeService{})
		body := map[string]any{
			"name":                 "John Updated",
			"identificationNumber": "ID456",
			"start_date":           "2026-01-01",
			"end_date":             "2025-12-31",
			"active":               true,
			"role":                 "admin",
		}
		b, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPut, "/flighthours/api/v1/employees/me", bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected status %d, got %d. body=%s", http.StatusBadRequest, w.Code, w.Body.String())
		}
	})
}

func TestHTTP_DeleteEmployee(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cache := newTestMessageCacheWithDelete(t)
	resp := middleware.NewResponseHandler(cache)
	errHandler := middleware.NewErrorHandler(cache)

	enc, err := idencoder.NewHashidsEncoder(idencoder.Config{Secret: "test-secret", MinLength: 10}, noopLogger{})
	if err != nil {
		t.Fatalf("failed to create encoder: %v", err)
	}

	existingEmployee := &domain.Employee{
		ID:                   "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
		Name:                 "John Doe",
		Email:                "john@example.com",
		Airline:              "ACME",
		Bp:                   "bp123",
		IdentificationNumber: "ID123",
		Role:                 "user",
		Active:               true,
		KeycloakUserID:       "kc-user-123",
	}

	newRouterWithAuth := func(svc input.Service) *gin.Engine {
		inter := interactor.NewInteractor(svc, noopLogger{})
		h := New(nil, inter, enc, resp, nil, nil, nil)

		r := gin.New()
		r.Use(middleware.RequestID())
		r.Use(errHandler.Handle())
		r.Use(func(c *gin.Context) {
			c.Set("authenticated_user", existingEmployee)
			c.Next()
		})
		r.DELETE("/flighthours/api/v1/employees/me", h.DeleteEmployee())
		return r
	}

	newRouterWithoutAuth := func(svc input.Service) *gin.Engine {
		inter := interactor.NewInteractor(svc, noopLogger{})
		h := New(nil, inter, enc, resp, nil, nil, nil)

		r := gin.New()
		r.Use(middleware.RequestID())
		r.Use(errHandler.Handle())
		r.DELETE("/flighthours/api/v1/employees/me", h.DeleteEmployee())
		return r
	}

	t.Run("success", func(t *testing.T) {
		r := newRouterWithAuth(&fakeServiceWithDelete{})
		req := httptest.NewRequest(http.MethodDelete, "/flighthours/api/v1/employees/me", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expected status %d, got %d. body=%s", http.StatusOK, w.Code, w.Body.String())
		}

		var out middleware.APIResponse
		if err := json.Unmarshal(w.Body.Bytes(), &out); err != nil {
			t.Fatalf("invalid json response: %v; body=%s", err, w.Body.String())
		}
		if !out.Success {
			t.Fatalf("expected success=true, got false")
		}
		if out.Code != domain.MsgUserDeleted {
			t.Fatalf("expected code %q, got %q", domain.MsgUserDeleted, out.Code)
		}
	})

	t.Run("not authenticated => 404", func(t *testing.T) {
		r := newRouterWithoutAuth(&fakeServiceWithDelete{})
		req := httptest.NewRequest(http.MethodDelete, "/flighthours/api/v1/employees/me", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		if w.Code != http.StatusNotFound {
			t.Fatalf("expected status %d, got %d. body=%s", http.StatusNotFound, w.Code, w.Body.String())
		}
	})

	t.Run("delete error => 500", func(t *testing.T) {
		r := newRouterWithAuth(&fakeServiceWithDeleteErr{err: domain.ErrUserCannotDelete})
		req := httptest.NewRequest(http.MethodDelete, "/flighthours/api/v1/employees/me", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected status %d, got %d. body=%s", http.StatusInternalServerError, w.Code, w.Body.String())
		}
	})

	t.Run("employee not found => 404", func(t *testing.T) {
		r := newRouterWithAuth(&fakeServiceWithDeleteErr{err: domain.ErrPersonNotFound})
		req := httptest.NewRequest(http.MethodDelete, "/flighthours/api/v1/employees/me", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		if w.Code != http.StatusNotFound {
			t.Fatalf("expected status %d, got %d. body=%s", http.StatusNotFound, w.Code, w.Body.String())
		}
	})
}

type fakeServiceWithDelete struct {
	fakeService
}

func (f *fakeServiceWithDelete) GetEmployeeByID(ctx context.Context, id string) (*domain.Employee, error) {
	return &domain.Employee{
		ID:             id,
		Email:          "test@example.com",
		KeycloakUserID: "kc-123",
	}, nil
}

func (f *fakeServiceWithDelete) DeleteEmployee(context.Context, string, string) error {
	return nil
}

type fakeServiceWithDeleteErr struct {
	fakeServiceErr
	err error
}

func (f *fakeServiceWithDeleteErr) GetEmployeeByID(ctx context.Context, id string) (*domain.Employee, error) {
	if f.err == domain.ErrPersonNotFound {
		return nil, domain.ErrPersonNotFound
	}
	return &domain.Employee{
		ID:             id,
		Email:          "test@example.com",
		KeycloakUserID: "kc-123",
	}, nil
}

func (f *fakeServiceWithDeleteErr) DeleteEmployee(context.Context, string, string) error {
	return f.err
}

func newTestMessageCacheWithDelete(t *testing.T) *messaging.MessageCache {
	t.Helper()

	repo := fakeMessageCacheRepo{messages: []cachetypes.CachedMessage{
		{Code: domain.MsgUserRegistered, Type: cachetypes.TypeSuccess, Content: "user registered"},
		{Code: domain.MsgValJSONInvalid, Type: cachetypes.TypeError, Content: "invalid json"},
		{Code: domain.MsgUserDuplicate, Type: cachetypes.TypeError, Content: "duplicate"},
		{Code: domain.MsgIncompleteRegistration, Type: cachetypes.TypeError, Content: "incomplete"},
		{Code: domain.MsgUserUpdated, Type: cachetypes.TypeSuccess, Content: "user updated"},
		{Code: domain.MsgValStartDateAfterEndDate, Type: cachetypes.TypeError, Content: "start date after end date"},
		{Code: domain.MsgUserDeleted, Type: cachetypes.TypeSuccess, Content: "user deleted"},
		{Code: domain.MsgUserNotFound, Type: cachetypes.TypeError, Content: "user not found"},
		{Code: domain.MsgUserCannotDelete, Type: cachetypes.TypeError, Content: "cannot delete user"},
		{Code: domain.MsgPersonNotFound, Type: cachetypes.TypeError, Content: "person not found"},
	}}
	cache := messaging.NewMessageCache(repo, 0)
	if err := cache.LoadMessages(context.Background()); err != nil {
		t.Fatalf("failed to load message cache: %v", err)
	}
	return cache
}

// ============================================
// HU19 - PASSWORD RESET TESTS
// ============================================

// fakeServicePasswordReset is a fake service for password reset tests
type fakeServicePasswordReset struct {
	fakeService
	sendPasswordResetErr error
	updatePasswordErr    error
	updatePasswordEmail  string
}

func (f *fakeServicePasswordReset) SendPasswordResetEmail(ctx context.Context, email string) error {
	return f.sendPasswordResetErr
}

func (f *fakeServicePasswordReset) UpdatePassword(ctx context.Context, token, password string) (string, error) {
	if f.updatePasswordErr != nil {
		return "", f.updatePasswordErr
	}
	return f.updatePasswordEmail, nil
}

func newTestMessageCacheWithPasswordReset(t *testing.T) *messaging.MessageCache {
	t.Helper()

	repo := fakeMessageCacheRepo{messages: []cachetypes.CachedMessage{
		{Code: domain.MsgKCPwdResetSent, Type: cachetypes.TypeSuccess, Content: "password reset email sent"},
		{Code: domain.MsgKCPwdUpdated, Type: cachetypes.TypeSuccess, Content: "password updated"},
		{Code: domain.MsgKCPwdMismatch, Type: cachetypes.TypeError, Content: "passwords do not match"},
		{Code: domain.MsgKCPwdUpdateTokenInvalid, Type: cachetypes.TypeError, Content: "invalid token"},
		{Code: domain.MsgKCPwdUpdateError, Type: cachetypes.TypeError, Content: "password update error"},
		{Code: domain.MsgValBadFormat, Type: cachetypes.TypeError, Content: "bad format"},
	}}
	cache := messaging.NewMessageCache(repo, 0)
	if err := cache.LoadMessages(context.Background()); err != nil {
		t.Fatalf("failed to load message cache: %v", err)
	}
	return cache
}

func TestHTTP_RequestPasswordReset(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cache := newTestMessageCacheWithPasswordReset(t)
	resp := middleware.NewResponseHandler(cache)
	errHandler := middleware.NewErrorHandler(cache)

	enc, err := idencoder.NewHashidsEncoder(idencoder.Config{Secret: "test-secret", MinLength: 10}, noopLogger{})
	if err != nil {
		t.Fatalf("failed to create encoder: %v", err)
	}

	newRouter := func(svc input.Service) *gin.Engine {
		inter := interactor.NewInteractor(svc, noopLogger{})
		h := New(nil, inter, enc, resp, nil, nil, nil)

		r := gin.New()
		r.Use(middleware.RequestID())
		r.Use(errHandler.Handle())
		r.POST("/flighthours/api/v1/auth/password-reset", h.PasswordReset())
		return r
	}

	t.Run("success", func(t *testing.T) {
		r := newRouter(&fakeServicePasswordReset{})
		body := map[string]any{
			"email": "test@example.com",
		}
		b, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/flighthours/api/v1/auth/password-reset", bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expected status %d, got %d. body=%s", http.StatusOK, w.Code, w.Body.String())
		}

		var out middleware.APIResponse
		if err := json.Unmarshal(w.Body.Bytes(), &out); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}
		if !out.Success {
			t.Errorf("expected success=true, got %v", out.Success)
		}
		if out.Code != domain.MsgKCPwdResetSent {
			t.Errorf("expected code %s, got %s", domain.MsgKCPwdResetSent, out.Code)
		}
	})

	t.Run("success_even_if_email_not_found (security pattern)", func(t *testing.T) {

		r := newRouter(&fakeServicePasswordReset{sendPasswordResetErr: domain.ErrUserNotFound})
		body := map[string]any{
			"email": "nonexistent@example.com",
		}
		b, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/flighthours/api/v1/auth/password-reset", bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expected status %d (security: always return success), got %d. body=%s", http.StatusOK, w.Code, w.Body.String())
		}
	})

	t.Run("invalid_json => 400", func(t *testing.T) {
		r := newRouter(&fakeServicePasswordReset{})
		req := httptest.NewRequest(http.MethodPost, "/flighthours/api/v1/auth/password-reset", bytes.NewReader([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected status %d, got %d. body=%s", http.StatusBadRequest, w.Code, w.Body.String())
		}
	})
}

func TestHTTP_UpdatePassword(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cache := newTestMessageCacheWithPasswordReset(t)
	resp := middleware.NewResponseHandler(cache)
	errHandler := middleware.NewErrorHandler(cache)

	enc, err := idencoder.NewHashidsEncoder(idencoder.Config{Secret: "test-secret", MinLength: 10}, noopLogger{})
	if err != nil {
		t.Fatalf("failed to create encoder: %v", err)
	}

	newRouter := func(svc input.Service) *gin.Engine {
		inter := interactor.NewInteractor(svc, noopLogger{})
		h := New(nil, inter, enc, resp, nil, nil, nil)

		r := gin.New()
		r.Use(middleware.RequestID())
		r.Use(errHandler.Handle())
		r.POST("/flighthours/api/v1/auth/update-password", h.UpdatePassword())
		return r
	}

	t.Run("success", func(t *testing.T) {
		r := newRouter(&fakeServicePasswordReset{updatePasswordEmail: "test@example.com"})
		body := map[string]any{
			"token":            "valid-jwt-token",
			"new_password":     "NewPassword123!",
			"confirm_password": "NewPassword123!",
		}
		b, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/flighthours/api/v1/auth/update-password", bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expected status %d, got %d. body=%s", http.StatusOK, w.Code, w.Body.String())
		}

		var out middleware.APIResponse
		if err := json.Unmarshal(w.Body.Bytes(), &out); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}
		if !out.Success {
			t.Errorf("expected success=true, got %v", out.Success)
		}
		if out.Code != domain.MsgKCPwdUpdated {
			t.Errorf("expected code %s, got %s", domain.MsgKCPwdUpdated, out.Code)
		}
	})

	t.Run("password_mismatch => 400", func(t *testing.T) {
		r := newRouter(&fakeServicePasswordReset{})
		body := map[string]any{
			"token":            "valid-jwt-token",
			"new_password":     "NewPassword123!",
			"confirm_password": "DifferentPassword123!",
		}
		b, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/flighthours/api/v1/auth/update-password", bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected status %d, got %d. body=%s", http.StatusBadRequest, w.Code, w.Body.String())
		}

		var out middleware.APIResponse
		if err := json.Unmarshal(w.Body.Bytes(), &out); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}
		if out.Code != domain.MsgKCPwdMismatch {
			t.Errorf("expected code %s, got %s", domain.MsgKCPwdMismatch, out.Code)
		}
	})

	t.Run("invalid_token => 401", func(t *testing.T) {
		r := newRouter(&fakeServicePasswordReset{updatePasswordErr: domain.ErrInvalidToken})
		body := map[string]any{
			"token":            "invalid-token",
			"new_password":     "NewPassword123!",
			"confirm_password": "NewPassword123!",
		}
		b, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/flighthours/api/v1/auth/update-password", bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		if w.Code != http.StatusUnauthorized {
			t.Fatalf("expected status %d, got %d. body=%s", http.StatusUnauthorized, w.Code, w.Body.String())
		}
	})

	t.Run("expired_token => 401", func(t *testing.T) {
		r := newRouter(&fakeServicePasswordReset{updatePasswordErr: domain.ErrTokenExpired})
		body := map[string]any{
			"token":            "expired-token",
			"new_password":     "NewPassword123!",
			"confirm_password": "NewPassword123!",
		}
		b, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/flighthours/api/v1/auth/update-password", bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		if w.Code != http.StatusUnauthorized {
			t.Fatalf("expected status %d, got %d. body=%s", http.StatusUnauthorized, w.Code, w.Body.String())
		}
	})

	t.Run("update_failed => 500", func(t *testing.T) {
		r := newRouter(&fakeServicePasswordReset{updatePasswordErr: domain.ErrPasswordUpdateFailed})
		body := map[string]any{
			"token":            "valid-token",
			"new_password":     "NewPassword123!",
			"confirm_password": "NewPassword123!",
		}
		b, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/flighthours/api/v1/auth/update-password", bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected status %d, got %d. body=%s", http.StatusInternalServerError, w.Code, w.Body.String())
		}
	})

	t.Run("invalid_json => 400", func(t *testing.T) {
		r := newRouter(&fakeServicePasswordReset{})
		req := httptest.NewRequest(http.MethodPost, "/flighthours/api/v1/auth/update-password", bytes.NewReader([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected status %d, got %d. body=%s", http.StatusBadRequest, w.Code, w.Body.String())
		}
	})
}
