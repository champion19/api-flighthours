package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	jwtPkg "github.com/champion19/api-flighthours/platform/jwt"
	"github.com/gin-gonic/gin"
)

const (
	testUserID       = "user-123"
	errMsgUnexpected = "unexpected error: %v"
)

func TestGetAuthenticatedUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("returns user when exists", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		expectedUser := &domain.Employee{ID: testUserID, Name: "Test User"}
		c.Set("authenticated_user", expectedUser)

		user, exists := GetAuthenticatedUser(c)

		if !exists {
			t.Error("expected user to exist")
		}
		if user.ID != testUserID {
			t.Errorf("expected 'user-123', got %q", user.ID)
		}
	})

	t.Run("returns false when user not set", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())

		_, exists := GetAuthenticatedUser(c)

		if exists {
			t.Error("expected user to not exist")
		}
	})

	t.Run("returns false when wrong type", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Set("authenticated_user", "not-an-employee")

		_, exists := GetAuthenticatedUser(c)

		if exists {
			t.Error("expected false for wrong type")
		}
	})
}

func TestRequireRole(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("allows user with matching role", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)

		// Set authenticated user with 'admin' role
		user := &domain.Employee{ID: "user-123", Role: "admin"}
		c.Set("authenticated_user", user)

		middleware := RequireRole("admin", "superadmin")
		middleware(c)

		if c.IsAborted() {
			t.Error("expected request to NOT be aborted for matching role")
		}
	})

	t.Run("allows user with any allowed role", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)

		// Set authenticated user with 'pilot' role
		user := &domain.Employee{ID: "user-456", Role: "pilot"}
		c.Set("authenticated_user", user)

		middleware := RequireRole("admin", "pilot")
		middleware(c)

		if c.IsAborted() {
			t.Error("expected request to NOT be aborted for matching role")
		}
	})

	t.Run("blocks user without matching role", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)

		// Set authenticated user with 'viewer' role
		user := &domain.Employee{ID: "user-789", Role: "viewer"}
		c.Set("authenticated_user", user)

		middleware := RequireRole("admin", "superadmin")
		middleware(c)

		if !c.IsAborted() {
			t.Error("expected request to be aborted for non-matching role")
		}
	})

	t.Run("blocks when user not authenticated", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)
		// No authenticated_user set

		middleware := RequireRole("admin")
		middleware(c)

		if !c.IsAborted() {
			t.Error("expected request to be aborted when no user")
		}
	})

	t.Run("allows pilot role specifically", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)

		user := &domain.Employee{ID: "pilot-001", Role: "pilot"}
		c.Set("authenticated_user", user)

		middleware := RequireRole("pilot")
		middleware(c)

		if c.IsAborted() {
			t.Error("expected pilot to be allowed")
		}
	})

	t.Run("blocks empty role", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)

		user := &domain.Employee{ID: "user-no-role", Role: ""}
		c.Set("authenticated_user", user)

		middleware := RequireRole("admin")
		middleware(c)

		if !c.IsAborted() {
			t.Error("expected empty role to be blocked")
		}
	})
}

// ═══════════════════════════════════════════
// Tests for extractBearerToken
// ═══════════════════════════════════════════

func TestExtractBearerToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("from Authorization header", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)
		c.Request.Header.Set("Authorization", "Bearer my-test-token")

		token, err := extractBearerToken(c)
		if err != nil {
			t.Fatalf(errMsgUnexpected, err)
		}
		if token != "my-test-token" {
			t.Errorf("expected 'my-test-token', got %q", token)
		}
	})

	t.Run("invalid Authorization format", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)
		c.Request.Header.Set("Authorization", "Basic abc123")

		_, err := extractBearerToken(c)
		if err == nil {
			t.Error("expected error for non-Bearer auth")
		}
	})

	t.Run("from cookie", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)
		c.Request.AddCookie(&http.Cookie{Name: "fh_access_token", Value: "cookie-token"})

		token, err := extractBearerToken(c)
		if err != nil {
			t.Fatalf(errMsgUnexpected, err)
		}
		if token != "cookie-token" {
			t.Errorf("expected 'cookie-token', got %q", token)
		}
	})

	t.Run("no token anywhere", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)

		_, err := extractBearerToken(c)
		if err == nil {
			t.Error("expected error when no token")
		}
	})
}

// ═══════════════════════════════════════════
// Tests for resolveTokenClaims + mapTokenError
// ═══════════════════════════════════════════

type fakeTokenValidator struct {
	validateFn func(token string) (map[string]interface{}, error)
}

func (f *fakeTokenValidator) ValidateToken(token string) (map[string]interface{}, error) {
	return f.validateFn(token)
}

// Close satisfies output.TokenValidator for test fakes.
func (f *fakeTokenValidator) Close() {
	// no-op: test fake does not manage real resources
}

func TestResolveTokenClaimsWithValidator(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		validator := &fakeTokenValidator{
			validateFn: func(token string) (map[string]interface{}, error) {
				return map[string]interface{}{"sub": testUserID}, nil
			},
		}
		claims, err := resolveTokenClaims("test-token", validator, nil)
		if err != nil {
			t.Fatalf(errMsgUnexpected, err)
		}
		if claims["sub"] != testUserID {
			t.Errorf("expected sub=user-123, got %v", claims["sub"])
		}
	})

	t.Run("validation error", func(t *testing.T) {
		validator := &fakeTokenValidator{
			validateFn: func(token string) (map[string]interface{}, error) {
				return nil, domain.ErrInvalidToken
			},
		}
		_, err := resolveTokenClaims("bad-token", validator, nil)
		if err == nil {
			t.Error("expected error")
		}
	})
}

func TestMapTokenError(t *testing.T) {
	t.Run("expired token", func(t *testing.T) {
		err := mapTokenError(jwtPkg.ErrTokenExpired)
		if err != domain.ErrTokenExpired {
			t.Errorf("expected ErrTokenExpired, got %v", err)
		}
	})

	t.Run("other error", func(t *testing.T) {
		err := mapTokenError(jwtPkg.ErrInvalidSignature)
		if err != domain.ErrInvalidToken {
			t.Errorf("expected ErrInvalidToken, got %v", err)
		}
	})
}
