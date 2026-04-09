package ui

import (
	"buc/internal/app"
	"buc/internal/config"
)

type DeviceModel struct {
	Device DeviceMeta   `json:"device"`
	Screen *ScreenModel `json:"screen"`
}

type DeviceMeta struct {
	Name           string                  `json:"name"`
	Mode           string                  `json:"mode"`
	Screen         string                  `json:"screen"`
	Theme          string                  `json:"theme"`
	Orientation    string                  `json:"orientation"`
	Resolution     config.ResolutionConfig `json:"resolution"`
	RefreshSeconds int                     `json:"refresh_seconds"`
}

func BuildDevice(a *app.App, deviceName string) (*DeviceModel, error) {
	devCfg, ok := a.Config.Devices.Devices[deviceName]
	if !ok {
		return nil, ErrUnknownDevice(deviceName)
	}

	screenModel, err := BuildScreen(a, devCfg.Screen, devCfg.Theme, devCfg.Mode)
	if err != nil {
		return nil, err
	}

	return &DeviceModel{
		Device: DeviceMeta{
			Name:           deviceName,
			Mode:           devCfg.Mode,
			Screen:         devCfg.Screen,
			Theme:          devCfg.Theme,
			Orientation:    devCfg.Orientation,
			Resolution:     devCfg.Resolution,
			RefreshSeconds: devCfg.RefreshSeconds,
		},
		Screen: screenModel,
	}, nil
}

func ErrUnknownDevice(name string) error {
	return screenError("unknown device: " + name)
}
