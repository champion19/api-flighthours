package config

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/champion19/api-flighthours/platform/logger"
	"github.com/champion19/api-flighthours/tools/utils"
	"github.com/joho/godotenv"
)

type Config struct {
	Environment  string          `json:"environment"`
	Database     Database        `json:"database"`
	Server       Server          `json:"server"`
	Resend       Resend          `json:"resend"`
	Verification Verification    `json:"verification"`
	Keycloak     KeycloakConfig  `json:"keycloak"`
	IDEncoder    IDEncoderConfig `json:"id_encoder"`
	Cookie       CookieConfig    `json:"cookie"`
}

type CookieConfig struct {
	Domain   string `json:"domain"`
	Secure   bool   `json:"secure"`
	SameSite string `json:"same_site"`
}

func (c *CookieConfig) GetSameSiteMode() http.SameSite {
	switch strings.ToLower(c.SameSite) {
	case "strict":
		return http.SameSiteStrictMode
	case "none":
		return http.SameSiteNoneMode
	default:
		return http.SameSiteLaxMode
	}
}

type Verification struct {
	BaseURL string `json:"base_url"`
}

type Database struct {
	Driver          string `json:"driver"`
	Host            string `json:"host"`
	Port            string `json:"port"`
	Username        string `json:"username"`
	Password        string `json:"password"`
	Name            string `json:"name"`
	SSL             string `json:"ssl,omitempty"`
	MaxOpenConns    int    `json:"max_open_conns,omitempty"`
	MaxIdleConns    int    `json:"max_idle_conns,omitempty"`
	ConnMaxLifetime int    `json:"conn_max_lifetime,omitempty"`
	ConnMaxIdleTime int    `json:"conn_max_idle_time,omitempty"`
}

type Server struct {
	Port string `json:"port"`
	Host string `json:"host"`
}

type Resend struct {
	APIKey    string `json:"api_key"`
	FromEmail string `json:"from_email"`
}

type KeycloakConfig struct {
	ServerURL    string `json:"server_url"`
	Realm        string `json:"realm"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	AdminUser    string `json:"admin_user"`
	AdminPass    string `json:"admin_pass"`
}

type IDEncoderConfig struct {
	Secret    string `json:"secret"`
	MinLength int    `json:"min_length"`
}

func LoadConfig() (*Config, error) {
	root, err := utils.FindModuleRoot()
	if err != nil {
		return nil, fmt.Errorf("error finding module root: %w", err)
	}

	// Load .env file from project root (non-fatal if not found)
	envPath := filepath.Join(root, ".env")
	if err := godotenv.Load(envPath); err != nil {
		slog.Warn(logger.LogConfigEnvFileNotFound,
			slog.String("path", envPath))
	}

	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "local"
	}

	const localConfigFile = "local-config.json"

	var configFile string
	switch env {
	case "railway":
		configFile = "railway-config.json"
	default:
		configFile = localConfigFile
	}

	configPath := filepath.Join(root, "config", configFile)

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		slog.Warn(logger.LogConfigFileNotFound,
			slog.String("requested_file", configFile),
			slog.String("fallback_file", localConfigFile))
		configPath = filepath.Join(root, "config", localConfigFile)
	}

	file, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("error reading config file %s: %w", configPath, err)
	}

	var config Config
	if err = json.Unmarshal(file, &config); err != nil {
		return nil, fmt.Errorf("error parsing JSON configuration: %w", err)
	}

	if serverURL := os.Getenv("KEYCLOAK_SERVER_URL"); serverURL != "" {
		config.Keycloak.ServerURL = serverURL
	}
	if realm := os.Getenv("KEYCLOAK_REALM"); realm != "" {
		config.Keycloak.Realm = realm
	}
	if clientID := os.Getenv("KEYCLOAK_CLIENT_ID"); clientID != "" {
		config.Keycloak.ClientID = clientID
	}
	if clientSecret := os.Getenv("KEYCLOAK_CLIENT_SECRET"); clientSecret != "" {
		config.Keycloak.ClientSecret = clientSecret
	}
	if adminUser := os.Getenv("KEYCLOAK_ADMIN"); adminUser != "" {
		config.Keycloak.AdminUser = adminUser
	}
	if adminPass := os.Getenv("KEYCLOAK_ADMIN_PASSWORD"); adminPass != "" {
		config.Keycloak.AdminPass = adminPass
	}

	// Override Resend config from env vars
	if resendKey := os.Getenv("RESEND_API_KEY"); resendKey != "" {
		config.Resend.APIKey = resendKey
	}

	// Override ID encoder config from env vars
	if idSecret := os.Getenv("ID_ENCODER_SECRET"); idSecret != "" {
		config.IDEncoder.Secret = idSecret
	}

	// Override database password from env var
	if dbPass := os.Getenv("DB_PASSWORD"); dbPass != "" {
		config.Database.Password = dbPass
	}

	// Override cookie config from env vars
	if cookieDomain := os.Getenv("COOKIE_DOMAIN"); cookieDomain != "" {
		config.Cookie.Domain = cookieDomain
	}
	if cookieSecure := os.Getenv("COOKIE_SECURE"); cookieSecure == "true" {
		config.Cookie.Secure = true
	}

	slog.Info(logger.LogAppConfigLoaded,
		slog.String("config_file", configFile),
		slog.String("environment", config.Environment),
		slog.String("config_path", configPath),
		slog.String("keycloak_server", config.Keycloak.ServerURL),
		slog.String("keycloak_realm", config.Keycloak.Realm))

	return &config, nil
}

func (c *Config) GetMySQLDSN() string {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local",
		c.Database.Username,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.Name,
	)

	if c.Database.SSL != "" {
		dsn += "&tls=" + c.Database.SSL
	}

	return dsn
}

func (c *Config) GetServerAddress() string {
	return fmt.Sprintf("%s:%s", c.Server.Host, c.Server.Port)
}

func (c *Config) IsProduction() bool {
	return c.Environment == "production" || c.Environment == "railway"
}

func (c *Config) GetKeycloakAuthURL() string {
	return fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token",
		c.Keycloak.ServerURL,
		c.Keycloak.Realm)
}

func (c *Config) GetKeycloakAdminURL() string {
	return fmt.Sprintf("%s/admin/realms/%s",
		c.Keycloak.ServerURL,
		c.Keycloak.Realm)
}

func (c *Config) GetKeycloakJWKSURL() string {
	return fmt.Sprintf("%s/realms/%s/protocol/openid-connect/certs",
		c.Keycloak.ServerURL,
		c.Keycloak.Realm)
}

func (c *Config) GetKeycloakIssuerURL() string {
	return fmt.Sprintf("%s/realms/%s",
		c.Keycloak.ServerURL,
		c.Keycloak.Realm)
}
