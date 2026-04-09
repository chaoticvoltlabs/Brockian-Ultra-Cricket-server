package httpapi

import (
	"net/http"
	"strings"

	"buc/internal/app"
	"buc/internal/support"
	"buc/internal/ui"
)

func DeviceHandler(a *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := strings.TrimPrefix(r.URL.Path, "/api/device/")

		model, err := ui.BuildDevice(a, name)
		if err != nil {
			support.JSON(w, http.StatusNotFound, map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		support.JSON(w, http.StatusOK, model)
	}
}
