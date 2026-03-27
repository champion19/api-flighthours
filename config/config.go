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
	IssuerURL    string `json:"issuer_url,omitempty"`
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
	if godotenv.Load(envPath) != nil {
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
	case "dokploy":
		configFile = "dokploy-config.json"
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

	// Override sensitive config from environment variables
	overrideFromEnv(&config.Keycloak.ServerURL, "KEYCLOAK_SERVER_URL")
	overrideFromEnv(&config.Keycloak.Realm, "KEYCLOAK_REALM")
	overrideFromEnv(&config.Keycloak.ClientID, "KEYCLOAK_CLIENT_ID")
	overrideFromEnv(&config.Keycloak.ClientSecret, "KEYCLOAK_CLIENT_SECRET")
	overrideFromEnv(&config.Keycloak.AdminUser, "KEYCLOAK_ADMIN")
	overrideFromEnv(&config.Keycloak.AdminPass, "KEYCLOAK_ADMIN_PASSWORD")
	overrideFromEnv(&config.Resend.APIKey, "RESEND_API_KEY")
	overrideFromEnv(&config.IDEncoder.Secret, "ID_ENCODER_SECRET")
	overrideFromEnv(&config.Database.Host, "DB_HOST")
	overrideFromEnv(&config.Database.Username, "DB_USER")
	overrideFromEnv(&config.Database.Password, "DB_PASSWORD")
	overrideFromEnv(&config.Keycloak.IssuerURL, "KEYCLOAK_ISSUER_URL")
	overrideFromEnv(&config.Cookie.Domain, "COOKIE_DOMAIN")
	config.Cookie.Secure = os.Getenv("COOKIE_SECURE") == "true"

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
	return c.Environment == "production" || c.Environment == "railway" || c.Environment == "dokploy"
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
	base := c.Keycloak.ServerURL
	if c.Keycloak.IssuerURL != "" {
		base = c.Keycloak.IssuerURL
	}
	return fmt.Sprintf("%s/realms/%s", base, c.Keycloak.Realm)
}

// overrideFromEnv sets the target value from the named environment variable if it is non-empty.
func overrideFromEnv(target *string, envKey string) {
	if v := os.Getenv(envKey); v != "" {
		*target = v
	}
}
