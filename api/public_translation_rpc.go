package api

import (
	"context"

	"connectrpc.com/connect"
	talv1 "github.com/TALPlatform/tal_api/proto_gen/tal/v1"
)

func (api *Api) TranslationDelete(ctx context.Context, req *connect.Request[talv1.TranslationDeleteRequest]) (*connect.Response[talv1.TranslationDeleteResponse], error) {
	if err := ctx.Err(); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	response, err := api.publicUsecase.TranslationDelete(ctx, req.Msg)
	return connect.NewResponse(response), err
}

func (api *Api) TranslationFindLocale(ctx context.Context, req *connect.Request[talv1.TranslationFindLocaleRequest]) (*connect.Response[talv1.TranslationFindLocaleResponse], error) {
	if err := ctx.Err(); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	response, err := api.publicUsecase.TranslationFindLocale(ctx, req.Msg)
	return connect.NewResponse(response), err
}

func (api *Api) TranslationList(ctx context.Context, req *connect.Request[talv1.TranslationListRequest]) (*connect.Response[talv1.TranslationListResponse], error) {
	if err := ctx.Err(); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	response, err := api.publicUsecase.TranslationList(ctx)
	return connect.NewResponse(response), err
}

func (api *Api) TranslationCreateUpdateBulk(ctx context.Context, req *connect.Request[talv1.TranslationCreateUpdateBulkRequest]) (*connect.Response[talv1.TranslationCreateUpdateBulkResponse], error) {
	if err := ctx.Err(); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	resp, err := api.publicUsecase.TranslationCreateUpdateBulk(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}
