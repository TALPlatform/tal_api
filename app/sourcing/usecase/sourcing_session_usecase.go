package usecase

import (
	"context"

	"connectrpc.com/connect"
	talv1 "github.com/TALPlatform/tal_api/proto_gen/tal/v1"
)

func (u *SourcingUsecase) SourcingSessionCreateUpdate(ctx context.Context, req *connect.Request[talv1.SourcingSessionCreateUpdateRequest]) (*talv1.SourcingSessionCreateUpdateResponse, error) {
	sqlReq := u.adapter.SourcingSessionCreateUpdateSqlFromGrpc(req.Msg)
	record, err := u.repo.SourcingSessionCreateUpdate(ctx, sqlReq)
	if err != nil {
		return nil, err
	}
	resp := u.adapter.SourcingSessionCreateUpdateGrpcFromSql(record)
	return resp, nil
}

func (u *SourcingUsecase) SourcingSessionFind(ctx context.Context, req *connect.Request[talv1.SourcingSessionFindRequest]) (*talv1.SourcingSessionFindResponse, error) {
	record, err := u.repo.SourcingSessionFind(ctx, req.Msg.SourcingSessionId)
	if err != nil {
		return nil, err
	}
	resp, err := u.adapter.SourcingSessionFindGrpcFromSql(record)
	if err != nil {
		return nil, err
	}
	return resp, nil

}

func (u *SourcingUsecase) SourcingSessionProfileCreateUpdate(ctx context.Context, req *connect.Request[talv1.SourcingSessionProfileCreateUpdateRequest]) (*talv1.SourcingSessionProfileCreateUpdateResponse, error) {
	sqlReq := u.adapter.SourcingSessionProfileCreateUpdateSqlFromGrpc(req.Msg)
	record, err := u.repo.SourcingSessionProfileCreateUpdate(ctx, sqlReq)
	if err != nil {
		return nil, err
	}
	resp := u.adapter.SourcingSessionProfileCreateUpdateGrpcFromSql(record)
	return resp, nil
}
