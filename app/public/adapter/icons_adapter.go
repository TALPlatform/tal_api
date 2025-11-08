package adapter

import (
	"github.com/TALPlatform/tal_api/db"
	talv1 "github.com/TALPlatform/tal_api/proto_gen/tal/v1"
)

func (a *PublicAdapter) IconFindSqlFromGrpc(icon *talv1.IconFindRequest) *db.IconFindParams {
	return &db.IconFindParams{
		IconID:   icon.IconId,
		IconName: icon.IconName,
	}

}
func (a *PublicAdapter) IconGrpcFromSql(icon *db.Icon) *talv1.Icon {
	return &talv1.Icon{
		IconId:      icon.IconID,
		IconName:    icon.IconName,
		IconContent: icon.IconContent,
	}

}

func (a *PublicAdapter) IconCreateUpdateBulkSqlFromGrpc(req *talv1.IconCreateUpdateBulkRequest) db.IconCreateUpdateBulkParams {
	names := make([]string, len(req.Icons))
	contents := make([]string, len(req.Icons))
	for index, v := range req.Icons {
		names[index] = v.IconName
		contents[index] = v.IconContent
	}
	return db.IconCreateUpdateBulkParams{
		IconsNames:    names,
		IconsContents: contents,
	}
}

func (a *PublicAdapter) IconListGrpcFromSql(resp []db.Icon) *talv1.IconListResponse {
	records := make([]*talv1.Icon, len(resp))
	for index, v := range resp {
		records[index] = a.IconGrpcFromSql(&v)
	}
	return &talv1.IconListResponse{
		Icons: records,
	}
}
