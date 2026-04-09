package httpapi

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"buc/internal/app"
	"buc/internal/config"
	"buc/internal/support"
)

const (
	panelControlActionActivate = "activate"
	panelControlActionToggle   = "toggle"
)

type panelControlRequest struct {
	Target string `json:"target"`
	Action string `json:"action"`
}

type panelControlCommand struct {
	Domain   string
	Service  string
	EntityID string
}

func PanelControlHandler(a *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		recordPanelRequest(r)

		if r.Method != http.MethodPost {
			support.JSON(w, http.StatusMethodNotAllowed, map[string]interface{}{
				"error": "method not allowed",
			})
			return
		}

		if a == nil || a.HAClient == nil {
			support.JSON(w, http.StatusServiceUnavailable, map[string]interface{}{
				"error": "ha client not configured",
			})
			return
		}

		var req panelControlRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			support.JSON(w, http.StatusBadRequest, map[string]interface{}{
				"error": "invalid json body",
			})
			return
		}

		req.Target = strings.TrimSpace(req.Target)
		req.Action = strings.TrimSpace(req.Action)

		commandKey := req.Target + ":" + req.Action
		command, ok := panelControlCommandForKey(a, commandKey)
		if !ok {
			log.Printf("panel control rejected target=%s action=%s", req.Target, req.Action)
			support.JSON(w, http.StatusBadRequest, map[string]interface{}{
				"error":  "unsupported control command",
				"target": req.Target,
				"action": req.Action,
			})
			return
		}

		log.Printf("panel control request target=%s action=%s -> %s.%s %s",
			req.Target, req.Action, command.Domain, command.Service, command.EntityID)

		if err := a.HAClient.CallService(command.Domain, command.Service, map[string]interface{}{
			"entity_id": command.EntityID,
		}); err != nil {
			log.Printf("panel control failed target=%s action=%s err=%v", req.Target, req.Action, err)
			support.JSON(w, http.StatusBadGateway, map[string]interface{}{
				"ok":      false,
				"error":   "ha service call failed",
				"target":  req.Target,
				"action":  req.Action,
				"details": err.Error(),
			})
			return
		}

		log.Printf("panel control ok target=%s action=%s", req.Target, req.Action)
		support.JSON(w, http.StatusOK, map[string]interface{}{
			"ok":     true,
			"target": req.Target,
			"action": req.Action,
		})
	}
}

func panelControlCommandForKey(a *app.App, key string) (panelControlCommand, bool) {
	if a == nil || a.Config == nil {
		return panelControlCommand{}, false
	}

	cfg, ok := a.Config.PanelCommands.Commands[key]
	if !ok {
		return panelControlCommand{}, false
	}

	return panelControlCommandFromConfig(cfg), true
}

func panelControlCommandFromConfig(cfg config.PanelCommandConfig) panelControlCommand {
	return panelControlCommand{
		Domain:   cfg.Domain,
		Service:  cfg.Service,
		EntityID: cfg.EntityID,
	}
}
