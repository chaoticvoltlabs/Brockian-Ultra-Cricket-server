package ui

import "buc/internal/config"

func BuildClimateOverview(
	_ *config.AllConfig,
	_ string,
	componentName string,
	componentType string,
	sourceName string,
	data map[string]interface{},
	options map[string]interface{},
) ComponentEnvelope {
	env := ComponentEnvelope{
		Component: componentName,
		Type:      componentType,
		Source:    sourceName,
		Status: Status{
			OK:       true,
			Warnings: []string{},
		},
		Data:     data,
		Options:  options,
		Resolved: map[string]interface{}{},
	}

	if err := RequireSourceData(sourceName, data); err != nil {
		env.Status.OK = false
		env.Status.Warnings = append(env.Status.Warnings, err.Error())
		env.Data = map[string]interface{}{}
	}

	return env
}
