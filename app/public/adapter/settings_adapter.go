package adapter

import (
	"github.com/TALPlatform/tal_api/db"
	talv1 "github.com/TALPlatform/tal_api/proto_gen/tal/v1"
)

func (a *PublicAdapter) SettingUpdateSqlFromGrpc(req *talv1.SettingUpdateRequest) *db.SettingUpdateParams {
	keys := make([]string, len(req.Settings))
	values := make([]string, len(req.Settings))
	for index, v := range req.Settings {
		keys[index] = v.SettingKey
		values[index] = v.SettingValue
	}
	return &db.SettingUpdateParams{
		Keys:   keys,
		Values: values,
	}
}
func (a *PublicAdapter) SettingEntityGrpcFromSql(resp []db.Setting) []*talv1.Setting {
	grpcResp := make([]*talv1.Setting, len(resp))
	for _, v := range resp {
		record := &talv1.Setting{
			SettingKey:   v.SettingKey,
			SettingValue: v.SettingValue,
		}
		grpcResp = append(grpcResp, record)
	}
	return grpcResp

}

func (a *PublicAdapter) SettingFindForUpdateGrpcFromSql(resp *[]db.SettingFindForUpdateRow) *talv1.SettingFindForUpdateResponse {
	grpcRows := make([]*talv1.SettingFindForUpdateRow, len(*resp))
	for index, v := range *resp {
		grpcRow := &talv1.SettingFindForUpdateRow{
			SettingKey:   v.SettingKey,
			SettingValue: v.SettingValue,
			InputType:    v.InputTypeName,
		}

		grpcRows[index] = grpcRow

	}

	return &talv1.SettingFindForUpdateResponse{
		Settings: grpcRows,
	}

}
