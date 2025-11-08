package repo

import (
	"context"

	"github.com/TALPlatform/tal_api/db"
)

func (repo *PublicRepo) SettingUpdate(ctx context.Context, req *db.SettingUpdateParams) error {
	_, err := repo.store.SettingUpdate(context.Background(), *req)

	if err != nil {
		return repo.store.DbErrorParser(err, repo.errorHandler)
	}
	return nil
}

func (repo *PublicRepo) SettingFindForUpdate(ctx context.Context) (*[]db.SettingFindForUpdateRow, error) {
	resp, err := repo.store.SettingFindForUpdate(context.Background())
	if err != nil {
		return nil, repo.store.DbErrorParser(err, repo.errorHandler)
	}
	return &resp, nil
}
