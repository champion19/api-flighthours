package handlers

import (
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGetBaseURL(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/test", func(c *gin.Context) {
		url := GetBaseURL(c)
		if url == "" {
			t.Error("expected non-empty base URL")
		}
		c.String(200, url)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
}

func TestSetLocationHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/test", func(c *gin.Context) {
		SetLocationHeader(c, "http://localhost:8080/api/v1", "resources", "123")
		c.Status(201)
	})

	req := httptest.NewRequest(http.MethodPost, "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	loc := w.Header().Get("Location")
	if loc == "" {
		t.Error("expected Location header to be set")
	}
}

func TestBuildResourceURL(t *testing.T) {
	url := BuildResourceURL("http://localhost:8080/api/v1", "airlines", "abc-123")
	if url == "" {
		t.Error("expected non-empty resource URL")
	}
	if !contains(url, "airlines") || !contains(url, "abc-123") {
		t.Errorf("expected URL to contain 'airlines' and 'abc-123', got '%s'", url)
	}
}

func TestBuildCollectionURL(t *testing.T) {
	url := BuildCollectionURL("http://localhost:8080/api/v1", "airlines")
	if url == "" {
		t.Error("expected non-empty collection URL")
	}
	if !contains(url, "airlines") {
		t.Errorf("expected URL to contain 'airlines', got '%s'", url)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsStr(s, substr))
}

func containsStr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestBuildResourceLinks(t *testing.T) {
	links := BuildResourceLinks("http://localhost:8080/api/v1", "airlines", "abc-123")
	if len(links) == 0 {
		t.Error("expected links to be generated")
	}
	selfFound := false
	for _, l := range links {
		if l.Rel == "self" {
			selfFound = true
		}
	}
	if !selfFound {
		t.Error("expected self link")
	}
}

func TestBuildAccountLinks(t *testing.T) {
	links := BuildAccountLinks("http://localhost:8080/api/v1", "acc-1")
	if len(links) == 0 {
		t.Error("expected account links")
	}
}

func TestBuildMessageLinks(t *testing.T) {
	links := BuildMessageLinks("http://localhost:8080/api/v1", "msg-1")
	if len(links) == 0 {
		t.Error("expected message links")
	}
}

func TestBuildMessageCreatedLinks(t *testing.T) {
	links := BuildMessageCreatedLinks("http://localhost:8080/api/v1", "msg-1")
	if len(links) == 0 {
		t.Error("expected message created links")
	}
}

func TestBuildMessageUpdatedLinks(t *testing.T) {
	links := BuildMessageUpdatedLinks("http://localhost:8080/api/v1", "msg-1")
	if len(links) == 0 {
		t.Error("expected message updated links")
	}
}

func TestBuildMessageListLinks(t *testing.T) {
	links := BuildMessageListLinks("http://localhost:8080/api/v1")
	if len(links) == 0 {
		t.Error("expected message list links")
	}
}

func TestBuildEmployeeLinks(t *testing.T) {
	links := BuildEmployeeLinks("http://localhost:8080/api/v1", "emp-1")
	if len(links) == 0 {
		t.Error("expected employee links")
	}
}

func TestBuildEmployeeMeLinks(t *testing.T) {
	links := BuildEmployeeMeLinks("http://localhost:8080/api/v1")
	if len(links) == 0 {
		t.Error("expected employee me links")
	}
}

func TestBuildAirlineLinks(t *testing.T) {
	links := BuildAirlineLinks("http://localhost:8080/api/v1", "airline-1")
	if len(links) == 0 {
		t.Error("expected airline links")
	}
}

func TestBuildAirlineListLinks(t *testing.T) {
	links := BuildAirlineListLinks("http://localhost:8080/api/v1")
	if len(links) == 0 {
		t.Error("expected airline list links")
	}
}

func TestBuildAirlineStatusLinks(t *testing.T) {
	links := BuildAirlineStatusLinks("http://localhost:8080/api/v1", "airline-1", true)
	if len(links) == 0 {
		t.Error("expected airline status links")
	}
	links2 := BuildAirlineStatusLinks("http://localhost:8080/api/v1", "airline-1", false)
	if len(links2) == 0 {
		t.Error("expected airline status links for inactive")
	}
}

func TestBuildAirportLinks(t *testing.T) {
	links := BuildAirportLinks("http://localhost:8080/api/v1", "airport-1")
	if len(links) == 0 {
		t.Error("expected airport links")
	}
}

func TestBuildAirportListLinks(t *testing.T) {
	links := BuildAirportListLinks("http://localhost:8080/api/v1")
	if len(links) == 0 {
		t.Error("expected airport list links")
	}
}

func TestBuildAirportStatusLinks(t *testing.T) {
	links := BuildAirportStatusLinks("http://localhost:8080/api/v1", "airport-1", true)
	if len(links) == 0 {
		t.Error("expected airport status links")
	}
	links2 := BuildAirportStatusLinks("http://localhost:8080/api/v1", "airport-1", false)
	if len(links2) == 0 {
		t.Error("expected airport status links for inactive")
	}
}

func TestBuildLicensePlateLinks(t *testing.T) {
	links := BuildLicensePlateLinks("http://localhost:8080/api/v1", "lp-1")
	if len(links) == 0 {
		t.Error("expected license plate links")
	}
}

func TestBuildLicensePlateListLinks(t *testing.T) {
	links := BuildLicensePlateListLinks("http://localhost:8080/api/v1")
	if len(links) == 0 {
		t.Error("expected license plate list links")
	}
}

func TestBuildLicensePlateCreatedLinks(t *testing.T) {
	links := BuildLicensePlateCreatedLinks("http://localhost:8080/api/v1", "lp-1")
	if len(links) == 0 {
		t.Error("expected license plate created links")
	}
}

