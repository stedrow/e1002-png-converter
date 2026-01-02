package main

import (
	"embed"
	"encoding/json"
	"fmt"
)

//go:embed devices.json
var devicesFS embed.FS

// Device represents a display device configuration
type Device struct {
	Name        string `json:"name"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	Colors      int    `json:"colors"`
	Description string `json:"description"`
}

// DeviceConfig holds all device definitions
type DeviceConfig struct {
	Devices map[string]Device  `json:"devices"`
	Aliases map[string]string  `json:"aliases"`
}

var deviceConfig *DeviceConfig

// LoadDeviceConfig loads the embedded devices.json file
func LoadDeviceConfig() (*DeviceConfig, error) {
	if deviceConfig != nil {
		return deviceConfig, nil
	}

	data, err := devicesFS.ReadFile("devices.json")
	if err != nil {
		return nil, fmt.Errorf("failed to read devices.json: %w", err)
	}

	var config DeviceConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse devices.json: %w", err)
	}

	deviceConfig = &config
	return deviceConfig, nil
}

// GetDevice returns a device by ID, resolving aliases
func GetDevice(deviceID string) (*Device, error) {
	config, err := LoadDeviceConfig()
	if err != nil {
		return nil, err
	}

	// Check if it's an alias
	if actualID, ok := config.Aliases[deviceID]; ok {
		deviceID = actualID
	}

	// Get device
	device, ok := config.Devices[deviceID]
	if !ok {
		return nil, fmt.Errorf("unknown device: %s", deviceID)
	}

	return &device, nil
}

// ListDevices returns all available device IDs and names
func ListDevices() ([]string, error) {
	config, err := LoadDeviceConfig()
	if err != nil {
		return nil, err
	}

	var devices []string
	for id, device := range config.Devices {
		devices = append(devices, fmt.Sprintf("  %-25s %s", id, device.Name))
	}

	return devices, nil
}
