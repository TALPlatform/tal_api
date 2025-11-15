package api

import (
	// INJECT IMPORTS
	"context"

	"connectrpc.com/connect"
	talv1 "github.com/TALPlatform/tal_api/proto_gen/tal/v1"
)

func (api *Api) RawProfileList(ctx context.Context, req *connect.Request[talv1.RawProfileListRequest], stream *connect.ServerStream[talv1.RawProfileListResponse]) error {
	err := api.peopleUsecase.RawProfileList(ctx, req, stream)
	if err != nil {
		return err
	}
	return nil
}
func (api *Api) RawProfileFind(ctx context.Context, req *connect.Request[talv1.RawProfileFindRequest]) (*connect.Response[talv1.RawProfileFindResponse], error) {
	resp, err := api.peopleUsecase.RawProfileFind(ctx, req)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

func (api *Api) RawProfileListRequestBuild(ctx context.Context, req *connect.Request[talv1.RawProfileListRequestBuildRequest]) (*connect.Response[talv1.RawProfileListRequestBuildResponse], error) {
	resp, err := api.peopleUsecase.RawProfileListRequestBuild(ctx, req)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}
