package adapter

import (
	"github.com/TALPlatform/tal_api/db"
	talv1 "github.com/TALPlatform/tal_api/proto_gen/tal/v1"
)

func (a *SourcingAdapter) ProjectInputListGrpcFromSql(resp *[]db.ProjectInputListRow) *talv1.ProjectInputListResponse {
	records := make([]*talv1.SelectInputOption, 0)
	for _, v := range *resp {
		records = append(records, &talv1.SelectInputOption{
			Value: v.Value,
			Label: v.Label,
		})
	}
	return &talv1.ProjectInputListResponse{
		Options: records,
	}
}
