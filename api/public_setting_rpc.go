package api

import (
	"context"

	"connectrpc.com/connect"
	talv1 "github.com/TALPlatform/tal_api/proto_gen/tal/v1"
)

func (api *Api) SettingUpdate(ctx context.Context, req *connect.Request[talv1.SettingUpdateRequest]) (*connect.Response[talv1.SettingUpdateResponse], error) {
	if err := ctx.Err(); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	err := api.publicUsecase.SettingUpdate(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&talv1.SettingUpdateResponse{}), nil
}

func (api *Api) SettingFindForUpdate(ctx context.Context, req *connect.Request[talv1.SettingFindForUpdateRequest]) (*connect.Response[talv1.SettingFindForUpdateResponse], error) {
	if err := ctx.Err(); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	resp, err := api.publicUsecase.SettingFindForUpdate(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}
