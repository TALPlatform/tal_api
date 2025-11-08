package repo

import (
	"context"

	"github.com/TALPlatform/tal_api/db"
	"github.com/TALPlatform/tal_api/pkg/contextkeys"
)

func (repo *SourcingRepo) ProjectInputList(ctx context.Context) (*[]db.ProjectInputListRow, error) {
	tenantID, _ := contextkeys.TenantID(ctx)
	resp, err := repo.store.ProjectInputList(ctx, tenantID)
	if err != nil {
		return nil, repo.store.DbErrorParser(err, repo.errorHandler)
	}
	return &resp, nil
}
