package api

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	talv1 "github.com/TALPlatform/tal_api/proto_gen/tal/v1"
)

func (api *Api) CommandPalleteSearch(ctx context.Context, req *connect.Request[talv1.CommandPalleteSearchRequest]) (*connect.Response[talv1.CommandPalleteSearchResponse], error) {
	resp ,  err := api.publicUsecase.CommandPalleteSearch(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve users list: %w", err)
	}
	return connect.NewResponse(resp), nil
}


func (api *Api) CommandPalleteSync(ctx context.Context, req *connect.Request[talv1.CommandPalleteSyncRequest]) (*connect.Response[talv1.CommandPalleteSyncResponse], error) {
	resp ,  err := api.publicUsecase.CommandPalleteSync(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve users list: %w", err)
	}
	return connect.NewResponse(resp), nil
}
