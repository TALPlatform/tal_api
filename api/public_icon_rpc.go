package api

import (
	"context"

	"connectrpc.com/connect"
	talv1 "github.com/TALPlatform/tal_api/proto_gen/tal/v1"
)

func (api *Api) IconCreateUpdateBulk(ctx context.Context, req *connect.Request[talv1.IconCreateUpdateBulkRequest]) (*connect.Response[talv1.IconCreateUpdateBulkResponse], error) {
	if err := ctx.Err(); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	_, err := api.publicUsecase.IconCreateUpdateBulk(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&talv1.IconCreateUpdateBulkResponse{}), nil
}

func (api *Api) IconFind(ctx context.Context, req *connect.Request[talv1.IconFindRequest]) (*connect.Response[talv1.IconFindResponse], error) {
	if err := ctx.Err(); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	resp, err := api.publicUsecase.IconFind(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}
func (api *Api) IconList(ctx context.Context, req *connect.Request[talv1.IconListRequest]) (*connect.Response[talv1.IconListResponse], error) {
	if err := ctx.Err(); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	resp, err := api.publicUsecase.IconList(ctx)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}
