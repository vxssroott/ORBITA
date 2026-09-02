package api_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/vxssroott/ORBITA/services/api"
)

func TestHealthEndpoint(t *testing.T) {
	server := api.NewServer(0)

	request := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	response := httptest.NewRecorder()

	server.Handler().ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected HTTP 200, got %d", response.Code)
	}

	if !strings.Contains(response.Body.String(), "ok") {
		t.Fatalf("expected health response to contain ok, got %q", response.Body.String())
	}
}

func TestReadyEndpoint(t *testing.T) {
	server := api.NewServer(0)

	request := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	response := httptest.NewRecorder()

	server.Handler().ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected HTTP 200, got %d", response.Code)
	}

	if !strings.Contains(response.Body.String(), "ready") {
		t.Fatalf("expected readiness response to contain ready, got %q", response.Body.String())
	}
}

func TestStatusEndpoint(t *testing.T) {
	server := api.NewServer(0)

	request := httptest.NewRequest(http.MethodGet, "/v1/status", nil)
	response := httptest.NewRecorder()

	server.Handler().ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected HTTP 200, got %d", response.Code)
	}

	if !strings.Contains(response.Header().Get("Content-Type"), "application/json") {
		t.Fatalf("expected JSON content type, got %q", response.Header().Get("Content-Type"))
	}

	if !strings.Contains(response.Body.String(), "orbita-api") {
		t.Fatalf("expected service identifier in response, got %q", response.Body.String())
	}
}

func TestUnknownEndpointReturnsNotFound(t *testing.T) {
	server := api.NewServer(0)

	request := httptest.NewRequest(http.MethodGet, "/v1/does-not-exist", nil)
	response := httptest.NewRecorder()

	server.Handler().ServeHTTP(response, request)

	if response.Code != http.StatusNotFound {
		t.Fatalf("expected HTTP 404, got %d", response.Code)
	}
}

func TestSecurityHeadersArePresent(t *testing.T) {
	server := api.NewServer(0)

	request := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	response := httptest.NewRecorder()

	server.Handler().ServeHTTP(response, request)

	if response.Header().Get("X-Content-Type-Options") == "" {
		t.Fatal("expected X-Content-Type-Options header")
	}

	if response.Header().Get("X-Frame-Options") == "" {
		t.Fatal("expected X-Frame-Options header")
	}

	if response.Header().Get("X-Request-ID") == "" {
		t.Fatal("expected X-Request-ID header")
	}
}
