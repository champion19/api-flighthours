package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/champion19/api-flighthours/core/interactor/services/domain"
	"github.com/gin-gonic/gin"
)

func TestGetAuthenticatedUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("returns user when exists", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		expectedUser := &domain.Employee{ID: "user-123", Name: "Test User"}
		c.Set("authenticated_user", expectedUser)

		user, exists := GetAuthenticatedUser(c)

		if !exists {
			t.Error("expected user to exist")
		}
		if user.ID != "user-123" {
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
