package httpapi

import (
	"net/http"

	"buc/internal/app"
)

func Router(a *app.App) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/health", HealthHandler(a))
	mux.HandleFunc("/api/log/stream", LogStreamHandler())
	mux.HandleFunc("/api/panel/config", PanelConfigHandler(a))
	mux.HandleFunc("/api/panel/weather", PanelWeatherHandler(a))
	mux.HandleFunc("/api/panel/control", PanelControlHandler(a))
	mux.HandleFunc("/api/panel/debug/last-seen", PanelDebugLastSeenHandler())
	mux.HandleFunc("/log/live", LiveLogPageHandler())
	mux.HandleFunc("/log/files", LogFilesIndexHandler())
	mux.HandleFunc("/log/files/", LogFileHandler())
	mux.HandleFunc("/api/source/", SourceHandler(a))
	mux.HandleFunc("/api/screen/", ScreenHandler(a))
	mux.HandleFunc("/api/device/", DeviceHandler(a))
	return mux
}
