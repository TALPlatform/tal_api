package adapter

import (
	"github.com/TALPlatform/tal_api/db"
	talv1 "github.com/TALPlatform/tal_api/proto_gen/tal/v1"
)

type TenantAdapterInterface interface {
	// INJECT INTERFACE

	SectionListInptGrpcFromSql(req *[]db.SectionListInptRow) *talv1.SectionListInptResponse

	PartialCreateUpdateSqlFromGrpc(req *talv1.PartialCreateUpdateRequest) *db.PartialCreateUpdateParams
	PartialListGrpcFromSql(resp *[]db.TenantsSchemaPartial) *talv1.PartialListResponse
	PartialEntityListGrpcFromSql(resp *[]db.TenantsSchemaPartial) *[]*talv1.TenantsSchemaPartial
	PartialFindForUpdateGrpcFromSql(resp *db.TenantsSchemaPartial) *talv1.PartialFindForUpdateResponse
	PartialEntityGrpcFromSql(resp *db.TenantsSchemaPartial) *talv1.TenantsSchemaPartial
	SectionEntityGrpcFromSql(resp *db.TenantsSchemaSection) *talv1.TenantsSchemaSection
	SectionEntityListGrpcFromSql(resp *[]db.TenantsSchemaSection) *[]*talv1.TenantsSchemaSection
	SectionListGrpcFromSql(resp *[]db.TenantsSchemaSection) *talv1.SectionListResponse
	SectionCreateUpdateSqlFromGrpc(req *talv1.SectionCreateUpdateRequest) *db.SectionCreateUpdateParams
	PageCreateUpdateSqlFromGrpc(req *talv1.PageCreateUpdateRequest) *db.PageCreateUpdateParams
	SectionFindForUpdateSqlFromGrpc(req *talv1.SectionFindForUpdateRequest) *db.SectionFindParams
	PageEntityGrpcFromSql(resp *db.TenantsSchemaPage) *talv1.TenantsSchemaPage
	SectionFindForUpdateGrpcFromSql(resp *db.TenantsSchemaSection) *talv1.SectionFindForUpdateResponse
	PageFindForUpdateGrpcFromSql(resp *db.TenantsSchemaPage) *talv1.PageFindForUpdateResponse
	PageListGrpcFromSql(resp *[]db.TenantsSchemaPage) *talv1.PageListResponse
	PageEntityListGrpcFromSql(req *[]db.TenantsSchemaPage) *[]*talv1.TenantsSchemaPage
	TenantDeleteRestoreGrpcFromSql(resp *[]db.TenantsSchemaTenant) *talv1.TenantDeleteRestoreResponse
	TenantListGrpcFromSql(resp *[]db.TenantsSchemaTenant) *talv1.TenantListResponse
	TenantFindGrpcFromSql(resp *db.TenantFindRow) *talv1.TenantFindResponse
	TenantListInputGrpcFromSql(resp *[]db.TenantListInputRow) *talv1.TenantListInputResponse
	PartialTypeListInputGrpcFromSql(resp []db.PartialTypeListInputRow) *talv1.PartialTypeListInputResponse
	TenantEntityGrpcFromSql(resp *db.TenantsSchemaTenant) *talv1.TenantsSchemaTenant
	TenantCreateUpdateSqlFromGrpc(req *talv1.TenantCreateUpdateRequest) *db.TenantCreateUpdateParams
	TenantDashboardGrpcFromSql(resp *[]db.TenantDashboardRow) *talv1.TenantDashboardResponse
}

type TenantAdapter struct {
}

func NewTenantAdapter() TenantAdapterInterface {
	return &TenantAdapter{}
}
