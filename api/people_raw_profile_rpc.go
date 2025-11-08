package api

import (
	// INJECT IMPORTS
	"context"

	"connectrpc.com/connect"
	talv1 "github.com/TALPlatform/tal_api/proto_gen/tal/v1"
)

func (api *Api) RawProfileList(ctx context.Context, req *connect.Request[talv1.RawProfileListRequest]) (*connect.Response[talv1.RawProfileListResponse], error) {
	resp, err := api.peopleUsecase.RawProfileList(ctx, req)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}
