package handlers

import (
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestGetBaseURL(t *testing.T) {
	t.Run("returns http scheme without TLS", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)
		c.Request.Host = "localhost:8080"

		result := GetBaseURL(c)
		expected := "http://localhost:8080"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("returns https scheme with TLS", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)
		c.Request.Host = "api.example.com"
		c.Request.TLS = &tls.ConnectionState{} // Simulate HTTPS

		result := GetBaseURL(c)
		expected := "https://api.example.com"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})
}

func TestSetLocationHeader(t *testing.T) {
	t.Run("sets Location header correctly", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		SetLocationHeader(c, "http://localhost:8080", "employees", "emp-123")

		location := w.Header().Get("Location")
		expected := "http://localhost:8080/flighthours/api/v1/employees/emp-123"
		if location != expected {
			t.Errorf("expected %q, got %q", expected, location)
		}
	})
}

func TestBuildResourceURL(t *testing.T) {
	tests := []struct {
		name       string
		baseURL    string
		resource   string
		resourceID string
		expected   string
	}{
		{
			name:       "builds employee URL",
			baseURL:    "http://localhost:8080",
			resource:   "employees",
			resourceID: "emp-123",
			expected:   "http://localhost:8080/flighthours/api/v1/employees/emp-123",
		},
		{
			name:       "builds airline URL",
			baseURL:    "https://api.example.com",
			resource:   "airlines",
			resourceID: "airline-456",
			expected:   "https://api.example.com/flighthours/api/v1/airlines/airline-456",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BuildResourceURL(tt.baseURL, tt.resource, tt.resourceID)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestBuildCollectionURL(t *testing.T) {
	result := BuildCollectionURL("http://localhost:8080", "messages")
	expected := "http://localhost:8080/flighthours/api/v1/messages"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestBuildResourceLinks(t *testing.T) {
	links := BuildResourceLinks("http://localhost:8080", "employees", "emp-123")

	if len(links) != 4 {
		t.Fatalf("expected 4 links, got %d", len(links))
	}

	// Check self link
	if links[0].Rel != "self" || links[0].Method != "GET" {
		t.Errorf("expected self GET link, got %s %s", links[0].Rel, links[0].Method)
	}

	// Check update link
	if links[1].Rel != "update" || links[1].Method != "PUT" {
		t.Errorf("expected update PUT link, got %s %s", links[1].Rel, links[1].Method)
	}

	// Check delete link
	if links[2].Rel != "delete" || links[2].Method != "DELETE" {
		t.Errorf("expected delete DELETE link, got %s %s", links[2].Rel, links[2].Method)
	}

	// Check collection link
	if links[3].Rel != "collection" || links[3].Method != "GET" {
		t.Errorf("expected collection GET link, got %s %s", links[3].Rel, links[3].Method)
	}
}

func TestBuildMessageLinks(t *testing.T) {
	links := BuildMessageLinks("http://localhost:8080", "msg-123")
	if len(links) != 4 {
		t.Errorf("expected 4 links, got %d", len(links))
	}
}

func TestBuildMessageCreatedLinks(t *testing.T) {
	links := BuildMessageCreatedLinks("http://localhost:8080", "msg-123")
	if len(links) != 4 {
		t.Errorf("expected 4 links, got %d", len(links))
	}
	// Check update link exists
	hasUpdate := false
	for _, l := range links {
		if l.Rel == "update" && l.Method == "PUT" {
			hasUpdate = true
		}
	}
	if !hasUpdate {
		t.Error("expected update link in created links")
	}
}

func TestBuildMessageUpdatedLinks(t *testing.T) {
	links := BuildMessageUpdatedLinks("http://localhost:8080", "msg-123")
	if len(links) != 3 {
		t.Errorf("expected 3 links, got %d", len(links))
	}
}

func TestBuildMessageListLinks(t *testing.T) {
	links := BuildMessageListLinks("http://localhost:8080")
	if len(links) != 2 {
		t.Errorf("expected 2 links, got %d", len(links))
	}
}

func TestBuildEmployeeLinks(t *testing.T) {
	links := BuildEmployeeLinks("http://localhost:8080", "emp-123")
	if len(links) != 4 {
		t.Errorf("expected 4 links, got %d", len(links))
	}
}

func TestBuildEmployeeMeLinks(t *testing.T) {
	links := BuildEmployeeMeLinks("http://localhost:8080")
	if len(links) != 3 {
		t.Errorf("expected 3 links, got %d", len(links))
	}
	if links[0].Href != "http://localhost:8080/flighthours/api/v1/employees" {
		t.Errorf("unexpected me URL: %s", links[0].Href)
	}
}

func TestBuildAirlineLinks(t *testing.T) {
	links := BuildAirlineLinks("http://localhost:8080", "airline-123")
	if len(links) != 4 {
		t.Errorf("expected 4 links, got %d", len(links))
	}
	// Check for activate link
	hasActivate := false
	for _, l := range links {
		if l.Rel == "activate" && l.Method == "PATCH" {
			hasActivate = true
		}
	}
	if !hasActivate {
		t.Error("expected activate link in airline links")
	}
}

func TestBuildAirlineListLinks(t *testing.T) {
	links := BuildAirlineListLinks("http://localhost:8080")
	if len(links) != 1 {
		t.Errorf("expected 1 link, got %d", len(links))
	}
}

func TestBuildAirlineStatusLinks(t *testing.T) {
	t.Run("when active shows deactivate link", func(t *testing.T) {
		links := BuildAirlineStatusLinks("http://localhost:8080", "airline-123", true)
		hasDeactivate := false
		for _, l := range links {
			if l.Rel == "deactivate" {
				hasDeactivate = true
			}
		}
		if !hasDeactivate {
			t.Error("expected deactivate link when active")
		}
	})

	t.Run("when inactive shows activate link", func(t *testing.T) {
		links := BuildAirlineStatusLinks("http://localhost:8080", "airline-123", false)
		hasActivate := false
		for _, l := range links {
			if l.Rel == "activate" {
				hasActivate = true
			}
		}
		if !hasActivate {
			t.Error("expected activate link when inactive")
		}
	})
}

func TestBuildRouteLinks(t *testing.T) {
	links := BuildRouteLinks("http://localhost:8080", "route-123")
	if len(links) != 2 {
		t.Errorf("expected 2 links, got %d", len(links))
	}
}

func TestBuildRouteListLinks(t *testing.T) {
	links := BuildRouteListLinks("http://localhost:8080")
	if len(links) != 1 {
		t.Errorf("expected 1 link, got %d", len(links))
	}
}

func TestBuildAirlineRouteLinks(t *testing.T) {
	t.Run("when active shows deactivate link", func(t *testing.T) {
		links := BuildAirlineRouteLinks("http://localhost:8080", "ar-123", true)
		hasDeactivate := false
		for _, l := range links {
			if l.Rel == "deactivate" {
				hasDeactivate = true
			}
		}
		if !hasDeactivate {
			t.Error("expected deactivate link when active")
		}
	})

	t.Run("when inactive shows activate link", func(t *testing.T) {
		links := BuildAirlineRouteLinks("http://localhost:8080", "ar-123", false)
		hasActivate := false
		for _, l := range links {
			if l.Rel == "activate" {
				hasActivate = true
			}
		}
		if !hasActivate {
			t.Error("expected activate link when inactive")
		}
	})
}

func TestBuildAirlineRouteListLinks(t *testing.T) {
	links := BuildAirlineRouteListLinks("http://localhost:8080")
	if len(links) != 1 {
		t.Errorf("expected 1 link, got %d", len(links))
	}
}

func TestBuildAirlineRouteStatusLinks(t *testing.T) {
	t.Run("when active", func(t *testing.T) {
		links := BuildAirlineRouteStatusLinks("http://localhost:8080", "ar-123", true)
		if len(links) != 3 {
			t.Errorf("expected 3 links, got %d", len(links))
		}
	})

	t.Run("when inactive", func(t *testing.T) {
		links := BuildAirlineRouteStatusLinks("http://localhost:8080", "ar-123", false)
		if len(links) != 3 {
			t.Errorf("expected 3 links, got %d", len(links))
		}
	})
}

func TestBuildEngineLinks(t *testing.T) {
	links := BuildEngineLinks("http://localhost:8080", "engine-123")
	if len(links) != 2 {
		t.Errorf("expected 2 links, got %d", len(links))
	}
}

func TestBuildEngineListLinks(t *testing.T) {
	links := BuildEngineListLinks("http://localhost:8080")
	if len(links) != 1 {
		t.Errorf("expected 1 link, got %d", len(links))
	}
}

func TestBuildAirlineEmployeeLinks(t *testing.T) {
	links := BuildAirlineEmployeeLinks("http://localhost:8080", "emp-123")
	if len(links) != 5 {
		t.Errorf("expected 5 links, got %d", len(links))
	}
}

func TestBuildAirlineEmployeeStatusLinks(t *testing.T) {
	t.Run("when active shows deactivate", func(t *testing.T) {
		links := BuildAirlineEmployeeStatusLinks("http://localhost:8080", "emp-123", true)
		hasDeactivate := false
		for _, l := range links {
			if l.Rel == "deactivate" {
				hasDeactivate = true
			}
		}
		if !hasDeactivate {
			t.Error("expected deactivate link when active")
		}
	})

	t.Run("when inactive shows activate", func(t *testing.T) {
		links := BuildAirlineEmployeeStatusLinks("http://localhost:8080", "emp-123", false)
		hasActivate := false
		for _, l := range links {
			if l.Rel == "activate" {
				hasActivate = true
			}
		}
		if !hasActivate {
			t.Error("expected activate link when inactive")
		}
	})
}

func TestBuildAirlineEmployeeCreatedLinks(t *testing.T) {
	links := BuildAirlineEmployeeCreatedLinks("http://localhost:8080", "emp-123")
	if len(links) != 3 {
		t.Errorf("expected 3 links, got %d", len(links))
	}
}

func TestBuildAccountLinks(t *testing.T) {
	links := BuildAccountLinks("http://localhost:8080", "acc-123")
	if len(links) != 4 {
		t.Errorf("expected 4 links, got %d", len(links))
	}
}
