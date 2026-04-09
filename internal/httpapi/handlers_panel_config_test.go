package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"buc/internal/app"
	"buc/internal/config"
)

func loadTestConfig(t *testing.T) *config.AllConfig {
	t.Helper()

	cfg, err := config.LoadAll("../../config")
	if err != nil {
		t.Fatalf("load config: %v", err)
	}

	return cfg
}

func TestPanelConfigHandlerReturnsRoomBetaProfileForKnownMAC(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/panel/config", nil)
	req.Header.Set("X-Panel-MAC", "11:22:33:44:55:66")
	rec := httptest.NewRecorder()

	Router(app.New(loadTestConfig(t), nil)).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var resp panelConfigResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if resp.Profile != "room_beta" {
		t.Fatalf("expected room_beta profile, got %#v", resp.Profile)
	}
	if len(resp.Page3.Scenes) != 3 {
		t.Fatalf("expected 3 room_beta scenes, got %d", len(resp.Page3.Scenes))
	}
	if len(resp.Page3.Targets) != 5 {
		t.Fatalf("expected 5 room_beta targets, got %d", len(resp.Page3.Targets))
	}
	if resp.Page3.Targets[0].Target != "task_light" {
		t.Fatalf("unexpected first room_beta target: %#v", resp.Page3.Targets[0])
	}
}

func TestPanelConfigHandlerReturnsRoomAlphaProfileByDefault(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/panel/config", nil)
	rec := httptest.NewRecorder()

	Router(app.New(loadTestConfig(t), nil)).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var resp panelConfigResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if resp.Profile != "room_alpha" {
		t.Fatalf("expected room_alpha profile, got %#v", resp.Profile)
	}
	if len(resp.Page3.Scenes) != 4 {
		t.Fatalf("expected 4 room_alpha scenes, got %d", len(resp.Page3.Scenes))
	}
}

func TestPanelConfigHandlerAcceptsQueryParamFallbackForBrowserTesting(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/panel/config?panel_mac=11:22:33:44:55:66&panel_ip=192.0.2.42", nil)
	rec := httptest.NewRecorder()

	Router(app.New(loadTestConfig(t), nil)).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var resp panelConfigResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if resp.Profile != "room_beta" {
		t.Fatalf("expected room_beta profile via query fallback, got %#v", resp.Profile)
	}
}
