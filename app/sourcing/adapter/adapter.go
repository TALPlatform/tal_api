package adapter

import (
	"github.com/TALPlatform/tal_api/db"
	talv1 "github.com/TALPlatform/tal_api/proto_gen/tal/v1"
)

type SourcingAdapterInterface interface {
	// INJECT INTERFACE

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
