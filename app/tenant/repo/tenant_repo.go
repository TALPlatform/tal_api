package repo

import (
	"context"

	"github.com/TALPlatform/tal_api/db"
	"github.com/TALPlatform/tal_api/pkg/contextkeys"
	"github.com/rs/zerolog/log"
)

func (repo *TenantRepo) TenantCreateUpdate(ctx context.Context, req *db.TenantCreateUpdateParams) (*db.TenantsSchemaTenant, error) {
	resp, err := repo.store.TenantCreateUpdate(ctx, *req)
	if err != nil {
		return nil, repo.store.DbErrorParser(err, repo.errorHandler)
	}
	return &resp, nil
}

func (repo *TenantRepo) TenantFind(ctx context.Context, tenantId int32) (*db.TenantFindRow, error) {
	log.Debug().Interface("cache not found and db hitted for ", tenantId).Msg("cahchow not found")
	tenateID, _ := contextkeys.TenantID(ctx)
	if tenateID > 0 {
		tenantId = tenateID
	}
	resp, err := repo.store.TenantFind(ctx, tenantId)
	if err != nil {
		return nil, repo.store.DbErrorParser(err, repo.errorHandler)
	}
	return &resp, nil
}

func (repo *TenantRepo) TenantList(ctx context.Context, tenantId int32) (*[]db.TenantsSchemaTenant, error) {
	tenateID, _ := contextkeys.TenantID(ctx)
	if tenateID > 0 {
		tenantId = tenateID
	}
	resp, err := repo.store.TenantList(ctx, tenantId)
	if err != nil {
		return nil, repo.store.DbErrorParser(err, repo.errorHandler)
	}
	return &resp, nil
}

func (repo *TenantRepo) TenantDashboard(ctx context.Context) (*[]db.TenantDashboardRow, error) {
	tenateID, _ := contextkeys.TenantID(ctx)
	resp, err := repo.store.TenantDashboard(ctx, tenateID)
	if err != nil {
		return nil, repo.store.DbErrorParser(err, repo.errorHandler)
	}
	return &resp, nil
}
func (repo *TenantRepo) TenantListInput(ctx context.Context, tenantId int32) (*[]db.TenantListInputRow, error) {
	resp, err := repo.store.TenantListInput(ctx)
	if err != nil {
		return nil, repo.store.DbErrorParser(err, repo.errorHandler)
	}
	return &resp, nil
}
func (repo *TenantRepo) TenantDeleteRestore(ctx context.Context, req []int32) (*[]db.TenantsSchemaTenant, error) {
	resp, err := repo.store.TenantDeleteRestore(ctx, req)
	if err != nil {
		return nil, repo.store.DbErrorParser(err, repo.errorHandler)
	}
	return &resp, nil
}
