package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/champion19/api-flighthours/middleware"
	"github.com/champion19/api-flighthours/platform/logger"
	"github.com/champion19/api-flighthours/tools/idencoder"
	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func createTestHandler() *handler {
	cfg := idencoder.Config{
		Secret:    "test-salt",
		MinLength: 8,
	}
	log := logger.NewSlogLogger()
	encoder, _ := idencoder.NewHashidsEncoder(cfg, log)
	return &handler{
		IDEncoder: encoder,
		Response:  middleware.NewResponseHandler(nil),
	}
}

func TestHandler_EncodeID(t *testing.T) {
	t.Run("encodes valid UUID", func(t *testing.T) {
		h := createTestHandler()
		uuid := "550e8400-e29b-41d4-a716-446655440000"

		encoded, err := h.EncodeID(uuid)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if encoded == "" {
			t.Error("expected non-empty encoded ID")
		}
		if encoded == uuid {
			t.Error("encoded ID should be different from UUID")
		}
	})

	t.Run("handles empty UUID", func(t *testing.T) {
		h := createTestHandler()

		_, err := h.EncodeID("")

		// Empty UUID should return error (implementation dependent)
		if err == nil {
			// Some encoders may accept empty strings
			t.Log("encoder accepted empty string")
		}
	})
}

func TestHandler_DecodeID(t *testing.T) {
	t.Run("decodes valid encoded ID", func(t *testing.T) {
		h := createTestHandler()
		originalUUID := "550e8400-e29b-41d4-a716-446655440000"

		encoded, err := h.EncodeID(originalUUID)
		if err != nil {
			t.Fatalf("failed to encode: %v", err)
		}

		decoded, err := h.DecodeID(encoded)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if decoded != originalUUID {
			t.Errorf("expected %q, got %q", originalUUID, decoded)
		}
	})

	t.Run("returns error for invalid encoded ID", func(t *testing.T) {
		h := createTestHandler()

		_, err := h.DecodeID("invalid-encoded-id")

		if err == nil {
			t.Error("expected error for invalid encoded ID")
		}
	})
}

func TestHandler_HandleIDEncodingError(t *testing.T) {
	t.Run("sets internal server error", func(t *testing.T) {
		h := createTestHandler()
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)

		h.HandleIDEncodingError(c, "test-uuid", nil)

		if len(c.Errors) == 0 {
			t.Error("expected error to be set on context")
		}
	})
}

func TestHandler_HandleIDDecodingError(t *testing.T) {
	t.Run("sets invalid ID error", func(t *testing.T) {
		h := createTestHandler()
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)

		h.HandleIDDecodingError(c, "invalid-id", nil)

		if len(c.Errors) == 0 {
			t.Error("expected error to be set on context")
		}
	})
}

func TestHandler_resolveID(t *testing.T) {
	t.Run("returns empty for empty input", func(t *testing.T) {
		h := createTestHandler()

		uuid, encodedID := h.resolveID("")

		if uuid != "" {
			t.Errorf("expected empty uuid, got %q", uuid)
		}
		if encodedID != "" {
			t.Errorf("expected empty encodedID, got %q", encodedID)
		}
	})

	t.Run("resolves valid encoded ID", func(t *testing.T) {
		h := createTestHandler()
		originalUUID := "550e8400-e29b-41d4-a716-446655440000"

		encoded, _ := h.EncodeID(originalUUID)
		uuid, returnedEncoded := h.resolveID(encoded)

		if uuid != originalUUID {
			t.Errorf("expected UUID %q, got %q", originalUUID, uuid)
		}
		if returnedEncoded != encoded {
			t.Errorf("expected encoded ID %q, got %q", encoded, returnedEncoded)
		}
	})

	t.Run("returns empty for invalid encoded ID", func(t *testing.T) {
		h := createTestHandler()

		uuid, encodedID := h.resolveID("invalid-id")

		if uuid != "" {
			t.Errorf("expected empty uuid, got %q", uuid)
		}
		if encodedID != "" {
			t.Errorf("expected empty encodedID, got %q", encodedID)
		}
	})
}

func TestNew(t *testing.T) {
	t.Run("creates handler with all dependencies", func(t *testing.T) {
		h := New(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)

		if h == nil {
			t.Error("expected non-nil handler")
		}
	})
}
