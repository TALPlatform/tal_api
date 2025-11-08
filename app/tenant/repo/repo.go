package repo

import (
	// INJECT IMPORTS
	"context"

	"github.com/TALPlatform/tal_api/db"
)

type TenantRepoInterface interface {
	PartialDeleteRestore(ctx context.Context, req *[]int32) (*[]db.TenantsSchemaPartial, error)
	PartialCreateUpdate(ctx context.Context, req *db.PartialCreateUpdateParams) (*db.TenantsSchemaPartial, error)
	PartialList(ctx context.Context, tenantId int32) (*[]db.TenantsSchemaPartial, error)
	// INJECT INTERFACE
	SectionListInpt(ctx context.Context) ([]db.SectionListInptRow, error)

	PartialTypeListInput(ctx context.Context) ([]db.PartialTypeListInputRow, error)

	PartialFindForUpdate(ctx context.Context, req *db.PartialFindForUpdateParams) (*db.TenantsSchemaPartial, error)

	PageFindForUpdate(ctx context.Context, req db.PageFindForUpdateParams) (*db.TenantsSchemaPage, error)
	TenantListInput(ctx context.Context, tenantId int32) (*[]db.TenantListInputRow, error)
	SectionDeleteRestore(ctx context.Context, req *[]int32) (*[]db.TenantsSchemaSection, error)
	SectionFind(ctx context.Context, req db.SectionFindParams) (*db.TenantsSchemaSection, error)
	SectionCreateUpdate(ctx context.Context, req *db.SectionCreateUpdateParams) (*db.TenantsSchemaSection, error)
	SectionList(ctx context.Context, tenantId int32) (*[]db.TenantsSchemaSection, error)
	PageDeleteRestore(ctx context.Context, req *[]int32) (*[]db.TenantsSchemaPage, error)
	PageCreateUpdate(ctx context.Context, req *db.PageCreateUpdateParams) (*db.TenantsSchemaPage, error)
	PageList(ctx context.Context, tenantId int32) (*[]db.TenantsSchemaPage, error)
	TenantDeleteRestore(ctx context.Context, req []int32) (*[]db.TenantsSchemaTenant, error)
	TenantList(ctx context.Context, tenantId int32) (*[]db.TenantsSchemaTenant, error)
	TenantFind(ctx context.Context, tenantId int32) (*db.TenantFindRow, error)
	TenantCreateUpdate(ctx context.Context, req *db.TenantCreateUpdateParams) (*db.TenantsSchemaTenant, error)
	TenantDashboard(ctx context.Context) (*[]db.TenantDashboardRow, error)
}

type TenantRepo struct {
	store        db.Store
	errorHandler map[string]string
}

func NewTenantRepo(store db.Store) TenantRepoInterface {
	errorHandler := map[string]string{
		// INJECT ERROR
	}
	return &TenantRepo{
		store:        store,
		errorHandler: errorHandler,
	}
}
