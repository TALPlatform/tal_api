package api

import (
	// INJECT IMPORTS
	"context"

	"connectrpc.com/connect"
	"github.com/TALPlatform/tal_api/pkg/headerkeys"
	talv1 "github.com/TALPlatform/tal_api/proto_gen/tal/v1"
)

func (api *Api) SectionListInpt(ctx context.Context, req *connect.Request[talv1.SectionListInptRequest]) (*connect.Response[talv1.SectionListInptResponse], error) {
	resp, err := api.tenantUsecase.SectionListInpt(ctx, req)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

func (api *Api) TenantCreateUpdate(ctx context.Context, req *connect.Request[talv1.TenantCreateUpdateRequest]) (*connect.Response[talv1.TenantCreateUpdateResponse], error) {
	resp, err := api.tenantUsecase.TenantCreateUpdate(ctx, req)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}
func (api *Api) TenantList(ctx context.Context, req *connect.Request[talv1.TenantListRequest]) (*connect.Response[talv1.TenantListResponse], error) {
	resp, err := api.tenantUsecase.TenantList(ctx, req)
	if err != nil {
		return nil, err
	}
	resp.Options = api.getAvailableOptions(req.Header())
	if resp.Options.UpdateHandler != nil {
		resp.Options.UpdateHandler.FindEndpoint = "tenantFind"
		resp.Options.UpdateHandler.FindRequestProperty = "tenantId"
		resp.Options.UpdateHandler.FindResponseProperty = "tenant"
	}
	return connect.NewResponse(resp), nil

}
func (api *Api) TenantListInput(ctx context.Context, req *connect.Request[talv1.TenantListInputRequest]) (*connect.Response[talv1.TenantListInputResponse], error) {
	resp, err := api.tenantUsecase.TenantListInput(ctx, req)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil

}

func (api *Api) TenantFind(ctx context.Context, req *connect.Request[talv1.TenantFindRequest]) (*connect.Response[talv1.TenantFindResponse], error) {
	resp, err := api.tenantUsecase.TenantFind(ctx, req)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil

}
func (api *Api) TenantDeleteRestore(ctx context.Context, req *connect.Request[talv1.TenantDeleteRestoreRequest]) (*connect.Response[talv1.TenantDeleteRestoreResponse], error) {
	resp, err := api.tenantUsecase.TenantDeleteRestore(ctx, req)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil

}

func (api *Api) TenantDashboard(ctx context.Context, req *connect.Request[talv1.TenantDashboardRequest]) (*connect.Response[talv1.TenantDashboardResponse], error) {
	resp, err := api.tenantUsecase.TenantDashboard(ctx, req)
	if err != nil {
		return nil, err
	}
	pageListReq := connect.NewRequest(&talv1.PageListRequest{})
	permissGroup, err := api.CheckForAccess(ctx, "PageList", "page")
	if err != nil {
		return nil, err
	}

	// pageListReq.Header().Add(key string, value string)
	headerkeys.WithPermissionGroup(pageListReq.Header(), "page")
	headerkeys.WithPermittedActions(pageListReq.Header(), *permissGroup)
	pagesResponse, err := api.PageList(ctx, pageListReq)
	if err != nil {
		return nil, err
	}
	resp.Pages = pagesResponse.Msg
	return connect.NewResponse(resp), nil
}
func (api *Api) PartialTypeListInput(ctx context.Context, req *connect.Request[talv1.PartialTypeListInputRequest]) (*connect.Response[talv1.PartialTypeListInputResponse], error) {
	resp, err := api.tenantUsecase.PartialTypeListInput(ctx, req)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}
