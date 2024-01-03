// Copyright (C) 2023 IOTech Ltd

package central

import (
	"context"
	"fmt"

	"github.com/edgexfoundry/go-mod-core-contracts/v3/errors"
	"github.com/edgexfoundry/go-mod-core-contracts/v3/xrtmodels"
)

func (c *xrtClient) AllSchedules(ctx context.Context) ([]string, errors.EdgeX) {
	request := xrtmodels.NewAllSchedulesRequest(clientName)
	var response xrtmodels.MultiSchedulesResponse

	err := c.sendXrtRequest(ctx, request.RequestId, request, &response)
	if err != nil {
		return nil, errors.NewCommonEdgeX(errors.Kind(err), "failed to query schedule list", err)
	}
	return response.Result.Schedules, nil
}

func (c *xrtClient) AddSchedule(ctx context.Context, schedule xrtmodels.Schedule) errors.EdgeX {
	request := xrtmodels.NewScheduleAddRequest(clientName, schedule)
	var response xrtmodels.CommonResponse

	err := c.sendXrtRequest(ctx, request.RequestId, request, &response)
	if err != nil {
		return errors.NewCommonEdgeX(errors.Kind(err), "failed to add schedule", err)
	}
	return nil
}

func (c *xrtClient) DeleteScheduleByName(ctx context.Context, name string) errors.EdgeX {
	request := xrtmodels.NewScheduleDeleteRequest(name, clientName)
	var response xrtmodels.CommonResponse

	err := c.sendXrtRequest(ctx, request.RequestId, request, &response)
	if err != nil {
		return errors.NewCommonEdgeX(errors.Kind(err), fmt.Sprintf("failed to delete schedule %s", name), err)
	}
	return nil
}
