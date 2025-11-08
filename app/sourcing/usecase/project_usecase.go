package usecase

import (
	"context"

	"connectrpc.com/connect"
	talv1 "github.com/TALPlatform/tal_api/proto_gen/tal/v1"
)

func (u *SourcingUsecase) ProjectInputList(ctx context.Context, req *connect.Request[talv1.ProjectInputListRequest]) (*talv1.ProjectInputListResponse, error) {
	record, err := u.repo.ProjectInputList(ctx)
	if err != nil {
		return nil, err
	}
	resp := u.adapter.ProjectInputListGrpcFromSql(record)
	return resp, nil

}
