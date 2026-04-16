package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"buc/internal/app"
	"buc/internal/ha"
)

func TestPanelControlHandlerCallsHALightAToggle(t *testing.T) {
	var gotMethod string
	var gotPath string
	var gotAuth string
	var gotBody map[string]interface{}

	haServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		gotAuth = r.Header.Get("Authorization")
		if err := json.NewDecoder(r.Body).Decode(&gotBody); err != nil {
			t.Fatalf("decode request body: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[]`))
	}))
	defer haServer.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/panel/control",
		strings.NewReader(`{"target":"light_a","action":"toggle"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Panel-MAC", "aa:bb:cc:dd:ee:ff")
	req.Header.Set("X-Panel-IP", "192.0.2.10")
	rec := httptest.NewRecorder()

	Router(app.New(loadTestConfig(t), ha.NewClient(haServer.URL, "test-token"))).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, rec.Code, rec.Body.String())
	}

	if gotMethod != http.MethodPost {
		t.Fatalf("expected POST to HA, got %s", gotMethod)
	}
	if gotPath != "/api/services/light/toggle" {
		t.Fatalf("expected HA path /api/services/light/toggle, got %s", gotPath)
	}
	if gotAuth != "Bearer test-token" {
		t.Fatalf("expected auth header to be set, got %q", gotAuth)
	}
	if gotBody["entity_id"] != "light.room_light_a" {
		t.Fatalf("expected entity_id %q, got %#v", "light.room_light_a", gotBody["entity_id"])
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if ok, _ := resp["ok"].(bool); !ok {
		t.Fatalf("expected ok=true, got %#v", resp)
	}
}

func TestPanelDebugLastSeenHandlerTracksPanelHeaders(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/panel/debug/last-seen", nil)
	rec := httptest.NewRecorder()

	PanelDebugLastSeenHandler().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var before map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &before); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	haServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[]`))
	}))
	defer haServer.Close()

	controlReq := httptest.NewRequest(http.MethodPost, "/api/panel/control",
		strings.NewReader(`{"target":"light_a","action":"toggle"}`))
	controlReq.Header.Set("Content-Type", "application/json")
	controlReq.Header.Set("X-Panel-MAC", "aa:bb:cc:dd:ee:ff")
	controlReq.Header.Set("X-Panel-IP", "192.0.2.10")
	controlRec := httptest.NewRecorder()

	Router(app.New(loadTestConfig(t), ha.NewClient(haServer.URL, "test-token"))).ServeHTTP(controlRec, controlReq)
	if controlRec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, controlRec.Code, controlRec.Body.String())
	}

	req = httptest.NewRequest(http.MethodGet, "/api/panel/debug/last-seen", nil)
	rec = httptest.NewRecorder()
	PanelDebugLastSeenHandler().ServeHTTP(rec, req)

	var after panelLastSeen
	if err := json.Unmarshal(rec.Body.Bytes(), &after); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if after.PanelMAC != "aa:bb:cc:dd:ee:ff" {
		t.Fatalf("expected panel MAC to be tracked, got %#v", after.PanelMAC)
	}
	if after.PanelIP != "192.0.2.10" {
		t.Fatalf("expected panel IP to be tracked, got %#v", after.PanelIP)
	}
	if after.Path != "/api/panel/control" {
		t.Fatalf("expected last path to be tracked, got %#v", after.Path)
	}
}

func TestPanelControlHandlerCallsHAWorkScene(t *testing.T) {
	var gotMethod string
	var gotPath string
	var gotBody map[string]interface{}

	haServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		if err := json.NewDecoder(r.Body).Decode(&gotBody); err != nil {
			t.Fatalf("decode request body: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[]`))
	}))
	defer haServer.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/panel/control",
		strings.NewReader(`{"target":"scene_work","action":"activate"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	Router(app.New(loadTestConfig(t), ha.NewClient(haServer.URL, "test-token"))).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, rec.Code, rec.Body.String())
	}

	if gotMethod != http.MethodPost {
		t.Fatalf("expected POST to HA, got %s", gotMethod)
	}
	if gotPath != "/api/services/scene/turn_on" {
		t.Fatalf("expected HA path /api/services/scene/turn_on, got %s", gotPath)
	}
	if gotBody["entity_id"] != "scene.room_work" {
		t.Fatalf("expected entity_id %q, got %#v", "scene.room_work", gotBody["entity_id"])
	}
}

func TestPanelControlHandlerSupportsAmbientBrightness(t *testing.T) {
	var gotBody map[string]interface{}

	haServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&gotBody); err != nil {
			t.Fatalf("decode request body: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[]`))
	}))
	defer haServer.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/panel/control",
		strings.NewReader(`{"target":"ambient","action":"set_brightness","value":42}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	Router(app.New(loadTestConfig(t), ha.NewClient(haServer.URL, "test-token"))).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, rec.Code, rec.Body.String())
	}
	if gotBody["entity_id"] != "light.room_ambient" {
		t.Fatalf("expected entity_id %q, got %#v", "light.room_ambient", gotBody["entity_id"])
	}
	if gotBody["brightness_pct"] != float64(42) {
		t.Fatalf("expected brightness_pct 42, got %#v", gotBody["brightness_pct"])
	}
}

func TestPanelControlHandlerSupportsAmbientRGB(t *testing.T) {
	var gotBody map[string]interface{}

	haServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&gotBody); err != nil {
			t.Fatalf("decode request body: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[]`))
	}))
	defer haServer.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/panel/control",
		strings.NewReader(`{"target":"ambient","action":"set_rgb","rgb":[255,136,32]}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	Router(app.New(loadTestConfig(t), ha.NewClient(haServer.URL, "test-token"))).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, rec.Code, rec.Body.String())
	}
	if gotBody["entity_id"] != "light.room_ambient" {
		t.Fatalf("expected entity_id %q, got %#v", "light.room_ambient", gotBody["entity_id"])
	}
	rgb, ok := gotBody["rgb_color"].([]interface{})
	if !ok || len(rgb) != 3 {
		t.Fatalf("expected rgb_color triplet, got %#v", gotBody["rgb_color"])
	}
}
func TestPanelControlHandlerRejectsUnsupportedCommand(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/panel/control",
		strings.NewReader(`{"target":"other","action":"toggle"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	Router(app.New(loadTestConfig(t), ha.NewClient("http://127.0.0.1:8123", "test-token"))).ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d: %s", http.StatusBadRequest, rec.Code, rec.Body.String())
	}
}
