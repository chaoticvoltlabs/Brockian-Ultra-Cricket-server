package ui

import (
	"fmt"

	"buc/internal/config"
)

func BuildComponent(
	cfg *config.AllConfig,
	themeName string,
	componentName string,
	componentCfg config.ComponentConfig,
	sourceData map[string]interface{},
) (ComponentEnvelope, error) {
	switch componentCfg.Type {
	case "outside_summary":
		return BuildOutsideSummary(cfg, themeName, componentName, componentCfg.Type, componentCfg.Source, sourceData, componentCfg.Options), nil
	case "wind_strip":
		return BuildWindStrip(cfg, themeName, componentName, componentCfg.Type, componentCfg.Source, sourceData, componentCfg.Options), nil
	case "daily_forecast":
		return BuildDailyForecast(cfg, themeName, componentName, componentCfg.Type, componentCfg.Source, sourceData, componentCfg.Options), nil
	case "climate_overview":
		return BuildClimateOverview(cfg, themeName, componentName, componentCfg.Type, componentCfg.Source, sourceData, componentCfg.Options), nil
	default:
		return ComponentEnvelope{
			Component: componentName,
			Type:      componentCfg.Type,
			Source:    componentCfg.Source,
			Status: Status{
				OK:       true,
				Warnings: []string{"component type not yet resolved; passing through raw data"},
			},
			Data:     sourceData,
			Options:  componentCfg.Options,
			Resolved: map[string]interface{}{},
		}, nil
	}
}

func RequireSourceData(sourceName string, sourceData map[string]interface{}) error {
	if sourceName == "" {
		return nil
	}
	if sourceData == nil {
		return fmt.Errorf("missing source data for %s", sourceName)
	}
	return nil
}
