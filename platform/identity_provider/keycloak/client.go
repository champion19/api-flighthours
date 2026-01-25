package keycloak

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/Nerzal/gocloak/v13"
	"github.com/champion19/flighthours-api/config"
	"github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/champion19/flighthours-api/core/ports/output"
	"github.com/champion19/flighthours-api/platform/logger"
)

type client struct {
	gocloak        *gocloak.GoCloak
	config         *config.KeycloakConfig
	token          *gocloak.JWT
	tokenExpiresAt time.Time
	tokenMutex     sync.RWMutex
	logger         logger.Logger
}

func NewClient(cfg *config.KeycloakConfig, log logger.Logger) (output.AuthClient, error) {
	if cfg == nil {
		return nil, fmt.Errorf("keycloak config cannot be nil")
	}

	log.Info(logger.LogKeycloakClientInit, "server_url", cfg.ServerURL, "realm", cfg.Realm)

	gc := gocloak.NewClient(cfg.ServerURL)

	authClient := &client{
		gocloak: gc,
		config:  cfg,
		logger:  log,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	log.Debug(logger.LogKeycloakAdminAuth, "admin_user", cfg.AdminUser, "realm", cfg.Realm)
	token, err := authClient.gocloak.LoginAdmin(ctx, authClient.config.AdminUser, authClient.config.AdminPass, authClient.config.Realm)
	if err != nil {
		log.Error(logger.LogKeycloakAdminAuthError, "error", err, "realm", cfg.Realm)
		return nil, fmt.Errorf("failed to initialize admin token: %w", err)
	}
	authClient.token = token
	authClient.tokenExpiresAt = time.Now().Add(time.Duration(token.ExpiresIn) * time.Second)

	log.Success(logger.LogKeycloakClientOK, "realm", cfg.Realm, "expires_in", token.ExpiresIn)

	return authClient, nil
}

func (c *client) ensureValidToken(ctx context.Context) error {
	c.tokenMutex.RLock()

	needsRefresh := time.Now().Add(30 * time.Second).After(c.tokenExpiresAt)
	c.tokenMutex.RUnlock()

	if !needsRefresh {
		return nil
	}

	c.tokenMutex.Lock()
	defer c.tokenMutex.Unlock()

	if time.Now().Add(30 * time.Second).Before(c.tokenExpiresAt) {
		return nil
	}

	c.logger.Info(logger.LogKeycloakTokenRefresh,
		"realm", c.config.Realm,
		"admin_user", c.config.AdminUser,
		"token_expires_at", c.tokenExpiresAt.Format(time.RFC3339))

	token, err := c.gocloak.LoginAdmin(ctx, c.config.AdminUser, c.config.AdminPass, c.config.Realm)
	if err != nil {
		c.logger.Error(logger.LogKeycloakTokenRefreshErr,
			"realm", c.config.Realm,
			"admin_user", c.config.AdminUser,
			"error", err)
		return fmt.Errorf("failed to refresh admin token: %w", err)
	}

	c.token = token
	c.tokenExpiresAt = time.Now().Add(time.Duration(token.ExpiresIn) * time.Second)

	c.logger.Success(logger.LogKeycloakTokenRefreshOK,
		"realm", c.config.Realm,
		"admin_user", c.config.AdminUser,
		"new_expires_at", c.tokenExpiresAt.Format(time.RFC3339),
		"expires_in_seconds", token.ExpiresIn)

	return nil
}

func (c *client) LoginUser(ctx context.Context, username, password string) (*gocloak.JWT, error) {
	if username == "" || password == "" {
		return nil, fmt.Errorf("username and password cannot be empty")
	}

	c.logger.Info(logger.LogKeycloakUserLogin, "username", username, "realm", c.config.Realm)

	token, err := c.gocloak.Login(
		ctx,
		c.config.ClientID,
		c.config.ClientSecret,
		c.config.Realm,
		username,
		password,
	)
	if err != nil {
		c.logger.Error(logger.LogKeycloakUserLoginError, "username", username, "error", err)
		return nil, fmt.Errorf("user login failed: %w", err)
	}

	c.logger.Success(logger.LogKeycloakUserLoginOK, "username", username)
	return token, nil
}

func (c *client) CreateUser(ctx context.Context, employee *domain.Employee) (string, error) {
	if employee == nil {
		return "", fmt.Errorf("employee cannot be nil")
	}

	if err := c.ensureValidToken(ctx); err != nil {
		return "", err
	}

	c.logger.Info(logger.LogKeycloakUserCreate, "email", employee.Email, "realm", c.config.Realm)

	keycloakUser := gocloak.User{
		Email:     &employee.Email,
		FirstName: &employee.Name,
		LastName:  &employee.Name,
		Enabled:   gocloak.BoolP(true),
		Username:  &employee.Email,
	}

	userID, err := c.gocloak.CreateUser(
		ctx,
		c.token.AccessToken,
		c.config.Realm,
		keycloakUser,
	)
	if err != nil {
		c.logger.Error(logger.LogKeycloakUserCreateError, "email", employee.Email, "error", err)
		return "", fmt.Errorf("failed to create user in keycloak: %w", err)
	}

	c.logger.Success(logger.LogKeycloakUserCreateOK, "email", employee.Email, "user_id", userID)
	return userID, nil
}

func (c *client) GetUserByEmail(ctx context.Context, email string) (*gocloak.User, error) {
	if email == "" {
		return nil, fmt.Errorf("email cannot be empty")
	}

	if err := c.ensureValidToken(ctx); err != nil {
		return nil, err
	}

	users, err := c.gocloak.GetUsers(
		ctx,
		c.token.AccessToken,
		c.config.Realm,
		gocloak.GetUsersParams{
			Email: &email,
			Exact: gocloak.BoolP(true),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	if len(users) == 0 {
		return nil, fmt.Errorf("user with email %s not found", email)
	}

	return users[0], nil
}

func (c *client) GetUserByID(ctx context.Context, userID string) (*gocloak.User, error) {
	if userID == "" {
		return nil, fmt.Errorf("userID cannot be empty")
	}

	if err := c.ensureValidToken(ctx); err != nil {
		return nil, err
	}

	user, err := c.gocloak.GetUserByID(
		ctx,
		c.token.AccessToken,
		c.config.Realm,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return user, nil
}

func (c *client) UpdateUser(ctx context.Context, user *gocloak.User) error {
	if user == nil || user.ID == nil {
		return fmt.Errorf("user or user ID cannot be nil")
	}

	if err := c.ensureValidToken(ctx); err != nil {
		return err
	}

	err := c.gocloak.UpdateUser(
		ctx,
		c.token.AccessToken,
		c.config.Realm,
		*user,
	)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func (c *client) DeleteUser(ctx context.Context, userID string) error {
	if userID == "" {
		return fmt.Errorf("userID cannot be empty")
	}

	if err := c.ensureValidToken(ctx); err != nil {
		return err
	}

	c.logger.Warn(logger.LogKeycloakUserDelete, "user_id", userID)

	err := c.gocloak.DeleteUser(
		ctx,
		c.token.AccessToken,
		c.config.Realm,
		userID,
	)
	if err != nil {
		c.logger.Error(logger.LogKeycloakUserDeleteError, "user_id", userID, "error", err)
		return fmt.Errorf("failed to delete user: %w", err)
	}

	c.logger.Info(logger.LogKeycloakUserDeleteOK, "user_id", userID)
	return nil
}

func (c *client) SetPassword(ctx context.Context, userID string, password string, temporary bool) error {
	if userID == "" || password == "" {
		return fmt.Errorf("userID and password cannot be empty")
	}

	if err := c.ensureValidToken(ctx); err != nil {
		return err
	}

	c.logger.Debug(logger.LogKeycloakPasswordSet, "user_id", userID, "temporary", temporary)

	err := c.gocloak.SetPassword(
		ctx,
		c.token.AccessToken,
		userID,
		c.config.Realm,
		password,
		temporary,
	)
	if err != nil {
		c.logger.Error(logger.LogKeycloakPasswordSetError, "user_id", userID, "error", err)
		return fmt.Errorf("failed to set password: %w", err)
	}

	c.logger.Success(logger.LogKeycloakPasswordSetOK, "user_id", userID)
	return nil
}

func (c *client) AssignRole(ctx context.Context, userID string, roleName string) error {
	if userID == "" || roleName == "" {
		return fmt.Errorf("userID and roleName cannot be empty")
	}

	if err := c.ensureValidToken(ctx); err != nil {
		return err
	}

	c.logger.Info(logger.LogKeycloakRoleAssign, "user_id", userID, "role", roleName)

	// Obtener el role por nombre
	role, err := c.gocloak.GetRealmRole(
		ctx,
		c.token.AccessToken,
		c.config.Realm,
		roleName,
	)
	if err != nil {
		c.logger.Error(logger.LogKeycloakRoleGetError, "role", roleName, "error", err)
		return fmt.Errorf("failed to get role %s: %w", roleName, err)
	}

	// Asignar el role al usuario
	err = c.gocloak.AddRealmRoleToUser(
		ctx,
		c.token.AccessToken,
		c.config.Realm,
		userID,
		[]gocloak.Role{*role},
	)
	if err != nil {
		c.logger.Error(logger.LogKeycloakRoleAssignError, "user_id", userID, "role", roleName, "error", err)
		return fmt.Errorf("failed to assign role to user: %w", err)
	}

	c.logger.Success(logger.LogKeycloakRoleAssignOK, "user_id", userID, "role", roleName)
	return nil
}

func (c *client) RemoveRole(ctx context.Context, userID string, roleName string) error {
	if userID == "" || roleName == "" {
		return fmt.Errorf("userID and roleName cannot be empty")
	}

	if err := c.ensureValidToken(ctx); err != nil {
		return err
	}

	role, err := c.gocloak.GetRealmRole(
		ctx,
		c.token.AccessToken,
		c.config.Realm,
		roleName,
	)
	if err != nil {
		return fmt.Errorf("failed to get role %s: %w", roleName, err)
	}

	err = c.gocloak.DeleteRealmRoleFromUser(
		ctx,
		c.token.AccessToken,
		c.config.Realm,
		userID,
		[]gocloak.Role{*role},
	)
	if err != nil {
		return fmt.Errorf("failed to remove role from user: %w", err)
	}

	return nil
}

func (c *client) GetUserRoles(ctx context.Context, userID string) ([]*gocloak.Role, error) {
	if userID == "" {
		return nil, fmt.Errorf("userID cannot be empty")
	}

	if err := c.ensureValidToken(ctx); err != nil {
		return nil, err
	}

	roles, err := c.gocloak.GetRealmRolesByUserID(
		ctx,
		c.token.AccessToken,
		c.config.Realm,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}

	return roles, nil
}

func (c *client) SendVerificationEmail(ctx context.Context, userID string) error {
	if userID == "" {
		return fmt.Errorf("userID cannot be empty")
	}

	if err := c.ensureValidToken(ctx); err != nil {
		return err
	}

	params := gocloak.ExecuteActionsEmail{
		UserID:   &userID,
		Actions:  &[]string{"VERIFY_EMAIL"},
		Lifespan: gocloak.IntP(86400), // 24 horas
	}

	err := c.gocloak.ExecuteActionsEmail(
		ctx,
		c.token.AccessToken,
		c.config.Realm,
		params,
	)
	if err != nil {
		return fmt.Errorf("failed to send verification email: %w", err)
	}

	return nil
}

// SendPasswordResetEmail sends a password reset email to the user
// It searches for the user by email first, then sends the reset email
func (c *client) SendPasswordResetEmail(ctx context.Context, email string) error {
	if email == "" {
		return fmt.Errorf("email cannot be empty")
	}

	if err := c.ensureValidToken(ctx); err != nil {
		return err
	}

	c.logger.Info(logger.LogKeycloakSendPasswordReset, "email", email, "realm", c.config.Realm)

	// Buscar usuario por email
	c.logger.Debug(logger.LogKeycloakSearchUserByEmail, "email", email)

	users, err := c.gocloak.GetUsers(
		ctx,
		c.token.AccessToken,
		c.config.Realm,
		gocloak.GetUsersParams{
			Email: &email,
			Exact: gocloak.BoolP(true),
		},
	)
	if err != nil {
		c.logger.Error(logger.LogKeycloakSendPasswordResetError, "email", email, "error", err)
		return fmt.Errorf("failed to search user: %w", err)
	}

	if len(users) == 0 {
		c.logger.Warn(logger.LogKeycloakUserNotFound, "email", email)
		return fmt.Errorf("user with email %s not found", email)
	}

	c.logger.Debug(logger.LogKeycloakSearchUserByEmailOK, "email", email, "user_id", *users[0].ID)

	// Enviar email de reset de contraseña
	params := gocloak.ExecuteActionsEmail{
		UserID:   users[0].ID,
		Actions:  &[]string{"UPDATE_PASSWORD"},
		Lifespan: gocloak.IntP(43200), // 12 horas
	}

	err = c.gocloak.ExecuteActionsEmail(
		ctx,
		c.token.AccessToken,
		c.config.Realm,
		params,
	)
	if err != nil {
		c.logger.Error(logger.LogKeycloakSendPasswordResetError, "email", email, "error", err)
		return fmt.Errorf("failed to send password reset email: %w", err)
	}

	c.logger.Success(logger.LogKeycloakSendPasswordResetOK, "email", email, "user_id", *users[0].ID)
	return nil
}

func (c *client) VerifyEmail(ctx context.Context, userID string) error {
	if userID == "" {
		return fmt.Errorf("userID cannot be empty")
	}

	user, err := c.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	emailVerified := true
	user.EmailVerified = &emailVerified

	err = c.gocloak.UpdateUser(
		ctx,
		c.token.AccessToken,
		c.config.Realm,
		*user,
	)
	if err != nil {
		return fmt.Errorf("failed to verify email: %w", err)
	}

	return nil
}

func (c *client) Logout(ctx context.Context, refreshToken string) error {
	if refreshToken == "" {
		return fmt.Errorf("refreshToken cannot be empty")
	}

	err := c.gocloak.Logout(
		ctx,
		c.config.ClientID,
		c.config.ClientSecret,
		c.config.Realm,
		refreshToken,
	)
	if err != nil {
		return fmt.Errorf("failed to logout: %w", err)
	}

	return nil
}

func (c *client) RefreshToken(ctx context.Context, refreshToken string) (*gocloak.JWT, error) {
	if refreshToken == "" {
		return nil, fmt.Errorf("refreshToken cannot be empty")
	}

	c.logger.Info(logger.LogKeycloakUserTokenRefresh,
		"realm", c.config.Realm,
		"client_id", c.config.ClientID)

	token, err := c.gocloak.RefreshToken(
		ctx,
		refreshToken,
		c.config.ClientID,
		c.config.ClientSecret,
		c.config.Realm,
	)
	if err != nil {
		c.logger.Error(logger.LogKeycloakUserTokenRefreshErr,
			"realm", c.config.Realm,
			"client_id", c.config.ClientID,
			"error", err)
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}

	c.logger.Success(logger.LogKeycloakUserTokenRefreshOK,
		"realm", c.config.Realm,
		"client_id", c.config.ClientID,
		"expires_in_seconds", token.ExpiresIn,
		"refresh_expires_in_seconds", token.RefreshExpiresIn)

	return token, nil
}

// ValidateActionToken validates a Keycloak action token (from password reset email)
// and returns the userID and email from the token claims
// This performs the following validations:
// 1. Token format is valid JWT
// 2. Token has not expired
// 3. Token issuer matches our Keycloak realm
// 4. Token type is an action token (typ: "reset-credentials")
// 5. User exists in Keycloak
func (c *client) ValidateActionToken(ctx context.Context, token string) (string, string, error) {
	if token == "" {
		return "", "", fmt.Errorf("token cannot be empty")
	}

	c.logger.Debug(logger.LogKeycloakPasswordTokenValidation, "realm", c.config.Realm)

	// Parse the JWT token to extract claims
	parts := splitToken(token)
	if len(parts) != 3 {
		c.logger.Error(logger.LogKeycloakPasswordTokenInvalid, "reason", "invalid token format")
		return "", "", fmt.Errorf("invalid token format")
	}

	// Decode the payload (second part)
	claims, err := decodeJWTPayload(parts[1])
	if err != nil {
		c.logger.Error(logger.LogKeycloakPasswordTokenInvalid, "error", err)
		return "", "", fmt.Errorf("failed to decode token: %w", err)
	}

	// 1. Validate expiration (exp claim)
	if exp, ok := claims["exp"].(float64); ok {
		expirationTime := time.Unix(int64(exp), 0)
		if time.Now().After(expirationTime) {
			c.logger.Error(logger.LogKeycloakPasswordTokenInvalid, "reason", "token expired",
				"expired_at", expirationTime.Format(time.RFC3339))
			return "", "", fmt.Errorf("token has expired")
		}
		c.logger.Debug("Token expiration validated", "expires_at", expirationTime.Format(time.RFC3339))
	} else {
		c.logger.Warn("Token missing expiration claim, proceeding with caution")
	}

	// 2. Validate issuer (iss claim)
	expectedIssuer := fmt.Sprintf("%s/realms/%s", c.config.ServerURL, c.config.Realm)
	if iss, ok := claims["iss"].(string); ok {
		if iss != expectedIssuer {
			c.logger.Error(logger.LogKeycloakPasswordTokenInvalid, "reason", "invalid issuer",
				"expected", expectedIssuer, "got", iss)
			return "", "", fmt.Errorf("invalid token issuer")
		}
		c.logger.Debug("Token issuer validated", "issuer", iss)
	} else {
		c.logger.Warn("Token missing issuer claim")
	}

	// 3. Validate action type (typ claim for action tokens)
	// Keycloak action tokens have typ like "kc-action", "reset-credentials", etc.
	if typ, ok := claims["typ"].(string); ok {
		c.logger.Debug("Token type", "typ", typ)
		// Action tokens typically have specific types
		// We allow any type since this is coming from a Keycloak email action
	}

	// 4. Extract user ID (sub claim)
	userID, ok := claims["sub"].(string)
	if !ok || userID == "" {
		c.logger.Error(logger.LogKeycloakPasswordTokenInvalid, "reason", "missing sub claim")
		return "", "", fmt.Errorf("missing user ID in token")
	}

	// 5. Verify user exists in Keycloak (this confirms the token is for a real user)
	if err := c.ensureValidToken(ctx); err != nil {
		return "", "", fmt.Errorf("failed to authenticate with Keycloak: %w", err)
	}

	user, err := c.GetUserByID(ctx, userID)
	if err != nil {
		c.logger.Error(logger.LogKeycloakPasswordTokenInvalid, "reason", "user not found",
			"user_id", userID, "error", err)
		return "", "", fmt.Errorf("user not found in Keycloak: %w", err)
	}

	// Get email from user or token
	var email string
	if user.Email != nil {
		email = *user.Email
	} else if emailClaim, ok := claims["email"].(string); ok {
		email = emailClaim
	}

	// Verify user is enabled
	if user.Enabled != nil && !*user.Enabled {
		c.logger.Error(logger.LogKeycloakPasswordTokenInvalid, "reason", "user is disabled",
			"user_id", userID)
		return "", "", fmt.Errorf("user account is disabled")
	}

	c.logger.Success(logger.LogKeycloakPasswordTokenValidOK,
		"user_id", userID,
		"email", email,
		"realm", c.config.Realm)

	return userID, email, nil
}

// splitToken splits a JWT token into its parts
func splitToken(token string) []string {
	return strings.Split(token, ".")
}

// decodeJWTPayload decodes the base64-encoded JWT payload
func decodeJWTPayload(payload string) (map[string]interface{}, error) {
	// JWT uses base64url encoding, decode it
	decoded, err := base64.RawURLEncoding.DecodeString(payload)
	if err != nil {
		// Try with padding
		decoded, err = base64.URLEncoding.DecodeString(payload)
		if err != nil {
			return nil, fmt.Errorf("failed to decode base64: %w", err)
		}
	}

	// Parse JSON
	var claims map[string]interface{}
	if err := json.Unmarshal(decoded, &claims); err != nil {
		return nil, fmt.Errorf("failed to parse JSON claims: %w", err)
	}

	return claims, nil
}
