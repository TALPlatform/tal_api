package api

import (
        // INJECT IMPORTS
        "context"
	"connectrpc.com/connect"
        "github.com/TALPlatform/tal_api/proto_gen/tal/v1"
)


func (api *Api) SourcingSessionCreateUpdate(ctx context.Context, req *connect.Request[talv1.SourcingSessionCreateUpdateRequest]) (*connect.Response[talv1.SourcingSessionCreateUpdateResponse], error) {
	resp, err := api.sourcingUsecase.SourcingSessionCreateUpdate(ctx, req)
	if err != nil {
	     return nil, err
	}
	return connect.NewResponse(resp), nil
}

func (api *Api) SourcingSessionFind(ctx context.Context, req *connect.Request[talv1.SourcingSessionFindRequest]) (*connect.Response[talv1.SourcingSessionFindResponse], error) {
	resp, err := api.sourcingUsecase.SourcingSessionFind(ctx, req)
	if err != nil {
	     return nil, err
	}
	return connect.NewResponse(resp), nil
}

func (api *Api) SourcingSessionProfileCreateUpdate(ctx context.Context, req *connect.Request[talv1.SourcingSessionProfileCreateUpdateRequest]) (*connect.Response[talv1.SourcingSessionProfileCreateUpdateResponse], error) {
	resp, err := api.sourcingUsecase.SourcingSessionProfileCreateUpdate(ctx, req)
	if err != nil {
	     return nil, err
	}
	return connect.NewResponse(resp), nil
}

