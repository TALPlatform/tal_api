package repo

import (
	"context"

	"github.com/TALPlatform/tal_api/db"
	"github.com/TALPlatform/tal_api/pkg/contextkeys"
)

func (repo *TenantRepo) PartialList(ctx context.Context, tenantId int32) (*[]db.TenantsSchemaPartial, error) {
	loggedInUserTenantId, _ := contextkeys.TenantID(ctx)
	if loggedInUserTenantId > 0 {
		tenantId = loggedInUserTenantId
	}

	resp, err := repo.store.PartialList(ctx, tenantId)

	if err != nil {
		return nil, repo.store.DbErrorParser(err, repo.errorHandler)
	}
	return &resp, nil
}
func (repo *TenantRepo) PartialCreateUpdate(ctx context.Context, req *db.PartialCreateUpdateParams) (*db.TenantsSchemaPartial, error) {
	resp, err := repo.store.PartialCreateUpdate(ctx, *req)

	if err != nil {
		return nil, repo.store.DbErrorParser(err, repo.errorHandler)
	}
	return &resp, nil
}
func (repo *TenantRepo) PartialDeleteRestore(ctx context.Context, req *[]int32) (*[]db.TenantsSchemaPartial, error) {
	resp, err := repo.store.PartialDeleteRestore(ctx, *req)
	if err != nil {
		return nil, repo.store.DbErrorParser(err, repo.errorHandler)
	}
	return &resp, nil
}
func (repo *TenantRepo) PartialFindForUpdate(ctx context.Context, req *db.PartialFindForUpdateParams) (*db.TenantsSchemaPartial, error) {
	resp, err := repo.store.PartialFindForUpdate(ctx, *req)

	if err != nil {
		return nil, repo.store.DbErrorParser(err, repo.errorHandler)
	}
	return &resp, nil
}
func (repo *TenantRepo) PartialTypeListInput(ctx context.Context) ([]db.PartialTypeListInputRow, error) {
	resp, err := repo.store.PartialTypeListInput(ctx)

	if err != nil {
		return nil, repo.store.DbErrorParser(err, repo.errorHandler)
	}
	return resp, nil
}
