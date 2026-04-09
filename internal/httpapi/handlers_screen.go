package httpapi

import (
	"net/http"
	"strings"

	"buc/internal/app"
	"buc/internal/support"
	"buc/internal/ui"
)

func ScreenHandler(a *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := strings.TrimPrefix(r.URL.Path, "/api/screen/")
		themeName := r.URL.Query().Get("theme")
		mode := r.URL.Query().Get("mode")
		if mode == "" {
			mode = "debug"
		}

		model, err := ui.BuildScreen(a, name, themeName, mode)
		if err != nil {
			support.JSON(w, http.StatusNotFound, map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		support.JSON(w, http.StatusOK, model)
	}
}
