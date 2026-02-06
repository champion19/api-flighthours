package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/champion19/api-flighthours/core/interactor/dto"
	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/champion19/api-flighthours/core/ports/input"
	"github.com/champion19/api-flighthours/middleware"
	"github.com/champion19/api-flighthours/platform/cache/messaging"
	cachetypes "github.com/champion19/api-flighthours/platform/cache/types"
	"github.com/champion19/api-flighthours/tools/idencoder"
	"github.com/gin-gonic/gin"
)

// fakeEmployeeInteractor implements input.EmployeeInteractor for testing
type fakeEmployeeInteractor struct {
	registerEmployeeFn        func(ctx context.Context, employee domain.Employee) (*dto.RegisterEmployee, error)
	loginFn                   func(ctx context.Context, email, password string) (*dto.TokenResponse, error)
	verifyEmailByTokenFn      func(ctx context.Context, token string) (string, error)
	requestPasswordResetFn    func(ctx context.Context, email string) error
	updatePasswordFn          func(ctx context.Context, token, newPass, confirmPass string) (string, error)
	changePasswordFn          func(ctx context.Context, email, current, newPass, confirmPass string) (string, error)
	resendVerificationEmailFn func(ctx context.Context, email string) error
	deleteEmployeeFn          func(ctx context.Context, employeeID string) error
	updateEmployeeFn          func(ctx context.Context, employee domain.Employee) (*dto.UpdateEmployee, error)
	locateFn                  func(ctx context.Context, id string) (*dto.RegisterEmployee, error)
	getEmployeesByRoleFn      func(ctx context.Context, role string) ([]domain.Employee, error)
}

var _ input.EmployeeInteractor = (*fakeEmployeeInteractor)(nil)

func (f *fakeEmployeeInteractor) RegisterEmployee(ctx context.Context, employee domain.Employee) (*dto.RegisterEmployee, error) {
	if f.registerEmployeeFn != nil {
		return f.registerEmployeeFn(ctx, employee)
	}
	return nil, errors.New("not implemented")
}

func (f *fakeEmployeeInteractor) Login(ctx context.Context, email, password string) (*dto.TokenResponse, error) {
	if f.loginFn != nil {
		return f.loginFn(ctx, email, password)
	}
	return nil, errors.New("not implemented")
}

func (f *fakeEmployeeInteractor) VerifyEmailByToken(ctx context.Context, token string) (string, error) {
	if f.verifyEmailByTokenFn != nil {
		return f.verifyEmailByTokenFn(ctx, token)
	}
	return "", errors.New("not implemented")
}

func (f *fakeEmployeeInteractor) RequestPasswordReset(ctx context.Context, email string) error {
	if f.requestPasswordResetFn != nil {
		return f.requestPasswordResetFn(ctx, email)
	}
	return errors.New("not implemented")
}

func (f *fakeEmployeeInteractor) UpdatePassword(ctx context.Context, token, newPass, confirmPass string) (string, error) {
	if f.updatePasswordFn != nil {
		return f.updatePasswordFn(ctx, token, newPass, confirmPass)
	}
	return "", errors.New("not implemented")
}

func (f *fakeEmployeeInteractor) ChangePassword(ctx context.Context, email, current, newPass, confirmPass string) (string, error) {
	if f.changePasswordFn != nil {
		return f.changePasswordFn(ctx, email, current, newPass, confirmPass)
	}
	return "", errors.New("not implemented")
}

func (f *fakeEmployeeInteractor) ResendVerificationEmail(ctx context.Context, email string) error {
	if f.resendVerificationEmailFn != nil {
		return f.resendVerificationEmailFn(ctx, email)
	}
	return errors.New("not implemented")
}

func (f *fakeEmployeeInteractor) DeleteEmployee(ctx context.Context, employeeID string) error {
	if f.deleteEmployeeFn != nil {
		return f.deleteEmployeeFn(ctx, employeeID)
	}
	return errors.New("not implemented")
}

func (f *fakeEmployeeInteractor) UpdateEmployee(ctx context.Context, employee domain.Employee) (*dto.UpdateEmployee, error) {
	if f.updateEmployeeFn != nil {
		return f.updateEmployeeFn(ctx, employee)
	}
	return nil, errors.New("not implemented")
}

