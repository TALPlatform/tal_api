package api

import (
	"connectrpc.com/connect"
	"context"
	"github.com/TALPlatform/tal_api/proto_gen/tal/v1"
)

func (api *Api) EmailSend(ctx context.Context, req *connect.Request[talv1.EmailSendRequest]) (*connect.Response[talv1.EmailSendResponse], error) {
	if err := ctx.Err(); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	response, err := api.publicUsecase.EmailSend(ctx, req.Msg)
	return connect.NewResponse(response), err
}
