package api

import (
	// INJECT IMPORTS
	"context"

	"connectrpc.com/connect"
	talv1 "github.com/TALPlatform/tal_api/proto_gen/tal/v1"
)

func (api *Api) SectionList(ctx context.Context, req *connect.Request[talv1.SectionListRequest]) (*connect.Response[talv1.SectionListResponse], error) {
	resp, err := api.tenantUsecase.SectionList(ctx, req)
	if err != nil {
		return nil, err
	}

	resp.Options = api.getAvailableOptions(req.Header())

	return connect.NewResponse(resp), nil
}

func (api *Api) SectionFindForUpdate(ctx context.Context, req *connect.Request[talv1.SectionFindForUpdateRequest]) (*connect.Response[talv1.SectionFindForUpdateResponse], error) {
	resp, err := api.tenantUsecase.SectionFindForUpdate(ctx, req)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(resp), nil
}

func (api *Api) SectionCreateUpdate(ctx context.Context, req *connect.Request[talv1.SectionCreateUpdateRequest]) (*connect.Response[talv1.SectionCreateUpdateResponse], error) {
	resp, err := api.tenantUsecase.SectionCreateUpdate(ctx, req)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}
func (api *Api) SectionDeleteRestore(ctx context.Context, req *connect.Request[talv1.SectionDeleteRestoreRequest]) (*connect.Response[talv1.SectionDeleteRestoreResponse], error) {
	resp, err := api.tenantUsecase.SectionDeleteRestore(ctx, req)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}
