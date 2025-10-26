package mail

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// setupRouter creates a test router with mail protocol routes
func setupRouter() *gin.Engine {
	router := gin.New()
	group := router.Group("/protocols")
	InitializeMailProtocolRoutes(group)
	return router
}

func TestGetLogsHandler_InvalidTailNonInteger(t *testing.T) {
	router := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/protocols/mail/test-id/logs?tail=invalid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if resp["error"] != "invalid tail parameter" {
		t.Fatalf("expected 'invalid tail parameter' error, got %v", resp["error"])
	}
}

func TestGetLogsHandler_InvalidTailOutOfRange(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		tail string
		want int
	}{
		{"tail too low", "0", http.StatusBadRequest},
		{"tail too high", "5001", http.StatusBadRequest},
		{"tail negative", "-1", http.StatusBadRequest},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			router := setupRouter()

			req := httptest.NewRequest(http.MethodGet, "/protocols/mail/test-id/logs?tail="+tc.tail, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tc.want {
				t.Fatalf("expected status %d, got %d", tc.want, w.Code)
			}

			var resp map[string]interface{}
			if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
				t.Fatalf("failed to parse response: %v", err)
			}

			if resp["error"] != "tail must be between 1 and 5000" {
				t.Fatalf("expected 'tail must be between 1 and 5000' error, got %v", resp["error"])
			}
		})
	}
}

func TestGetLogsHandler_ContainerNotFound(t *testing.T) {
	router := setupRouter()

	// Use a container ID that definitely doesn't exist
	req := httptest.NewRequest(http.MethodGet, "/protocols/mail/nonexistent-container-id-12345/logs", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Expecting 404 since container won't exist in database
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", w.Code)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if resp["error"] != "container not found" {
		t.Fatalf("expected 'container not found' error, got %v", resp["error"])
	}
}

func TestGetLogsHandler_ValidTailDefault(t *testing.T) {
	router := setupRouter()

	// No tail parameter should use default of 500
	req := httptest.NewRequest(http.MethodGet, "/protocols/mail/test-id/logs", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Will return 404 for non-existent container, but proves tail parsing succeeded
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status 404 for non-existent container, got %d", w.Code)
	}
}

func TestGetLogsHandler_ValidTailInRange(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		tail string
	}{
		{"tail minimum", "1"},
		{"tail medium", "500"},
		{"tail maximum", "5000"},
		{"tail typical", "100"},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			router := setupRouter()

			req := httptest.NewRequest(http.MethodGet, "/protocols/mail/test-id/logs?tail="+tc.tail, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			// Should pass validation (tail parsing) and fail at container lookup
			// Status should NOT be 400 (bad request for tail parameter)
			if w.Code == http.StatusBadRequest {
				var resp map[string]interface{}
				_ = json.Unmarshal(w.Body.Bytes(), &resp)
				// If we got 400, it shouldn't be because of tail parameter
				if err, ok := resp["error"].(string); ok {
					if err == "invalid tail parameter" || err == "tail must be between 1 and 5000" {
						t.Fatalf("valid tail %s was rejected: %s", tc.tail, err)
					}
				}
			}
		})
	}
}
