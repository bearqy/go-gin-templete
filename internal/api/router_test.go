package api

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-gin-templete/internal/webserver"
)

func TestHealthRoute(t *testing.T) {
	engine := webserver.NewEngine(Router, io.Discard)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/health", nil)

	engine.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Code)
	}
	if recorder.Body.String() != "OK" {
		t.Fatalf("expected OK body, got %q", recorder.Body.String())
	}
}

func TestAPIV1HomeRoute(t *testing.T) {
	engine := webserver.NewEngine(Router, io.Discard)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/v1/home", nil)

	engine.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Code)
	}
	var body struct {
		Data string `json:"data"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body.Data != "Hello, World!" {
		t.Fatalf("unexpected body: %+v", body)
	}
}
