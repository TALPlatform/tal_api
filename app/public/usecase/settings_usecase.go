package usecase

import (
	"context"
	talv1 "github.com/TALPlatform/tal_api/proto_gen/tal/v1"
)

func (s *PublicUsecase) SettingUpdate(ctx context.Context, req *talv1.SettingUpdateRequest) error {
	params := s.adapter.SettingUpdateSqlFromGrpc(req)
	err := s.repo.SettingUpdate(ctx, params)
	if err != nil {
		return err
	}
	return nil

}

func (u *PublicUsecase) SettingFindForUpdate(ctx context.Context, req *talv1.SettingFindForUpdateRequest) (*talv1.SettingFindForUpdateResponse, error) {
	settings, err := u.repo.SettingFindForUpdate(ctx)

	if err != nil {
		return nil, err
	}
	resp := u.adapter.SettingFindForUpdateGrpcFromSql(settings)

	return resp, nil
}
