package repo

import (
	"context"

	"github.com/TALPlatform/tal_api/db"
)

func (repo *PublicRepo) TranslationList(ctx context.Context) ([]db.Translation, error) {
	response, err := repo.store.TranslationList(ctx)
	if err != nil {
		return nil, repo.store.DbErrorParser(err, repo.errorHandler)
	}
	return response, nil
}

func (repo *PublicRepo) TranslationCreateUpdateBulk(ctx context.Context, req db.TranslationCreateUpdateBulkParams) ([]db.TranslationCreateUpdateBulkRow, error) {
	response, err := repo.store.TranslationCreateUpdateBulk(ctx, req)
	if err != nil {
		return nil, repo.store.DbErrorParser(err, repo.errorHandler)
	}
	return response, nil
}
func (repo *PublicRepo) TranslationDelete(ctx context.Context, req []string) ([]db.Translation, error) {
	response, err := repo.store.TranslationDelete(ctx, req)
	if err != nil {
		return nil, repo.store.DbErrorParser(err, repo.errorHandler)
	}
	return response, nil
}
