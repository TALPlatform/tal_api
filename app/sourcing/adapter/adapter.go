package adapter

import (
	"github.com/TALPlatform/tal_api/db"
	"github.com/TALPlatform/tal_api/pkg/agenthub"
	"github.com/TALPlatform/tal_api/pkg/crustdata"
	talv1 "github.com/TALPlatform/tal_api/proto_gen/tal/v1"
)

type SourcingAdapterInterface interface {
	// INJECT INTERFACE
	SourcingSessionFiltersInferedGrpcFromAgent(req *agenthub.SourcingSessionFiltersInfered) *talv1.SourcingSessionFiltersInferedResponse

	SourcingSessionApplyGrpcFromCrustdata(row *crustdata.PersonDBProfile) *talv1.SourcingSessionApplyResponse
	SourcingSessionApplyCrustDataFromSql(req *db.SourcingSessionFindRow) *crustdata.PeopleSearchRequest
	SourcingSessionApplyGrpcFromSql(req *db.SourcingSessionApplyRow) *talv1.SourcingSessionApplyResponse
	SourcingSessionCriteriaCreateSqlFromGrpc(req *talv1.SourcingSessionCriteriaCreateRequest) *db.SourcingSessionCriteriaCreateParams
	SourcingSessionCriteriaCreateGrpcFromSql(req *db.SourcingSessionCriteriaProfilesBulkInsertParams) *talv1.SourcingSessionCriteriaCreateResponse
	// SourcingSessionProfileFindForAISqlFromGrpc(req *talv1.SourcingSessionProfileJustifyRequest) *db.SourcingSessionProfileFindForAIParams
	SourcingSessionListGrpcFromSql(req *[]db.SourcingSessionListRow) *talv1.SourcingSessionListResponse

	SourcingSessionProfileCreateUpdateSqlFromGrpc(req *talv1.SourcingSessionProfileCreateUpdateRequest) *db.SourcingSessionProfileCreateUpdateParams
	SourcingSessionProfileCreateUpdateGrpcFromSql(req *db.SourcingSessionProfileCreateUpdateRow) *talv1.SourcingSessionProfileCreateUpdateResponse
	SourcingSessionFindGrpcFromSql(req *db.SourcingSessionFindRow) (*talv1.SourcingSessionFindResponse, error)
	SourcingSessionCreateUpdateSqlFromGrpc(req *talv1.SourcingSessionCreateUpdateRequest) *db.SourcingSessionCreateUpdateParams
	SourcingSessionCreateUpdateGrpcFromSql(req *db.SourcingSessionCreateUpdateRow) *talv1.SourcingSessionCreateUpdateResponse
	ProjectInputListGrpcFromSql(resp *[]db.ProjectInputListRow) *talv1.ProjectInputListResponse
}

type SourcingAdapter struct {
}

func NewSourcingAdapter() SourcingAdapterInterface {
	return &SourcingAdapter{}
}
