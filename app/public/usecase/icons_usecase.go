package usecase

import (
	"context"

	talv1 "github.com/TALPlatform/tal_api/proto_gen/tal/v1"
)

func (u *PublicUsecase) IconFind(ctx context.Context, req *talv1.IconFindRequest) (*talv1.IconFindResponse, error) {
	params := u.adapter.IconFindSqlFromGrpc(req)
	icons, err := u.repo.IconFind(ctx, *params)
	if err != nil {
		return nil, err
	}
	res := u.adapter.IconGrpcFromSql(icons)
	return &talv1.IconFindResponse{Icon: res}, nil
}
func (u *PublicUsecase) IconCreateUpdateBulk(ctx context.Context, req *talv1.IconCreateUpdateBulkRequest) (*talv1.IconListResponse, error) {
	params := u.adapter.IconCreateUpdateBulkSqlFromGrpc(req)
	icons, err := u.repo.IconCreateUpdateBulk(ctx, params)
	if err != nil {
		return nil, err
	}
	res := u.adapter.IconListGrpcFromSql(*icons)
	return res, nil
}
func (u *PublicUsecase) IconList(ctx context.Context) (*talv1.IconListResponse, error) {
	icons, err := u.repo.IconList(ctx)
	if err != nil {
		return nil, err
	}
	res := u.adapter.IconListGrpcFromSql(*icons)
	return res, nil
}
