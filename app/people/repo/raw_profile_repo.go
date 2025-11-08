package repo

import (
	"context"

	"github.com/TALPlatform/tal_api/db"
)

func (repo *PeopleRepo) RawProfileList(ctx context.Context, req *db.RawProfileListParams) (*[]*db.RawProfileListRow, error) {
	resp, err := repo.store.RawProfileList(ctx, *req)
	if err != nil {
		return nil, repo.store.DbErrorParser(err, repo.errorHandler)
	}
	var response = make([]*db.RawProfileListRow, len(resp))
	for i, v := range resp {
		response[i] = &v
	}
	return &response, nil
}

func (repo *PeopleRepo) RawProfileBulkCreateUpdate(ctx context.Context, params *db.RawProfilesBulkCreateUpdateParams) error {
	err := repo.store.RawProfilesBulkCreateUpdate(ctx, *params)
	if err != nil {
		return repo.store.DbErrorParser(err, repo.errorHandler)
	}
	return nil
}
func (repo *PeopleRepo) RawProfileFind(ctx context.Context, req int32) (*db.RawProfileFindRow, error) {
	resp, err := repo.store.RawProfileFind(ctx, int64(req))
	if err != nil {
		return nil, repo.store.DbErrorParser(err, repo.errorHandler)
	}
	return &resp, nil
}
