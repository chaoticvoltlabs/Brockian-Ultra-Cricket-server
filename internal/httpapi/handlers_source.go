package httpapi

import (
	"net/http"
	"strings"
	"time"

	"buc/internal/app"
	"buc/internal/ha"
	"buc/internal/support"
)

func SourceHandler(a *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := strings.TrimPrefix(r.URL.Path, "/api/source/")
		srcCfg, ok := a.Config.Sources.Sources[name]
		if !ok {
			support.JSON(w, http.StatusNotFound, map[string]interface{}{
				"error":        "unknown source",
				"source":       name,
				"generated_at": time.Now().Format(time.RFC3339),
			})
			return
		}

		entity, err := a.HAClient.GetState(srcCfg.EntityID)
		if err != nil {
			support.JSON(w, http.StatusOK, map[string]interface{}{
				"source":       name,
				"type":         srcCfg.Type,
				"entity_id":    srcCfg.EntityID,
				"generated_at": time.Now().Format(time.RFC3339),
				"data":         nil,
				"status": map[string]interface{}{
					"ok":       false,
					"warnings": []string{err.Error()},
				},
			})
			return
		}

		result, err := ha.BuildSource(name, srcCfg.EntityID, entity)
		if err != nil {
			support.JSON(w, http.StatusInternalServerError, map[string]interface{}{
				"error":        err.Error(),
				"generated_at": time.Now().Format(time.RFC3339),
			})
			return
		}

		support.JSON(w, http.StatusOK, map[string]interface{}{
			"source":       result.Name,
			"type":         result.Type,
			"entity_id":    result.EntityID,
			"generated_at": time.Now().Format(time.RFC3339),
			"data":         result.Data,
			"status":       result.Status,
		})
	}
}
