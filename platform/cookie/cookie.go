package cookie

import (
	"net/http"

	"github.com/champion19/api-flighthours/config"
	"github.com/gin-gonic/gin"
)

const (
	AccessTokenCookie  = "fh_access_token"
	RefreshTokenCookie = "fh_refresh_token"
	RefreshTokenPath   = "/flighthours/api/v1/auth"
)

// Manager handles HttpOnly cookie operations for authentication tokens.
type Manager struct {
	domain   string
	secure   bool
	sameSite http.SameSite
}

// NewManager creates a cookie Manager from the application config.
func NewManager(cfg config.CookieConfig) *Manager {
	return &Manager{
		domain:   cfg.Domain,
		secure:   cfg.Secure,
		sameSite: cfg.GetSameSiteMode(),
	}
}

// SetAccessToken sets the access token as an HttpOnly cookie on the root path.
func (m *Manager) SetAccessToken(c *gin.Context, token string, maxAge int) {
	c.SetSameSite(m.sameSite)
	c.SetCookie(AccessTokenCookie, token, maxAge, "/", m.domain, m.secure, true)
}

// SetRefreshToken sets the refresh token as an HttpOnly cookie on the auth path only.
func (m *Manager) SetRefreshToken(c *gin.Context, token string, maxAge int) {
	c.SetSameSite(m.sameSite)
	c.SetCookie(RefreshTokenCookie, token, maxAge, RefreshTokenPath, m.domain, m.secure, true)
}

// SetTokens sets both access and refresh token cookies.
func (m *Manager) SetTokens(c *gin.Context, accessToken string, accessMaxAge int, refreshToken string, refreshMaxAge int) {
	m.SetAccessToken(c, accessToken, accessMaxAge)
	m.SetRefreshToken(c, refreshToken, refreshMaxAge)
}

// ClearTokens removes both token cookies by setting MaxAge to -1.
func (m *Manager) ClearTokens(c *gin.Context) {
	c.SetSameSite(m.sameSite)
	c.SetCookie(AccessTokenCookie, "", -1, "/", m.domain, m.secure, true)
	c.SetCookie(RefreshTokenCookie, "", -1, RefreshTokenPath, m.domain, m.secure, true)
}

// GetAccessToken reads the access token from the cookie.
func GetAccessToken(c *gin.Context) (string, error) {
	return c.Cookie(AccessTokenCookie)
}

// GetRefreshToken reads the refresh token from the cookie.
func GetRefreshToken(c *gin.Context) (string, error) {
	return c.Cookie(RefreshTokenCookie)
}
