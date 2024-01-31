// Copyright (C) 2023 IOTech Ltd

package central

import (
	"context"

	"github.com/edgexfoundry/go-mod-core-contracts/v3/errors"
	"github.com/edgexfoundry/go-mod-core-contracts/v3/xrtmodels"
)

const (
	discoveryDurationOption = "DiscoveryDuration"
)

func (c *xrtClient) TriggerDiscovery(ctx context.Context) errors.EdgeX {
	if c.clientOptions == nil || c.clientOptions.DiscoveryOptions == nil {
		return errors.NewCommonEdgeX(errors.KindContractInvalid, "please provide DiscoveryOptions for the discovery request", nil)
	}
	options := map[string]interface{}{
		discoveryDurationOption: c.clientOptions.DiscoveryOptions.DiscoveryDuration.Milliseconds()}
	request := xrtmodels.NewDiscoveryRequest(clientName, options)
	var response xrtmodels.CommonResponse

	err := c.sendXrtDiscoveryRequest(ctx, request.RequestId, request, &response)
	if err != nil {
		return errors.NewCommonEdgeX(errors.Kind(err), "failed to trigger discovery", err)
	}
	return nil
}
