package repo

import (
	"context"

	"github.com/TALPlatform/tal_api/db"
	"github.com/TALPlatform/tal_api/pkg/contextkeys"
)

func (repo *TenantRepo) PageList(ctx context.Context, tenantId int32) (*[]db.TenantsSchemaPage, error) {
	loggedInUserTenantId, _ := contextkeys.TenantID(ctx)
	if loggedInUserTenantId > 0 {
		tenantId = loggedInUserTenantId
	}
	resp, err := repo.store.PageList(ctx, tenantId)

	if err != nil {
		return nil, repo.store.DbErrorParser(err, repo.errorHandler)
	}
	return &resp, nil
}
func (repo *TenantRepo) PageCreateUpdate(ctx context.Context, req *db.PageCreateUpdateParams) (*db.TenantsSchemaPage, error) {
	loggedInUserTenantId, _ := contextkeys.TenantID(ctx)
	if loggedInUserTenantId > 0 {
		req.TenantID = loggedInUserTenantId
	}
	resp, err := repo.store.PageCreateUpdate(ctx, *req)

	if err != nil {
		return nil, repo.store.DbErrorParser(err, repo.errorHandler)
	}
	return &resp, nil
}
func (repo *TenantRepo) PageDeleteRestore(ctx context.Context, req *[]int32) (*[]db.TenantsSchemaPage, error) {

	resp, err := repo.store.PageDeleteRestore(ctx, *req)
	if err != nil {
		return nil, repo.store.DbErrorParser(err, repo.errorHandler)
	}
	return &resp, nil
}
func (repo *TenantRepo) PageFindForUpdate(ctx context.Context, req db.PageFindForUpdateParams) (*db.TenantsSchemaPage, error) {
	loggedInUserTenantId, _ := contextkeys.TenantID(ctx)
	if loggedInUserTenantId > 0 {
		req.TenantID = loggedInUserTenantId
	}
	resp, err := repo.store.PageFindForUpdate(ctx, req)

	if err != nil {
		return nil, repo.store.DbErrorParser(err, repo.errorHandler)
	}
	return &resp, nil
}
