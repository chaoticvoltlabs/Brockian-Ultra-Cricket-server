package main

import (
	"log"
	"net/http"
	"os"

	"buc/internal/app"
	"buc/internal/config"
	"buc/internal/ha"
	"buc/internal/httpapi"
	"buc/internal/logview"
)

func main() {
	configDir := getenv("UX_CONFIG_DIR", "./config")
	haURL := getenv("HA_BASE_URL", "http://127.0.0.1:8123")
	haToken := os.Getenv("HA_TOKEN")
	listenAddr := getenv("UX_LISTEN_ADDR", "127.0.0.1:9100")
	logDir := getenv("BUC_LOG_DIR", "./logs")

	logManager, err := logview.NewManager(logDir)
	if err != nil {
		log.Fatalf("log init failed: %v", err)
	}
	logview.SetDefault(logManager)
	log.SetOutput(logManager)

	cfg, err := config.LoadAll(configDir)
	if err != nil {
		log.Fatalf("config load failed: %v", err)
	}

	haClient := ha.NewClient(haURL, haToken)
	application := app.New(cfg, haClient)

	router := httpapi.Router(application)

	log.Printf("buc-server listening on %s", listenAddr)
	if err := http.ListenAndServe(listenAddr, router); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

func getenv(key, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return v
}
