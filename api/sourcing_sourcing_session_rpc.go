package api

import (
	// INJECT IMPORTS
	"context"

	"connectrpc.com/connect"
	talv1 "github.com/TALPlatform/tal_api/proto_gen/tal/v1"
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

func (api *Api) SourcingSessionList(ctx context.Context, req *connect.Request[talv1.SourcingSessionListRequest]) (*connect.Response[talv1.SourcingSessionListResponse], error) {
	resp, err := api.sourcingUsecase.SourcingSessionList(ctx, req)
	if err != nil {
		return nil, err
	}
	resp.Options = api.getAvailableOptions(req.Header())
	return connect.NewResponse(resp), nil
}
func (api *Api) SourcingSessionProfileJustify(ctx context.Context, req *connect.Request[talv1.SourcingSessionProfileJustifyRequest], stream *connect.ServerStream[talv1.SourcingSessionProfileJustifyResponse]) error {
	err := api.sourcingUsecase.SourcingSessionProfileJustify(ctx, req, stream)
	if err != nil {
		return err
	}
	return nil
}
func (api *Api) SourcingSessionCriteriaCreate(ctx context.Context, req *connect.Request[talv1.SourcingSessionCriteriaCreateRequest], stream *connect.ServerStream[talv1.SourcingSessionCriteriaCreateResponse]) error {
	// err := api.sourcingUsecase.SourcingSessionCriteriaCreate(ctx, req, stream)
	// if err != nil {
	// 	return err
	// }
	return nil
}
func (api *Api) SourcingSessionApply(ctx context.Context, req *connect.Request[talv1.SourcingSessionApplyRequest], stream *connect.ServerStream[talv1.SourcingSessionApplyResponse]) error {
	err := api.sourcingUsecase.SourcingSessionApply(ctx, req, stream)
	if err != nil {
		return err
	}
	return nil
}
func (api *Api) SourcingSessionFiltersBuilder(ctx context.Context, req *connect.Request[talv1.SourcingSessionFiltersBuilderRequest]) (*connect.Response[talv1.SourcingSessionFiltersBuilderResponse], error) {
	resp, err := api.sourcingUsecase.SourcingSessionFiltersBuilder(ctx, req)
	if err != nil {
	     return nil, err
	}
	return connect.NewResponse(resp), nil
}

func (api *Api) SourcingSessionFiltersInfered(ctx context.Context, req *connect.Request[talv1.SourcingSessionFiltersInferedRequest]) (*connect.Response[talv1.SourcingSessionFiltersInferedResponse], error) {
	resp, err := api.sourcingUsecase.SourcingSessionFiltersInfered(ctx, req)
	if err != nil {
	     return nil, err
	}
	return connect.NewResponse(resp), nil
}

