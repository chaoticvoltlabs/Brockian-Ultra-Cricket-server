package theme

import (
	"fmt"

	"buc/internal/config"
)

type ResolvedTheme struct {
	Name   string            `json:"name"`
	Tokens map[string]string `json:"tokens"`
}

func ResolveTheme(cfg *config.AllConfig, themeName string) (*ResolvedTheme, error) {
	t, ok := cfg.Themes.Themes[themeName]
	if !ok {
		return nil, fmt.Errorf("unknown theme %q", themeName)
	}

	out := &ResolvedTheme{
		Name:   themeName,
		Tokens: map[string]string{},
	}

	for k, v := range t.Tokens {
		out.Tokens[k] = v
	}

	return out, nil
}

func ResolveTemperatureToken(cfg *config.AllConfig, themeName, scaleName string, value float64) (string, string, error) {
	themeCfg, ok := cfg.Themes.Themes[themeName]
	if !ok {
		return "", "", fmt.Errorf("unknown theme %q", themeName)
	}

	scale, ok := themeCfg.TemperatureScales[scaleName]
	if !ok {
		return "", "", fmt.Errorf("unknown temperature scale %q for theme %q", scaleName, themeName)
	}

	for _, band := range scale.Bands {
		if value <= band.Max {
			color, ok := themeCfg.Tokens[band.Token]
			if !ok {
				return "", "", fmt.Errorf("unknown token %q in theme %q", band.Token, themeName)
			}
			return band.Token, color, nil
		}
	}

	return "", "", fmt.Errorf("no matching temperature band in scale %q", scaleName)
}
