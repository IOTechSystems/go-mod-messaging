// Copyright (C) 2023 IOTech Ltd

package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/edgexfoundry/go-mod-core-contracts/v2/errors"
)

func FetchXRTResponse(ctx context.Context, requestId string, requestMap RequestMap, responseTimeout time.Duration) ([]byte, errors.EdgeX) {
	resChan, ok := requestMap.Get(requestId)
	if !ok {
		return nil, errors.NewCommonEdgeX(errors.KindServerError, fmt.Sprintf("the corresponding ResponseChan not found by requestId %s", requestId), nil)
	}
	defer func() {
		close(resChan)
		requestMap.Delete(requestId)
	}()

	timeout := time.After(responseTimeout)
	select {
	case <-ctx.Done():
		return nil, nil
	case <-timeout:
		return nil, errors.NewCommonEdgeX(errors.KindServerError, "timed out fetching command response", nil)
	case commandResponse := <-resChan:
		return commandResponse, nil
	}
}
