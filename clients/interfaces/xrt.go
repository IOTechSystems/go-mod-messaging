// Copyright (C) 2023 IOTech Ltd

package interfaces

import (
	"context"
	"time"

	"github.com/edgexfoundry/go-mod-core-contracts/v2/errors"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/models"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/xrtmodels"
)

// XrtClient defines the interface for interactions with the XRT MQTT Management API.
type XrtClient interface {
	AllDevices(ctx context.Context) ([]string, errors.EdgeX)
	DeviceByName(ctx context.Context, name string) (xrtmodels.DeviceInfo, errors.EdgeX)
	AddDevice(ctx context.Context, device models.Device) errors.EdgeX
	UpdateDevice(ctx context.Context, device models.Device) errors.EdgeX
	DeleteDeviceByName(ctx context.Context, name string) errors.EdgeX
	AddDiscoveredDevice(ctx context.Context, device models.Device) errors.EdgeX
	ScanDevice(ctx context.Context, device models.Device) errors.EdgeX

	ReadDeviceResources(ctx context.Context, deviceName string, resourceNames []string) (map[string]xrtmodels.Reading, errors.EdgeX)
	WriteDeviceResources(ctx context.Context, deviceName string, resourceValuePairs, options map[string]interface{}) errors.EdgeX

	AllSchedules(ctx context.Context) ([]string, errors.EdgeX)
	AddSchedule(ctx context.Context, schedule xrtmodels.Schedule) errors.EdgeX
	DeleteScheduleByName(ctx context.Context, scheduleName string) errors.EdgeX

	AllDeviceProfiles(ctx context.Context) ([]string, errors.EdgeX)
	DeviceProfileByName(ctx context.Context, name string) (models.DeviceProfile, errors.EdgeX)
	AddDeviceProfile(ctx context.Context, device models.DeviceProfile) errors.EdgeX
	UpdateDeviceProfile(ctx context.Context, device models.DeviceProfile) errors.EdgeX
	DeleteDeviceProfileByName(ctx context.Context, name string) errors.EdgeX

	UpdateLuaScript(ctx context.Context, luaScript string) errors.EdgeX

	TriggerDiscovery(ctx context.Context) errors.EdgeX

	// SetResponseTimeout sets responseTimeout to XrtClient
	SetResponseTimeout(responseTimeout time.Duration)

	// Close closes the connection of XRT client
	Close() errors.EdgeX
}
