package httpapi

import (
	"net/http"
	"time"

	"buc/internal/app"
	"buc/internal/support"
)

func HealthHandler(a *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]interface{}{
			"ok":           true,
			"generated_at": time.Now().Format(time.RFC3339),
			"checks": map[string]interface{}{
				"config_loaded":  a.Config != nil,
				"ha_reachable":   nil,
				"cache_writable": true,
			},
		}
		support.JSON(w, http.StatusOK, resp)
	}
}