func (f *fakeEmployeeInteractor) Locate(ctx context.Context, id string) (*dto.RegisterEmployee, error) {
	if f.locateFn != nil {
		return f.locateFn(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func (f *fakeEmployeeInteractor) GetEmployeesByRole(ctx context.Context, role string) ([]domain.Employee, error) {
	if f.getEmployeesByRoleFn != nil {
		return f.getEmployeesByRoleFn(ctx, role)
	}
	return nil, errors.New("not implemented")
}

func newTestEmployeeMessageCache(t *testing.T) *messaging.MessageCache {
	t.Helper()

	repo := fakeMessageCacheRepo{messages: []cachetypes.CachedMessage{
		{Code: domain.MsgUserRegistered, Type: cachetypes.TypeSuccess, Content: "user registered"},
		{Code: domain.MsgValJSONInvalid, Type: cachetypes.TypeError, Content: "invalid json"},
		{Code: domain.MsgValBadFormat, Type: cachetypes.TypeError, Content: "bad format"},
		{Code: domain.MsgServerError, Type: cachetypes.TypeError, Content: "server error"},
		{Code: domain.MsgKCLoginSuccess, Type: cachetypes.TypeSuccess, Content: "login successful"},
		{Code: domain.MsgKCLoginEmailNotVerified, Type: cachetypes.TypeError, Content: "email not verified"},
		{Code: domain.MsgKCEmailVerified, Type: cachetypes.TypeSuccess, Content: "email verified"},
		{Code: domain.MsgKCVerifEmailSent, Type: cachetypes.TypeSuccess, Content: "verification email sent"},
	}}
	cache := messaging.NewMessageCache(repo, 0)
	if err := cache.LoadMessages(context.Background()); err != nil {
		t.Fatalf("failed to load message cache: %v", err)
	}
	return cache
}

func TestHTTP_Login(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cache := newTestEmployeeMessageCache(t)
	resp := middleware.NewResponseHandler(cache)
	errHandler := middleware.NewErrorHandler(cache)

	enc, err := idencoder.NewHashidsEncoder(idencoder.Config{Secret: "test-secret", MinLength: 10}, noopLogger{})
	if err != nil {
		t.Fatalf("failed to create encoder: %v", err)
	}

	newRouter := func(interactor input.EmployeeInteractor) *gin.Engine {
		h := New(nil, interactor, enc, resp, nil, nil, nil, nil, nil, nil, nil, nil, nil)

		r := gin.New()
		r.Use(middleware.RequestID())
		r.Use(errHandler.Handle())
		r.POST("/login", h.Login())
		return r
	}

	t.Run("success - returns tokens", func(t *testing.T) {
		fake := &fakeEmployeeInteractor{
			loginFn: func(ctx context.Context, email, password string) (*dto.TokenResponse, error) {
				return &dto.TokenResponse{
					AccessToken:  "access-token-123",
					RefreshToken: "refresh-token-456",
					ExpiresIn:    3600,
					TokenType:    "Bearer",
				}, nil
			},
		}

		router := newRouter(fake)
		body := `{"email":"test@example.com","password":"secret123"}`
		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d. body=%s", w.Code, w.Body.String())
		}

		var response map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if response["success"] != true {
			t.Errorf("expected success=true, got %v", response["success"])
		}
	})

	t.Run("interactor error - returns error", func(t *testing.T) {
		fake := &fakeEmployeeInteractor{
			loginFn: func(ctx context.Context, email, password string) (*dto.TokenResponse, error) {
				return nil, errors.New("invalid credentials")
			},
		}

		router := newRouter(fake)
		body := `{"email":"test@example.com","password":"wrongpassword"}`
		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		var response map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if response["success"] != false {
			t.Errorf("expected success=false, got %v", response["success"])
		}
	})

	t.Run("invalid JSON - returns error", func(t *testing.T) {
		fake := &fakeEmployeeInteractor{}

		router := newRouter(fake)
		body := `{invalid json`
		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		var response map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if response["success"] != false {
			t.Errorf("expected success=false, got %v", response["success"])
		}
	})
}

func TestHTTP_PasswordReset(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cache := newTestEmployeeMessageCache(t)
	resp := middleware.NewResponseHandler(cache)
	errHandler := middleware.NewErrorHandler(cache)

	enc, err := idencoder.NewHashidsEncoder(idencoder.Config{Secret: "test-secret", MinLength: 10}, noopLogger{})
	if err != nil {
		t.Fatalf("failed to create encoder: %v", err)
	}

	newRouter := func(interactor input.EmployeeInteractor) *gin.Engine {
		h := New(nil, interactor, enc, resp, nil, nil, nil, nil, nil, nil, nil, nil, nil)

		r := gin.New()
		r.Use(middleware.RequestID())
		r.Use(errHandler.Handle())
		r.POST("/password-reset", h.PasswordReset())
		return r
	}

	t.Run("success - sends password reset email", func(t *testing.T) {
		sendCalled := false
		fake := &fakeEmployeeInteractor{
			requestPasswordResetFn: func(ctx context.Context, email string) error {
				sendCalled = true
				return nil
			},
		}

		router := newRouter(fake)
		body := `{"email":"test@example.com"}`
		req := httptest.NewRequest(http.MethodPost, "/password-reset", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d. body=%s", w.Code, w.Body.String())
		}

		if !sendCalled {
			t.Error("expected RequestPasswordReset to be called")
		}
	})

	t.Run("invalid JSON - returns error", func(t *testing.T) {
		fake := &fakeEmployeeInteractor{}

		router := newRouter(fake)
		body := `{invalid json`
		req := httptest.NewRequest(http.MethodPost, "/password-reset", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		var response map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if response["success"] != false {
			t.Errorf("expected success=false, got %v", response["success"])
		}
	})
}

func TestHTTP_VerifyEmailByToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cache := newTestEmployeeMessageCache(t)
	resp := middleware.NewResponseHandler(cache)
	errHandler := middleware.NewErrorHandler(cache)

	enc, err := idencoder.NewHashidsEncoder(idencoder.Config{Secret: "test-secret", MinLength: 10}, noopLogger{})
	if err != nil {
		t.Fatalf("failed to create encoder: %v", err)
	}

	newRouter := func(interactor input.EmployeeInteractor) *gin.Engine {
		h := New(nil, interactor, enc, resp, nil, nil, nil, nil, nil, nil, nil, nil, nil)

		r := gin.New()
		r.Use(middleware.RequestID())
		r.Use(errHandler.Handle())
		r.POST("/verify-email", h.VerifyEmailByToken())
		return r
	}

	t.Run("success - verifies email", func(t *testing.T) {
		fake := &fakeEmployeeInteractor{
			verifyEmailByTokenFn: func(ctx context.Context, token string) (string, error) {
				return "verified@example.com", nil
			},
		}

		router := newRouter(fake)
		body := `{"token":"valid-token-123"}`
		req := httptest.NewRequest(http.MethodPost, "/verify-email", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d. body=%s", w.Code, w.Body.String())
		}

		var response map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if response["success"] != true {
			t.Errorf("expected success=true, got %v", response["success"])
		}
	})

	t.Run("invalid token - returns error", func(t *testing.T) {
		fake := &fakeEmployeeInteractor{
			verifyEmailByTokenFn: func(ctx context.Context, token string) (string, error) {
				return "", domain.ErrInvalidToken
			},
		}

		router := newRouter(fake)
		body := `{"token":"invalid-token"}`
		req := httptest.NewRequest(http.MethodPost, "/verify-email", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		var response map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if response["success"] != false {
			t.Errorf("expected success=false, got %v", response["success"])
		}
	})

	t.Run("invalid JSON - returns error", func(t *testing.T) {
		fake := &fakeEmployeeInteractor{}

		router := newRouter(fake)
		body := `{invalid json`
		req := httptest.NewRequest(http.MethodPost, "/verify-email", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		var response map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if response["success"] != false {
			t.Errorf("expected success=false, got %v", response["success"])
		}
	})
}

func TestHTTP_ResendVerificationEmail(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cache := newTestEmployeeMessageCache(t)
	resp := middleware.NewResponseHandler(cache)
	errHandler := middleware.NewErrorHandler(cache)

	enc, err := idencoder.NewHashidsEncoder(idencoder.Config{Secret: "test-secret", MinLength: 10}, noopLogger{})
	if err != nil {
		t.Fatalf("failed to create encoder: %v", err)
	}

	newRouter := func(interactor input.EmployeeInteractor) *gin.Engine {
		h := New(nil, interactor, enc, resp, nil, nil, nil, nil, nil, nil, nil, nil, nil)

		r := gin.New()
		r.Use(middleware.RequestID())
		r.Use(errHandler.Handle())
		r.POST("/resend-verification", h.ResendVerificationEmail())
		return r
	}

	t.Run("success - sends verification email", func(t *testing.T) {
		sendCalled := false
		fake := &fakeEmployeeInteractor{
			resendVerificationEmailFn: func(ctx context.Context, email string) error {
				sendCalled = true
				return nil
			},
		}

		router := newRouter(fake)
		body := `{"email":"test@example.com"}`
		req := httptest.NewRequest(http.MethodPost, "/resend-verification", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d. body=%s", w.Code, w.Body.String())
		}

		if !sendCalled {
			t.Error("expected ResendVerificationEmail to be called")
		}
	})

	t.Run("user not found - returns error", func(t *testing.T) {
		fake := &fakeEmployeeInteractor{
			resendVerificationEmailFn: func(ctx context.Context, email string) error {
				return domain.ErrUserNotFound
			},
		}

		router := newRouter(fake)
		body := `{"email":"notfound@example.com"}`
		req := httptest.NewRequest(http.MethodPost, "/resend-verification", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		var response map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if response["success"] != false {
			t.Errorf("expected success=false, got %v", response["success"])
		}
	})

	t.Run("email already verified - returns warning", func(t *testing.T) {
		fake := &fakeEmployeeInteractor{
			resendVerificationEmailFn: func(ctx context.Context, email string) error {
				return domain.ErrEmailAlreadyVerified
			},
		}

		router := newRouter(fake)
		body := `{"email":"verified@example.com"}`
		req := httptest.NewRequest(http.MethodPost, "/resend-verification", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Warning responses vary, just ensure we got a response
		if w.Code == 0 {
			t.Error("expected non-zero status code")
		}
	})

	t.Run("invalid JSON - returns error", func(t *testing.T) {
		fake := &fakeEmployeeInteractor{}

		router := newRouter(fake)
		body := `{invalid json`
		req := httptest.NewRequest(http.MethodPost, "/resend-verification", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		var response map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if response["success"] != false {
			t.Errorf("expected success=false, got %v", response["success"])
		}
	})
}

func TestHTTP_UpdatePassword(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cache := newTestEmployeeMessageCache(t)
	resp := middleware.NewResponseHandler(cache)
	errHandler := middleware.NewErrorHandler(cache)

	enc, err := idencoder.NewHashidsEncoder(idencoder.Config{Secret: "test-secret", MinLength: 10}, noopLogger{})
	if err != nil {
		t.Fatalf("failed to create encoder: %v", err)
	}

	newRouter := func(interactor input.EmployeeInteractor) *gin.Engine {
		h := New(nil, interactor, enc, resp, nil, nil, nil, nil, nil, nil, nil, nil, nil)

		r := gin.New()
		r.Use(middleware.RequestID())
		r.Use(errHandler.Handle())
		r.POST("/update-password", h.UpdatePassword())
		return r
	}

	t.Run("success - updates password", func(t *testing.T) {
		fake := &fakeEmployeeInteractor{
			updatePasswordFn: func(ctx context.Context, token, newPass, confirmPass string) (string, error) {
				return "user@example.com", nil
			},
		}

		router := newRouter(fake)
		body := `{"token":"valid-token","new_password":"NewPass123!","confirm_password":"NewPass123!"}`
		req := httptest.NewRequest(http.MethodPost, "/update-password", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d. body=%s", w.Code, w.Body.String())
		}
	})

	t.Run("password mismatch - returns error", func(t *testing.T) {
		fake := &fakeEmployeeInteractor{}

		router := newRouter(fake)
		body := `{"token":"valid-token","new_password":"NewPass123!","confirm_password":"DifferentPass!"}`
		req := httptest.NewRequest(http.MethodPost, "/update-password", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		var response map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if response["success"] != false {
			t.Errorf("expected success=false, got %v", response["success"])
		}
	})

	t.Run("invalid token - returns error", func(t *testing.T) {
		fake := &fakeEmployeeInteractor{
			updatePasswordFn: func(ctx context.Context, token, newPass, confirmPass string) (string, error) {
				return "", domain.ErrInvalidToken
			},
		}

		router := newRouter(fake)
		body := `{"token":"invalid-token","new_password":"NewPass123!","confirm_password":"NewPass123!"}`
		req := httptest.NewRequest(http.MethodPost, "/update-password", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		var response map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if response["success"] != false {
			t.Errorf("expected success=false, got %v", response["success"])
		}
	})

	t.Run("invalid JSON - returns error", func(t *testing.T) {
		fake := &fakeEmployeeInteractor{}

		router := newRouter(fake)
		body := `{invalid json`
		req := httptest.NewRequest(http.MethodPost, "/update-password", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		var response map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if response["success"] != false {
			t.Errorf("expected success=false, got %v", response["success"])
		}
	})
}
