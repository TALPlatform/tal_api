package api

import (
	// INJECT IMPORTS
	"context"

	"connectrpc.com/connect"
	talv1 "github.com/TALPlatform/tal_api/proto_gen/tal/v1"
	"github.com/rs/zerolog/log"
)

func (api *Api) PartialList(ctx context.Context, req *connect.Request[talv1.PartialListRequest]) (*connect.Response[talv1.PartialListResponse], error) {
	resp, err := api.tenantUsecase.PartialList(ctx, req)
	if err != nil {
		return nil, err
	}

	resp.Options = api.getAvailableOptions(req.Header())
	return connect.NewResponse(resp), nil
}

func (api *Api) PartialCreateUpdate(ctx context.Context, req *connect.Request[talv1.PartialCreateUpdateRequest]) (*connect.Response[talv1.PartialCreateUpdateResponse], error) {
	var err error
	if req.Msg.Uploads != nil {
		_, err := api.publicUsecase.FileCreateBulk(ctx, req.Msg.Uploads)
		if err != nil {
			log.Error().Err(err).Msg("error uploading files")
		}
	}
	resp, err := api.tenantUsecase.PartialCreateUpdate(ctx, req)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}
func (api *Api) PartialDeleteRestore(ctx context.Context, req *connect.Request[talv1.PartialDeleteRestoreRequest]) (*connect.Response[talv1.PartialDeleteRestoreResponse], error) {
	resp, err := api.tenantUsecase.PartialDeleteRestore(ctx, req)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}
func (api *Api) PartialFindForUpdate(ctx context.Context, req *connect.Request[talv1.PartialFindForUpdateRequest]) (*connect.Response[talv1.PartialFindForUpdateResponse], error) {
	resp, err := api.tenantUsecase.PartialFindForUpdate(ctx, req)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}
