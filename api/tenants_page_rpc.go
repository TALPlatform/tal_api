package api

import (
	// INJECT IMPORTS
	"context"

	"connectrpc.com/connect"
	talv1 "github.com/TALPlatform/tal_api/proto_gen/tal/v1"
)

func (api *Api) PageList(ctx context.Context, req *connect.Request[talv1.PageListRequest]) (*connect.Response[talv1.PageListResponse], error) {
	resp, err := api.tenantUsecase.PageList(ctx, req)

	if err != nil {
		return nil, err
	}

	resp.Options = api.getAvailableOptions(req.Header())
	return connect.NewResponse(resp), nil
}

func (api *Api) PageCreateUpdate(ctx context.Context, req *connect.Request[talv1.PageCreateUpdateRequest]) (*connect.Response[talv1.PageCreateUpdateResponse], error) {
	resp, err := api.tenantUsecase.PageCreateUpdate(ctx, req)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}
func (api *Api) PageDeleteRestore(ctx context.Context, req *connect.Request[talv1.PageDeleteRestoreRequest]) (*connect.Response[talv1.PageDeleteRestoreResponse], error) {
	resp, err := api.tenantUsecase.PageDeleteRestore(ctx, req)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}
func (api *Api) PageFindForUpdate(ctx context.Context, req *connect.Request[talv1.PageFindForUpdateRequest]) (*connect.Response[talv1.PageFindForUpdateResponse], error) {
	resp, err := api.tenantUsecase.PageFindForUpdate(ctx, req)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}
