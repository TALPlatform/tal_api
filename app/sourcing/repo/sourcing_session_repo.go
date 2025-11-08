package repo

import (
	"context"

	"github.com/TALPlatform/tal_api/db"
	"github.com/TALPlatform/tal_api/pkg/contextkeys"
)

func (repo *SourcingRepo) SourcingSessionCreateUpdate(ctx context.Context, req *db.SourcingSessionCreateUpdateParams) (*db.SourcingSessionCreateUpdateRow, error) {
	callerId, _ := contextkeys.CallerID(ctx)
	req.CreatedBy = callerId
	resp, err := repo.store.SourcingSessionCreateUpdate(ctx, *req)

	if err != nil {
		return nil, repo.store.DbErrorParser(err, repo.errorHandler)
	}
	return &resp, nil
}
func (repo *SourcingRepo) SourcingSessionFind(ctx context.Context, req int32) (*db.SourcingSessionFindRow, error) {
	resp, err := repo.store.SourcingSessionFind(ctx, req)
	if err != nil {
		return nil, repo.store.DbErrorParser(err, repo.errorHandler)
	}
	return &resp, nil
}
func (repo *SourcingRepo) SourcingSessionProfileCreateUpdate(ctx context.Context, req *db.SourcingSessionProfileCreateUpdateParams) (*db.SourcingSessionProfileCreateUpdateRow, error) {
	resp, err := repo.store.SourcingSessionProfileCreateUpdate(ctx, *req)
	if err != nil {
		return nil, repo.store.DbErrorParser(err, repo.errorHandler)
	}
	return &resp, nil
}
