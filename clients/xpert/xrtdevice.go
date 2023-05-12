// Copyright (C) 2023 IOTech Ltd

package xpert

import (
	"context"
	"fmt"

	"github.com/edgexfoundry/go-mod-core-contracts/v2/errors"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/models"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/xrtmodels"
)

func (c *xrtClient) AllDevices(ctx context.Context) ([]string, errors.EdgeX) {
	request := xrtmodels.NewAllDevicesRequest(clientName)
	var response xrtmodels.MultiDevicesResponse

	err := c.sendXrtRequest(ctx, request.RequestId, request, &response)
	if err != nil {
		return nil, errors.NewCommonEdgeX(errors.Kind(err), "failed to query device list", err)
	}
	return response.Result.Devices, nil
}

func (c *xrtClient) DeviceByName(ctx context.Context, name string) (xrtmodels.DeviceInfo, errors.EdgeX) {
	request := xrtmodels.NewDeviceGetRequest(name, clientName)
	var response xrtmodels.DeviceResponse

	err := c.sendXrtRequest(ctx, request.RequestId, request, &response)
	if err != nil {
		return xrtmodels.DeviceInfo{}, errors.NewCommonEdgeX(errors.Kind(err), "failed to query device", err)
	}
	return response.Result.Device, nil
}

func (c *xrtClient) AddDevice(ctx context.Context, device models.Device) errors.EdgeX {
	xrtDevice, err := xrtmodels.ToXrtDevice(device)
	if err != nil {
		return errors.NewCommonEdgeX(errors.KindServerError, "failed to convert Edgex device to XRT device data", err)
	}
	request := xrtmodels.NewDeviceAddRequest(xrtDevice, clientName)
	var response xrtmodels.CommonResponse

	err = c.sendXrtRequest(ctx, request.RequestId, request, &response)
	if err != nil {
		return errors.NewCommonEdgeX(errors.Kind(err), "failed to add device", err)
	}
	return nil
}

func (c *xrtClient) UpdateDevice(ctx context.Context, device models.Device) errors.EdgeX {
	xrtDevice, err := xrtmodels.ToXrtDevice(device)
	if err != nil {
		return errors.NewCommonEdgeX(errors.KindServerError, "failed to convert Edgex device to XRT device data", err)
	}
	request := xrtmodels.NewDeviceUpdateRequest(xrtDevice, clientName)
	var response xrtmodels.CommonResponse

	err = c.sendXrtRequest(ctx, request.RequestId, request, &response)
	if err != nil {
		return errors.NewCommonEdgeX(errors.Kind(err), "failed to update device", err)
	}
	return nil
}

func (c *xrtClient) DeleteDeviceByName(ctx context.Context, name string) errors.EdgeX {
	request := xrtmodels.NewDeviceDeleteRequest(name, clientName)
	var response xrtmodels.CommonResponse

	err := c.sendXrtRequest(ctx, request.RequestId, request, &response)
	if err != nil {
		return errors.NewCommonEdgeX(errors.Kind(err), fmt.Sprintf("failed to delete device %s", name), err)
	}
	return nil
}

// AddDiscoveredDevice adds discovered device without profile, which means the device is not usable until the profile is set or generate by device:scan
func (c *xrtClient) AddDiscoveredDevice(ctx context.Context, device models.Device) errors.EdgeX {
	xrtDevice, err := xrtmodels.ToXrtDevice(device)
	if err != nil {
		return errors.NewCommonEdgeX(errors.KindServerError, "failed to convert Edgex device to XRT device data", err)
	}
	request := xrtmodels.NewDiscoveredDeviceAddRequest(xrtDevice, clientName)
	var response xrtmodels.CommonResponse

	err = c.sendXrtRequest(ctx, request.RequestId, request, &response)
	if err != nil {
		return errors.NewCommonEdgeX(errors.Kind(err), "failed to add discovered device", err)
	}
	return nil
}

// ScanDevice checks a device profile for updates.
func (c *xrtClient) ScanDevice(ctx context.Context, device models.Device) errors.EdgeX {
	xrtDevice, err := xrtmodels.ToXrtDevice(device)
	if err != nil {
		return errors.NewCommonEdgeX(errors.KindServerError, "failed to convert Edgex device to XRT device data", err)
	}
	request := xrtmodels.NewDeviceScanRequest(xrtDevice, clientName)
	var response xrtmodels.CommonResponse

	// use discovery request for auto-generate or updating the profile
	err = c.sendXrtDiscoveryRequest(ctx, request.RequestId, request, &response)
	if err != nil {
		return errors.NewCommonEdgeX(errors.Kind(err), "failed to scan device", err)
	}
	return nil
}

func (c *xrtClient) ReadDeviceResources(ctx context.Context, deviceName string, resourceNames []string) (map[string]xrtmodels.Reading, errors.EdgeX) {
	request := xrtmodels.NewDeviceResourceGetRequest(deviceName, clientName, resourceNames)
	var response xrtmodels.MultiResourcesResponse

	err := c.sendXrtRequest(ctx, request.RequestId, request, &response)
	if err != nil {
		return nil, errors.NewCommonEdgeX(errors.Kind(err), "failed to read device resources", err)
	}
	return response.Result.Readings, nil
}

func (c *xrtClient) WriteDeviceResources(ctx context.Context, deviceName string, resourceValuePairs, options map[string]interface{}) errors.EdgeX {
	request := xrtmodels.NewDeviceResourceSetRequest(deviceName, clientName, resourceValuePairs, options)
	var response xrtmodels.CommonResponse

	err := c.sendXrtRequest(ctx, request.RequestId, request, &response)
	if err != nil {
		return errors.NewCommonEdgeX(errors.Kind(err), "failed to write device resources", err)
	}
	return nil
}
