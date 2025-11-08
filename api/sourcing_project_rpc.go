package api

import (
        // INJECT IMPORTS
        "context"
	"connectrpc.com/connect"
        "github.com/TALPlatform/tal_api/proto_gen/tal/v1"
)


func (api *Api) ProjectInputList(ctx context.Context, req *connect.Request[talv1.ProjectInputListRequest]) (*connect.Response[talv1.ProjectInputListResponse], error) {
	resp, err := api.sourcingUsecase.ProjectInputList(ctx, req)
	if err != nil {
	     return nil, err
	}
	return connect.NewResponse(resp), nil
}

