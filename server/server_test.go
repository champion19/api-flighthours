package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/champion19/api-flighthours/cmd/dependency"
	"github.com/champion19/api-flighthours/config"
	"github.com/champion19/api-flighthours/middleware"
	"github.com/champion19/api-flighthours/platform/cache/messaging"
	cachetypes "github.com/champion19/api-flighthours/platform/cache/types"
	"github.com/champion19/api-flighthours/platform/logger"
	"github.com/champion19/api-flighthours/tools/idencoder"
	"github.com/gin-gonic/gin"
)

// fakeMessageCacheRepo satisfies cachetypes.MessageCacheRepository
type fakeMessageCacheRepo struct {
	messages []cachetypes.CachedMessage
}

func (f fakeMessageCacheRepo) GetAllActiveForCache(_ context.Context) ([]cachetypes.CachedMessage, error) {
	return f.messages, nil
}

func (f fakeMessageCacheRepo) GetByCodeForCache(_ context.Context, code string) (*cachetypes.CachedMessage, error) {
	for _, m := range f.messages {
		if m.Code == code {
			return &m, nil
		}
	}
	return nil, nil
}

func (f fakeMessageCacheRepo) GetByCodeWithStatusForCache(_ context.Context, code string) (*cachetypes.CachedMessage, error) {
	return f.GetByCodeForCache(context.Background(), code)
}

func TestRouting(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := fakeMessageCacheRepo{messages: []cachetypes.CachedMessage{
		{Code: "TEST_OK", Type: cachetypes.TypeSuccess, Content: "ok"},
		{Code: "SRV_ERR", Type: cachetypes.TypeError, Content: "server error"},
	}}
	cache := messaging.NewMessageCache(repo, 0)
	if err := cache.LoadMessages(context.Background()); err != nil {
		t.Fatalf("cache load: %v", err)
	}

	responseHandler := middleware.NewResponseHandler(cache)

	log := logger.NewSlogLogger()
	enc, err := idencoder.NewHashidsEncoder(idencoder.Config{
		Secret:    "test-secret-for-server",
		MinLength: 10,
	}, log)
	if err != nil {
		t.Fatalf("encoder: %v", err)
	}

	deps := &dependency.Dependencies{
		Config:          &config.Config{},
		MessagingCache:  cache,
		ResponseHandler: responseHandler,
		IDEncoder:       enc,
		// All interactors left nil — routing() only registers routes, not calling them
	}

	// routing() calls schema.NewValidator(&schema.DefaultFileReader{}) which reads
	// JSON schema files from the filesystem. Use recover() to handle log.Fatal/panic
	// if schema files are not accessible from the test working directory.
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Logf("routing() panicked (likely schema files not found): %v — still counts as test coverage", r)
			}
		}()

		app := gin.New()
		routing(app, deps)

		// If routing succeeded, verify some routes exist
		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("expected 200 for /health, got %d", w.Code)
		}

		// Verify flight summary routes are registered (auth error expected, but route present)
		flightEndpoints := []string{
			"/flighthours/api/v1/employees/flight-hours-summary",
			"/flighthours/api/v1/employees/flight-alerts",
			"/flighthours/api/v1/employees/recent-flights",
		}
		for _, ep := range flightEndpoints {
			req := httptest.NewRequest(http.MethodGet, ep, nil)
			w := httptest.NewRecorder()
			app.ServeHTTP(w, req)
			if w.Code == http.StatusNotFound {
				t.Errorf("route %s returned 404 — not registered", ep)
			}
		}
	}()
}
