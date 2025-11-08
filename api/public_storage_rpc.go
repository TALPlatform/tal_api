package api

import (
	"bytes"
	"context"

	"connectrpc.com/connect"
	talv1 "github.com/TALPlatform/tal_api/proto_gen/tal/v1"
)

func (api *Api) FileDelete(ctx context.Context, req *connect.Request[talv1.FileDeleteRequest]) (*connect.Response[talv1.FileDeleteResponse], error) {
	if err := ctx.Err(); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	response, err := api.publicUsecase.FileDelete(ctx, req.Msg)
	return connect.NewResponse(response), err
}

func (api *Api) FileDeleteByBucket(ctx context.Context, req *connect.Request[talv1.FileDeleteByBucketRequest]) (*connect.Response[talv1.FileDeleteByBucketResponse], error) {
	if err := ctx.Err(); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	response, err := api.publicUsecase.FileDeleteByBucket(ctx, req.Msg)
	return connect.NewResponse(response), err
}
func (api *Api) GalleryList(ctx context.Context, req *connect.Request[talv1.GalleryListRequest]) (*connect.Response[talv1.GalleryListResponse], error) {
	if err := ctx.Err(); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	response, err := api.publicUsecase.GalleryList(ctx, req.Msg)
	response.Options = api.getAvailableOptions(req.Header())
	response.Options.DeleteRestoreHandler = nil
	response.Options.UpdateHandler = nil
	return connect.NewResponse(response), err
}
func (api *Api) FileList(ctx context.Context, req *connect.Request[talv1.FileListRequest]) (*connect.Response[talv1.FileListResponse], error) {
	if err := ctx.Err(); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	response, err := api.publicUsecase.FileList(ctx, req.Msg)
	return connect.NewResponse(response), err
}
func (api *Api) BucketCreateUpdate(ctx context.Context, req *connect.Request[talv1.BucketCreateUpdateRequest]) (*connect.Response[talv1.BucketCreateUpdateResponse], error) {
	if err := ctx.Err(); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	response, err := api.publicUsecase.BucketCreateUpdate(ctx, req.Msg)
	return connect.NewResponse(response), err
}

func (api *Api) BucketList(ctx context.Context, req *connect.Request[talv1.BucketListRequest]) (*connect.Response[talv1.BucketListResponse], error) {
	if err := ctx.Err(); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	response, err := api.publicUsecase.BucketList(ctx, req.Msg)
	return connect.NewResponse(response), err
}

func (api *Api) FileUploadUrlFind(ctx context.Context, req *connect.Request[talv1.FileUploadUrlFindRequest]) (*connect.Response[talv1.FileUploadUrlFindResponse], error) {
	if err := ctx.Err(); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	response, err := api.publicUsecase.FileUploadUrlFind(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	maxAge := int(api.config.RefreshTokenDuration.Seconds())
	resp := connect.NewResponse(response)
	api.WithCookie(resp, api.config.SupabaseTokenCookieName, response.Token, maxAge)
	api.WithCookie(resp, api.config.SupabaseRefreshTokenCookieName, response.RefreshToken, maxAge)
	return resp, nil
}
func (api *Api) FileCreate(ctx context.Context, req *connect.Request[talv1.FileCreateRequest]) (*connect.Response[talv1.FileCreateResponse], error) {
	if err := ctx.Err(); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	resp, err := api.publicUsecase.FileCreate(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

func (api *Api) FileCreateBulk(ctx context.Context, req *connect.Request[talv1.FileCreateBulkRequest]) (*connect.Response[talv1.FileCreateBulkResponse], error) {
	if err := ctx.Err(); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	resp, err := api.publicUsecase.FileCreateBulk(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

func (api *Api) ImportTable(ctx context.Context, req *connect.Request[talv1.ImportTableRequest]) (*connect.Response[talv1.ImportTableResponse], error) {
	if err := ctx.Err(); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	buffer := bytes.NewBuffer(req.Msg.Reader)
	_, err := api.sqlSeeder.SeedFromExcel(*buffer, req.Msg.SchemaName, req.Msg.TableName, req.Msg.SchemaName)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&talv1.ImportTableResponse{
		Message: "imported",
	}), nil
}
