package adapter

import (
	"context"

	"github.com/TALPlatform/tal_api/db"
	"github.com/TALPlatform/tal_api/pkg/crustdata"
	talv1 "github.com/TALPlatform/tal_api/proto_gen/tal/v1"
	"github.com/darwishdev/genaiclient"
)

type PeopleAdapterInterface interface {
	// INJECT INTERFACE

	RawProfileFindGrpcFromSql(req *db.RawProfileFindRow) *talv1.RawProfileFindResponse
	RawProfileListSqlFromGrpc(req *talv1.RawProfileListRequest) *db.RawProfileListParams
	RawProfileListCrustDataFromGrpc(req *talv1.RawProfileListRequest) *crustdata.PeopleSearchRequest
	RawProfileListDbFromCrustData(req *[]*crustdata.PersonDBProfile) *[]*db.RawProfileListRow
	RawProfileListGrpcFromSql(req *[]*db.RawProfileListRow, crustDataReq ...*[]*db.RawProfileListRow) *talv1.RawProfileListResponse
	RawProfileListEnrichAndMarshal(
		ctx context.Context,
		session_id int32,
		profiles *crustdata.PeopleSearchResponse,
	) (*db.RawProfilesBulkCreateUpdateParams, error)
}

type PeopleAdapter struct {
	embed     func(ctx context.Context, text string, options ...*genaiclient.EmbedOptions) ([]float32, error)
	embedBulk func(ctx context.Context, text []string, options ...*genaiclient.EmbedOptions) ([][]float32, error)
}

func NewPeopleAdapter(embed func(ctx context.Context, text string, options ...*genaiclient.EmbedOptions) ([]float32, error), embedBulk func(ctx context.Context, text []string, options ...*genaiclient.EmbedOptions) ([][]float32, error)) PeopleAdapterInterface {
	return &PeopleAdapter{
		embed:     embed,
		embedBulk: embedBulk,
	}
}