func TestBuildAircraftModelLinks(t *testing.T) {
	links := BuildAircraftModelLinks("http://localhost:8080/api/v1", "model-1")
	if len(links) == 0 {
		t.Error("expected aircraft model links")
	}
}

func TestBuildAircraftModelListLinks(t *testing.T) {
	links := BuildAircraftModelListLinks("http://localhost:8080/api/v1")
	if len(links) == 0 {
		t.Error("expected aircraft model list links")
	}
}

func TestBuildAircraftFamilyListLinks(t *testing.T) {
	links := BuildAircraftFamilyListLinks("http://localhost:8080/api/v1")
	if len(links) == 0 {
		t.Error("expected aircraft family list links")
	}
}

func TestBuildAircraftModelStatusLinks(t *testing.T) {
	links := BuildAircraftModelStatusLinks("http://localhost:8080/api/v1", "model-1", true)
	if len(links) == 0 {
		t.Error("expected aircraft model status links")
	}
	links2 := BuildAircraftModelStatusLinks("http://localhost:8080/api/v1", "model-1", false)
	if len(links2) == 0 {
		t.Error("expected aircraft model status links for inactive")
	}
}

func TestBuildRouteLinks(t *testing.T) {
	links := BuildRouteLinks("http://localhost:8080/api/v1", "route-1")
	if len(links) == 0 {
		t.Error("expected route links")
	}
}

func TestBuildRouteListLinks(t *testing.T) {
	links := BuildRouteListLinks("http://localhost:8080/api/v1")
	if len(links) == 0 {
		t.Error("expected route list links")
	}
}

func TestBuildAirlineRouteLinks(t *testing.T) {
	links := BuildAirlineRouteLinks("http://localhost:8080/api/v1", "ar-1", true)
	if len(links) == 0 {
		t.Error("expected airline route links")
	}
	links2 := BuildAirlineRouteLinks("http://localhost:8080/api/v1", "ar-1", false)
	if len(links2) == 0 {
		t.Error("expected airline route links for inactive")
	}
}

func TestBuildAirlineRouteListLinks(t *testing.T) {
	links := BuildAirlineRouteListLinks("http://localhost:8080/api/v1")
	if len(links) == 0 {
		t.Error("expected airline route list links")
	}
}

func TestBuildAirlineRouteStatusLinks(t *testing.T) {
	links := BuildAirlineRouteStatusLinks("http://localhost:8080/api/v1", "ar-1", true)
	if len(links) == 0 {
		t.Error("expected airline route status links")
	}
	links2 := BuildAirlineRouteStatusLinks("http://localhost:8080/api/v1", "ar-1", false)
	if len(links2) == 0 {
		t.Error("expected airline route status links for inactive")
	}
}

func TestBuildDailyLogbookDetailLinks_Gin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/test", func(c *gin.Context) {
		links := BuildDailyLogbookDetailLinks(c, "detail-1")
		if len(links) == 0 {
			t.Error("expected daily logbook detail links")
		}
		c.Status(200)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
}

func TestBuildEngineLinks(t *testing.T) {
	links := BuildEngineLinks("http://localhost:8080/api/v1", "engine-1")
	if len(links) == 0 {
		t.Error("expected engine links")
	}
}

func TestBuildEngineListLinks(t *testing.T) {
	links := BuildEngineListLinks("http://localhost:8080/api/v1")
	if len(links) == 0 {
		t.Error("expected engine list links")
	}
}

func TestBuildManufacturerLinks(t *testing.T) {
	links := BuildManufacturerLinks("http://localhost:8080/api/v1", "mfr-1")
	if len(links) == 0 {
		t.Error("expected manufacturer links")
	}
}

func TestBuildManufacturerListLinks(t *testing.T) {
	links := BuildManufacturerListLinks("http://localhost:8080/api/v1")
	if len(links) == 0 {
		t.Error("expected manufacturer list links")
	}
}

func TestBuildAirlineEmployeeLinks(t *testing.T) {
	links := BuildAirlineEmployeeLinks("http://localhost:8080/api/v1", "ae-1")
	if len(links) == 0 {
		t.Error("expected airline employee links")
	}
}

func TestBuildAirlineEmployeeStatusLinks(t *testing.T) {
	links := BuildAirlineEmployeeStatusLinks("http://localhost:8080/api/v1", "ae-1", true)
	if len(links) == 0 {
		t.Error("expected airline employee status links")
	}
	// Also test deactivate variant
	links2 := BuildAirlineEmployeeStatusLinks("http://localhost:8080/api/v1", "ae-1", false)
	if len(links2) == 0 {
		t.Error("expected airline employee status links for inactive")
	}
}

func TestBuildAirlineEmployeeCreatedLinks(t *testing.T) {
	links := BuildAirlineEmployeeCreatedLinks("http://localhost:8080/api/v1", "ae-1")
	if len(links) == 0 {
		t.Error("expected airline employee created links")
	}
}

func TestBuildDailyLogbookDeletedLinks(t *testing.T) {
	links := BuildDailyLogbookDeletedLinks("http://localhost:8080/api/v1")
	if len(links) != 2 {
		t.Errorf("expected 2 links, got %d", len(links))
	}
	if links[0].Rel != "list" {
		t.Errorf("expected first link rel to be 'list', got %q", links[0].Rel)
	}
	if links[1].Rel != "create" {
		t.Errorf("expected second link rel to be 'create', got %q", links[1].Rel)
	}
}

func TestGetBaseURL_TLS(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "https://example.com/test", nil)
	c.Request.TLS = &tls.ConnectionState{}

	baseURL := GetBaseURL(c)
	expected := "https://example.com"
	if baseURL != expected {
		t.Errorf("expected %s, got %s", expected, baseURL)
	}
}
