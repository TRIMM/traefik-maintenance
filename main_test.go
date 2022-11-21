package traefik_maintenance

import (
	"context"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
)

func TestMaintenancePage(t *testing.T) {
	maintenancePage, err := filepath.Abs("./maintenance_test.html")
	if err != nil {
		t.Fatal(err)
	}

	maintenanceTrigger, err := filepath.Abs("./maintenance_test.trigger")
	if err != nil {
		t.Fatal(err)
	}

	cfg := CreateConfig()
	cfg.Enabled = true
	cfg.Filename = maintenancePage
	cfg.TriggerFilename = maintenanceTrigger
	cfg.HttpResponseCode = http.StatusServiceUnavailable
	cfg.HttpContentType = "text/html; charset=utf-8"

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, err := New(ctx, next, cfg, "traefik-maintenance")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler.ServeHTTP(recorder, req)

	assertResponseStatus(t, recorder, http.StatusServiceUnavailable)
	assertResponseHeader(t, recorder, "Content-Type", "text/html; charset=utf-8")
	assertResponseBody(t, recorder, "<html><head></head><body>Maintenance</body></html>")
}

func TestMaintenancePageWithOtherStatusCodeAndContentType(t *testing.T) {
	maintenancePage, err := filepath.Abs("./maintenance_test.json")
	if err != nil {
		t.Fatal(err)
	}

	maintenanceTrigger, err := filepath.Abs("./maintenance_test.trigger")
	if err != nil {
		t.Fatal(err)
	}

	cfg := CreateConfig()
	cfg.Enabled = true
	cfg.Filename = maintenancePage
	cfg.TriggerFilename = maintenanceTrigger
	cfg.HttpResponseCode = http.StatusTeapot
	cfg.HttpContentType = "application/json; charset=utf-8"

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, err := New(ctx, next, cfg, "traefik-maintenance")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler.ServeHTTP(recorder, req)

	assertResponseStatus(t, recorder, http.StatusTeapot)
	assertResponseHeader(t, recorder, "Content-Type", "application/json; charset=utf-8")
	assertResponseBody(t, recorder, "{\"detail\": \"This endpoint is currently in maintenance mode\"}")
}

func TestMaintenancePageWithoutTrigger(t *testing.T) {
	maintenancePage, err := filepath.Abs("./maintenance_test.html")
	if err != nil {
		t.Fatal(err)
	}

	cfg := CreateConfig()
	cfg.Enabled = true
	cfg.Filename = maintenancePage
	cfg.TriggerFilename = "./missing.trigger"
	cfg.HttpResponseCode = http.StatusServiceUnavailable
	cfg.HttpContentType = "text/html; charset=utf-8"

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, err := New(ctx, next, cfg, "traefik-maintenance")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler.ServeHTTP(recorder, req)

	assertEmptyContentTypeHeader(t, recorder)
	assertEmptyResponseBody(t, recorder)
}

func TestMaintenancePageWithMissingTrigger(t *testing.T) {
	maintenancePage, err := filepath.Abs("./maintenance_test.html")
	if err != nil {
		t.Fatal(err)
	}

	cfg := CreateConfig()
	cfg.Enabled = true
	cfg.Filename = maintenancePage
	cfg.HttpResponseCode = http.StatusServiceUnavailable
	cfg.HttpContentType = "text/html; charset=utf-8"

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, err := New(ctx, next, cfg, "traefik-maintenance")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler.ServeHTTP(recorder, req)

	assertResponseHeader(t, recorder, "Content-Type", "text/html; charset=utf-8")
	assertResponseBody(t, recorder, "<html><head></head><body>Maintenance</body></html>")
}

func TestDisabledMaintenancePage(t *testing.T) {
	maintenancePage, err := filepath.Abs("./maintenance_test.html")
	if err != nil {
		t.Fatal(err)
	}

	maintenanceTrigger, err := filepath.Abs("./maintenance_test.trigger")
	if err != nil {
		t.Fatal(err)
	}

	cfg := CreateConfig()
	cfg.Enabled = false
	cfg.Filename = maintenancePage
	cfg.TriggerFilename = maintenanceTrigger
	cfg.HttpResponseCode = http.StatusServiceUnavailable
	cfg.HttpContentType = "text/html; charset=utf-8"

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, err := New(ctx, next, cfg, "traefik-maintenance")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler.ServeHTTP(recorder, req)

	assertEmptyContentTypeHeader(t, recorder)
	assertEmptyResponseBody(t, recorder)
}

func assertEmptyResponseBody(t *testing.T, recorder *httptest.ResponseRecorder) {
	t.Helper()

	responseBodyValue := recorder.Body.String()
	if responseBodyValue != "" {
		t.Errorf("unexpected response body value: %s", responseBodyValue)
	}
}

func assertEmptyContentTypeHeader(t *testing.T, recorder *httptest.ResponseRecorder) {
	t.Helper()

	contentTypeHeaderValue := recorder.Header().Get("Content-Type")
	if contentTypeHeaderValue != "" {
		t.Errorf("unexpected header value: %s", contentTypeHeaderValue)
	}
}

func assertResponseStatus(t *testing.T, resp *httptest.ResponseRecorder, expected int) {
	t.Helper()

	if resp.Code != expected {
		t.Errorf("invalid resonse status [%d] was expecting [%d]", resp.Code, expected)
	}
}

func assertResponseHeader(t *testing.T, resp *httptest.ResponseRecorder, key, expected string) {
	t.Helper()

	if resp.Header().Get(key) != expected {
		t.Errorf("invalid header value [%s] was expecting [%s]", resp.Header().Get(key), expected)
	}
}

func assertResponseBody(t *testing.T, resp *httptest.ResponseRecorder, expected string) {
	t.Helper()

	if resp.Body.String() != expected {
		t.Errorf("invalid response value [%s] was expecting [%s]", resp.Body.String(), expected)
	}
}
