// Copyright (C) 2023 IOTech Ltd

package xpert

import (
	"context"

	"github.com/edgexfoundry/go-mod-core-contracts/v2/errors"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/xrtmodels"
)

const (
	luaTransformComponent = "lua"
	componentConfigScript = "Script"
)

func (c *xrtClient) UpdateLuaScript(ctx context.Context, luaScript string) errors.EdgeX {
	config := map[string]interface{}{
		componentConfigScript: luaScript,
	}
	request := xrtmodels.NewComponentUpdateRequest(luaTransformComponent, clientName, config)
	var response xrtmodels.CommonResponse

	err := c.sendXrtCommandRequest(ctx, request.RequestId, request, &response)
	if err != nil {
		return errors.NewCommonEdgeX(errors.Kind(err), "failed to update the Lua script to Lua transform component", err)
	}
	return nil
}
