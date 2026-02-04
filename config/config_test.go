package config

import (
	"testing"
)

func TestConfig_GetMySQLDSN(t *testing.T) {
	t.Run("returns correct DSN format", func(t *testing.T) {
		cfg := &Config{
			Database: Database{
				Username: "testuser",
				Password: "testpass",
				Host:     "localhost",
				Port:     "3306",
				Name:     "testdb",
			},
		}

		dsn := cfg.GetMySQLDSN()
		expected := "testuser:testpass@tcp(localhost:3306)/testdb?parseTime=true&loc=Local"
		if dsn != expected {
			t.Errorf("expected %q, got %q", expected, dsn)
		}
	})

	t.Run("includes SSL when configured", func(t *testing.T) {
		cfg := &Config{
			Database: Database{
				Username: "user",
				Password: "pass",
				Host:     "db.example.com",
				Port:     "3306",
				Name:     "mydb",
				SSL:      "true",
			},
		}

		dsn := cfg.GetMySQLDSN()
		if dsn == "" {
			t.Error("expected non-empty DSN")
		}
		// Should contain tls parameter
		if !contains(dsn, "&tls=true") {
			t.Errorf("expected DSN to contain '&tls=true', got %q", dsn)
		}
	})
}

func TestConfig_GetServerAddress(t *testing.T) {
	t.Run("returns correct address format", func(t *testing.T) {
		cfg := &Config{
			Server: Server{
				Host: "0.0.0.0",
				Port: "8080",
			},
		}

		addr := cfg.GetServerAddress()
		expected := "0.0.0.0:8080"
		if addr != expected {
			t.Errorf("expected %q, got %q", expected, addr)
		}
	})
}

func TestConfig_IsProduction(t *testing.T) {
	t.Run("returns true for production environment", func(t *testing.T) {
		cfg := &Config{Environment: "production"}
		if !cfg.IsProduction() {
			t.Error("expected true for production")
		}
	})

	t.Run("returns true for railway environment", func(t *testing.T) {
		cfg := &Config{Environment: "railway"}
		if !cfg.IsProduction() {
			t.Error("expected true for railway")
		}
	})

	t.Run("returns false for local environment", func(t *testing.T) {
		cfg := &Config{Environment: "local"}
		if cfg.IsProduction() {
			t.Error("expected false for local")
		}
	})

	t.Run("returns false for development environment", func(t *testing.T) {
		cfg := &Config{Environment: "development"}
		if cfg.IsProduction() {
			t.Error("expected false for development")
		}
	})
}

func TestConfig_GetKeycloakAuthURL(t *testing.T) {
	t.Run("returns correct auth URL format", func(t *testing.T) {
		cfg := &Config{
			Keycloak: KeycloakConfig{
				ServerURL: "https://keycloak.example.com",
				Realm:     "myrealm",
			},
		}

		url := cfg.GetKeycloakAuthURL()
		expected := "https://keycloak.example.com/realms/myrealm/protocol/openid-connect/token"
		if url != expected {
			t.Errorf("expected %q, got %q", expected, url)
		}
	})
}

func TestConfig_GetKeycloakAdminURL(t *testing.T) {
	t.Run("returns correct admin URL format", func(t *testing.T) {
		cfg := &Config{
			Keycloak: KeycloakConfig{
				ServerURL: "https://keycloak.example.com",
				Realm:     "myrealm",
			},
		}

		url := cfg.GetKeycloakAdminURL()
		expected := "https://keycloak.example.com/admin/realms/myrealm"
		if url != expected {
			t.Errorf("expected %q, got %q", expected, url)
		}
	})
}

func TestConfig_GetKeycloakJWKSURL(t *testing.T) {
	t.Run("returns correct JWKS URL format", func(t *testing.T) {
		cfg := &Config{
			Keycloak: KeycloakConfig{
				ServerURL: "https://keycloak.example.com",
				Realm:     "myrealm",
			},
		}

		url := cfg.GetKeycloakJWKSURL()
		expected := "https://keycloak.example.com/realms/myrealm/protocol/openid-connect/certs"
		if url != expected {
			t.Errorf("expected %q, got %q", expected, url)
		}
	})
}

func TestConfig_GetKeycloakIssuerURL(t *testing.T) {
	t.Run("returns correct issuer URL format", func(t *testing.T) {
		cfg := &Config{
			Keycloak: KeycloakConfig{
				ServerURL: "https://keycloak.example.com",
				Realm:     "myrealm",
			},
		}

		url := cfg.GetKeycloakIssuerURL()
		expected := "https://keycloak.example.com/realms/myrealm"
		if url != expected {
			t.Errorf("expected %q, got %q", expected, url)
		}
	})
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsSubstring(s, substr))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
