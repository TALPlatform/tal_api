package usecase

import (
	"context"

	"connectrpc.com/connect"
	"github.com/TALPlatform/tal_api/db"
	talv1 "github.com/TALPlatform/tal_api/proto_gen/tal/v1"
	"github.com/rs/zerolog/log"
)

func (u *TenantUsecase) PageList(ctx context.Context, req *connect.Request[talv1.PageListRequest]) (*talv1.PageListResponse, error) {
	record, err := u.repo.PageList(ctx, 0)
	if err != nil {
		return nil, err
	}

	resp := u.adapter.PageListGrpcFromSql(record)
	return resp, nil

}

func (u *TenantUsecase) PageCreateUpdate(ctx context.Context, req *connect.Request[talv1.PageCreateUpdateRequest]) (*talv1.PageCreateUpdateResponse, error) {

	sqlReq := u.adapter.PageCreateUpdateSqlFromGrpc(req.Msg)
	record, err := u.repo.PageCreateUpdate(ctx, sqlReq)
	if err != nil {
		return nil, err
	}
	resp := u.adapter.PageEntityGrpcFromSql(record)

	err = u.redisClient.TenantDelete(ctx, req.Msg.GetTenantId())
	if err != nil {
		log.Error().Str("message", "clear cache failed :").Err(err).Msg("Cache Clear Failed")
	}
	return &talv1.PageCreateUpdateResponse{Record: resp}, nil

}

func (u *TenantUsecase) PageDeleteRestore(ctx context.Context, req *connect.Request[talv1.PageDeleteRestoreRequest]) (*talv1.PageDeleteRestoreResponse, error) {
	record, err := u.repo.PageDeleteRestore(ctx, &req.Msg.Records)
	if err != nil {
		return nil, err
	}
	for _, r := range *record {
		if r.TenantID.Valid {
			err = u.redisClient.TenantDelete(ctx, r.TenantID.Int32)
			if err != nil {
				log.Error().Str("message", "clear cache failed :").Err(err).Msg("Cache Clear Failed")
			}
		}
	}
	resp := u.adapter.PageEntityListGrpcFromSql(record)
	return &talv1.PageDeleteRestoreResponse{Records: *resp}, nil

}

func (u *TenantUsecase) PageFindForUpdate(ctx context.Context, req *connect.Request[talv1.PageFindForUpdateRequest]) (*talv1.PageFindForUpdateResponse, error) {
	record, err := u.repo.PageFindForUpdate(ctx, db.PageFindForUpdateParams{PageID: req.Msg.RecordId})
	if err != nil {
		return nil, err
	}
	resp := u.adapter.PageFindForUpdateGrpcFromSql(record)
	return resp, nil
}
