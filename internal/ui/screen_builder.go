package ui

import (
	"time"

	"buc/internal/app"
	"buc/internal/ha"
	"buc/internal/theme"
)

func BuildScreen(a *app.App, screenName, themeName, mode string) (*ScreenModel, error) {
	screenCfg, ok := a.Config.Screens.Screens[screenName]
	if !ok {
		return nil, ErrUnknownScreen(screenName)
	}

	if themeName == "" {
		themeName = "dark_default"
	}

	resolvedTheme, err := theme.ResolveTheme(a.Config, themeName)
	if err != nil {
		return nil, err
	}

	model := &ScreenModel{
		Screen: ScreenMeta{
			Name:        screenName,
			Title:       screenCfg.Title,
			Layout:      screenCfg.Layout,
			Theme:       themeName,
			GeneratedAt: time.Now().Format(time.RFC3339),
			DeviceMode:  mode,
		},
		Layout: LayoutMeta{
			Regions: orderedRegionNames(screenCfg.Regions),
		},
		Theme: ThemeMeta{
			Name:   resolvedTheme.Name,
			Tokens: resolvedTheme.Tokens,
		},
		Regions: map[string][]ComponentEnvelope{},
	}

	for regionName, componentNames := range screenCfg.Regions {
		model.Regions[regionName] = []ComponentEnvelope{}

		for _, componentName := range componentNames {
			componentCfg := a.Config.Components.Components[componentName]

			var sourceData map[string]interface{}
			if componentCfg.Source != "" {
				srcCfg := a.Config.Sources.Sources[componentCfg.Source]
				entity, err := a.HAClient.GetState(srcCfg.EntityID)
				if err != nil {
					env := ComponentEnvelope{
						Component: componentName,
						Type:      componentCfg.Type,
						Source:    componentCfg.Source,
						Status: Status{
							OK:       false,
							Warnings: []string{err.Error()},
						},
						Data:     map[string]interface{}{},
						Options:  componentCfg.Options,
						Resolved: map[string]interface{}{},
					}
					model.Regions[regionName] = append(model.Regions[regionName], env)
					continue
				}

				srcResult, err := ha.BuildSource(componentCfg.Source, srcCfg.EntityID, entity)
				if err != nil {
					env := ComponentEnvelope{
						Component: componentName,
						Type:      componentCfg.Type,
						Source:    componentCfg.Source,
						Status: Status{
							OK:       false,
							Warnings: []string{err.Error()},
						},
						Data:     map[string]interface{}{},
						Options:  componentCfg.Options,
						Resolved: map[string]interface{}{},
					}
					model.Regions[regionName] = append(model.Regions[regionName], env)
					continue
				}
				sourceData = srcResult.Data
			}

			env, err := BuildComponent(a.Config, themeName, componentName, componentCfg, sourceData)
			if err != nil {
				env = ComponentEnvelope{
					Component: componentName,
					Type:      componentCfg.Type,
					Source:    componentCfg.Source,
					Status: Status{
						OK:       false,
						Warnings: []string{err.Error()},
					},
					Data:     map[string]interface{}{},
					Options:  componentCfg.Options,
					Resolved: map[string]interface{}{},
				}
			}

			model.Regions[regionName] = append(model.Regions[regionName], env)
		}
	}

	return model, nil
}

func orderedRegionNames(regions map[string][]string) []string {
	order := []string{
		"left", "right_top", "right_bottom", "footer",
		"left_top", "left_middle", "right_top", "right_middle", "footer",
	}

	exists := map[string]bool{}
	for name := range regions {
		exists[name] = true
	}

	out := []string{}
	for _, name := range order {
		if exists[name] {
			out = append(out, name)
			delete(exists, name)
		}
	}
	for name := range exists {
		out = append(out, name)
	}
	return out
}

type screenError string

func (e screenError) Error() string { return string(e) }

func ErrUnknownScreen(name string) error {
	return screenError("unknown screen: " + name)
}
