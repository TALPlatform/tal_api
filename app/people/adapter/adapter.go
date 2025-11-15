package adapter

import (
	"context"

	"github.com/TALPlatform/tal_api/db"
	"github.com/TALPlatform/tal_api/pkg/crustdata"
	talv1 "github.com/TALPlatform/tal_api/proto_gen/tal/v1"
)

type PeopleAdapterInterface interface {
	// INJECT INTERFACE
	RawProfileListGrpcFromSql(r *db.RawProfileListRow) *talv1.RawProfileListResponse
	RawProfileListGrpcFromCrustdata(row *crustdata.PersonDBProfile) *talv1.RawProfileListResponse
	RawProfileFindGrpcFromSql(req *db.RawProfileFindRow) *talv1.RawProfileFindResponse
	RawProfileListSqlFromGrpc(req *talv1.RawProfileListRequest) *db.RawProfileListParams
	RawProfileListCrustDataFromGrpc(req *talv1.RawProfileListRequest) *crustdata.PeopleSearchRequest
	RawProfileListEnrichAndMarshal(
		ctx context.Context,
		session_id int32,
		profiles *crustdata.PeopleSearchResponse,
	) (*db.RawProfilesBulkCreateUpdateParams, error)
}

type PeopleAdapter struct {
	// embed     func(ctx context.Context, text string, options ...*genaiclient.EmbedOptions) ([]float32, error)
	// embedBulk func(ctx context.Context, text []string, options ...*genaiclient.EmbedOptions) ([][]float32, error)
}

func NewPeopleAdapter() PeopleAdapterInterface {
	return &PeopleAdapter{
		// embed:     embed,
		// embedBulk: embedBulk,
	}
}
