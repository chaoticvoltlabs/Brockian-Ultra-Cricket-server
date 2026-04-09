package app

import (
	"buc/internal/config"
	"buc/internal/ha"
)

type App struct {
	Config   *config.AllConfig
	HAClient *ha.Client
}

func New(cfg *config.AllConfig, hc *ha.Client) *App {
	return &App{
		Config:   cfg,
		HAClient: hc,
	}
}
