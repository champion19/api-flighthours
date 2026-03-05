package cookie

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/champion19/api-flighthours/config"
	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func newTestManager(secure bool, sameSite string) *Manager {
	return NewManager(config.CookieConfig{
		Domain:   "",
		Secure:   secure,
		SameSite: sameSite,
	})
}

func setupTestContext() (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	return c, w
}

func TestSetAccessToken(t *testing.T) {
	m := newTestManager(false, "lax")
	c, w := setupTestContext()

	m.SetAccessToken(c, "test-access-token", 3600)

	cookies := w.Result().Cookies()
	found := false
	for _, ck := range cookies {
		if ck.Name == AccessTokenCookie {
			found = true
			if ck.Value != "test-access-token" {
				t.Errorf("expected value 'test-access-token', got '%s'", ck.Value)
			}
			if ck.Path != "/" {
				t.Errorf("expected path '/', got '%s'", ck.Path)
			}
			if !ck.HttpOnly {
				t.Error("expected HttpOnly=true")
			}
			if ck.MaxAge != 3600 {
				t.Errorf("expected MaxAge=3600, got %d", ck.MaxAge)
			}
		}
	}
	if !found {
		t.Errorf("cookie '%s' not found in response", AccessTokenCookie)
	}
}

func TestSetRefreshToken(t *testing.T) {
	m := newTestManager(false, "lax")
	c, w := setupTestContext()

	m.SetRefreshToken(c, "test-refresh-token", 86400)

	cookies := w.Result().Cookies()
	found := false
	for _, ck := range cookies {
		if ck.Name == RefreshTokenCookie {
			found = true
			if ck.Value != "test-refresh-token" {
				t.Errorf("expected value 'test-refresh-token', got '%s'", ck.Value)
			}
			if ck.Path != RefreshTokenPath {
				t.Errorf("expected path '%s', got '%s'", RefreshTokenPath, ck.Path)
			}
			if !ck.HttpOnly {
				t.Error("expected HttpOnly=true")
			}
		}
	}
	if !found {
		t.Errorf("cookie '%s' not found in response", RefreshTokenCookie)
	}
}

func TestSetTokens(t *testing.T) {
	m := newTestManager(false, "lax")
	c, w := setupTestContext()

	m.SetTokens(c, "access-123", 3600, "refresh-456", 86400)

	cookies := w.Result().Cookies()
	foundAccess, foundRefresh := false, false
	for _, ck := range cookies {
		if ck.Name == AccessTokenCookie {
			foundAccess = true
		}
		if ck.Name == RefreshTokenCookie {
			foundRefresh = true
		}
	}
	if !foundAccess {
		t.Error("access token cookie not found")
	}
	if !foundRefresh {
		t.Error("refresh token cookie not found")
	}
}

func TestClearTokens(t *testing.T) {
	m := newTestManager(false, "lax")
	c, w := setupTestContext()

	m.ClearTokens(c)

	cookies := w.Result().Cookies()
	for _, ck := range cookies {
		if ck.Name == AccessTokenCookie || ck.Name == RefreshTokenCookie {
			if ck.MaxAge != -1 {
				t.Errorf("expected MaxAge=-1 for %s, got %d", ck.Name, ck.MaxAge)
			}
		}
	}
}

func TestGetAccessToken(t *testing.T) {
	c, _ := setupTestContext()
	c.Request.AddCookie(&http.Cookie{Name: AccessTokenCookie, Value: "my-access-token"})

	token, err := GetAccessToken(c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if token != "my-access-token" {
		t.Errorf("expected 'my-access-token', got '%s'", token)
	}
}

func TestGetRefreshToken(t *testing.T) {
	c, _ := setupTestContext()
	c.Request.AddCookie(&http.Cookie{Name: RefreshTokenCookie, Value: "my-refresh-token"})

	token, err := GetRefreshToken(c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if token != "my-refresh-token" {
		t.Errorf("expected 'my-refresh-token', got '%s'", token)
	}
}

func TestGetAccessToken_NoCookie(t *testing.T) {
	c, _ := setupTestContext()

	_, err := GetAccessToken(c)
	if err == nil {
		t.Error("expected error when cookie is missing")
	}
}

func TestSecureMode(t *testing.T) {
	m := newTestManager(true, "strict")
	c, w := setupTestContext()

	m.SetAccessToken(c, "secure-token", 3600)

	cookies := w.Result().Cookies()
	for _, ck := range cookies {
		if ck.Name == AccessTokenCookie {
			if !ck.Secure {
				t.Error("expected Secure=true in production mode")
			}
		}
	}
}
