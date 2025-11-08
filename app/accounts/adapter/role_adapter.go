package adapter

import (
	"github.com/TALPlatform/tal_api/db"
	"github.com/TALPlatform/tal_api/pkg/dateutils"
	talv1 "github.com/TALPlatform/tal_api/proto_gen/tal/v1"
)

func (a *AccountsAdapter) RoleFindForUpdateUpdateGrpcFromSql(resp *db.RoleFindForUpdateRow) *talv1.RoleCreateUpdateRequest {
	return &talv1.RoleCreateUpdateRequest{
		RoleId:            resp.RoleID,
		RoleSecurityLevel: resp.RoleSecurityLevel,
		TenantId:          resp.TenantID.Int32,
		RoleName:          resp.RoleName,
		RoleDescription:   resp.RoleDescription.String,
		Permissions:       resp.Permissions,
	}
}
func (a *AccountsAdapter) RoleListInputGrpcFromSql(resp *[]db.RoleListInputRow) *talv1.RoleListInputResponse {
	records := make([]*talv1.SelectInputOption, 0)
	for _, v := range *resp {
		records = append(records, &talv1.SelectInputOption{
			Value: v.Value,
			Note:  v.Note,
			Label: v.Label,
		})
	}
	return &talv1.RoleListInputResponse{
		Options: records,
	}
}
func (a *AccountsAdapter) RoleEntityGrpcFromSql(resp *db.AccountsSchemaRole) *talv1.AccountsSchemaRole {
	return &talv1.AccountsSchemaRole{
		RoleId:          int32(resp.RoleID),
		RoleName:        resp.RoleName,
		TenantId:        resp.TenantID.Int32,
		RoleDescription: resp.RoleDescription.String,
		CreatedAt:       db.TimeToString(resp.CreatedAt.Time),
		DeletedAt:       db.TimeToString(resp.DeletedAt.Time),
	}
}

func (a *AccountsAdapter) RoleCreateUpdateSqlFromGrpc(req *talv1.RoleCreateUpdateRequest) *db.RoleCreateUpdateParams {
	resp := &db.RoleCreateUpdateParams{
		RoleID:            req.RoleId,
		TenantID:          req.TenantId,
		RoleName:          req.RoleName,
		RoleSecurityLevel: req.RoleSecurityLevel,
		RoleDescription:   req.RoleDescription,
		Permissions:       req.Permissions,
	}

	return resp
}

func (a *AccountsAdapter) RoleListRowGrpcFromSql(resp *db.RoleListRow) *talv1.RoleListRow {
	return &talv1.RoleListRow{
		RoleId:            int32(resp.RoleID),
		RoleName:          resp.RoleName,
		TenantName:        resp.TenantName,
		RoleSecurityLevel: resp.RoleSecurityLevel,
		UserCount:         int32(resp.UserCount),
		PermissionCount:   int32(resp.PermissionCount),
		TenantId:          resp.TenantID,
		CreatedAt:         dateutils.DateTimeToStringDigit(resp.CreatedAt.Time),
		UpdatedAt:         dateutils.DateTimeToStringDigit(resp.UpdatedAt.Time),
		DeletedAt:         dateutils.DateTimeToStringDigit(resp.DeletedAt.Time),
	}
}
func (a *AccountsAdapter) RoleListGrpcFromSql(resp *[]db.RoleListRow) *talv1.RoleListResponse {
	records := make([]*talv1.RoleListRow, 0)
	deletedRecords := make([]*talv1.RoleListRow, 0)
	for _, v := range *resp {
		record := a.RoleListRowGrpcFromSql(&v)
		if v.DeletedAt.Valid {
			deletedRecords = append(deletedRecords, record)
			continue
		}
		records = append(records, record)
	}
	return &talv1.RoleListResponse{
		DeletedRecords: deletedRecords,
		Records:        records,
	}
}
func (a *AccountsAdapter) RoleCreateUpdateGrpcFromSql(resp *db.AccountsSchemaRole) *talv1.RoleCreateUpdateResponse {
	return &talv1.RoleCreateUpdateResponse{
		Role: a.RoleEntityGrpcFromSql(resp),
	}
}
