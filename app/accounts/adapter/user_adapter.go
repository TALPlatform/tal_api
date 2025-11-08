package adapter

import (
	"encoding/json"

	"github.com/TALPlatform/tal_api/db"
	"github.com/TALPlatform/tal_api/pkg/dateutils"
	"github.com/TALPlatform/tal_api/pkg/redisclient"
	talv1 "github.com/TALPlatform/tal_api/proto_gen/tal/v1"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

func (a *AccountsAdapter) UserEntityGrpcFromSql(resp *db.AccountsSchemaUser) *talv1.AccountsSchemaUser {
	return &talv1.AccountsSchemaUser{
		UserId:     int32(resp.UserID),
		UserName:   resp.UserName,
		UserImage:  resp.UserImage.String,
		UserTypeId: resp.UserTypeID,
		TenantId:   resp.TenantID.Int32,
		UserPhone:  resp.UserPhone.String,
		UserEmail:  resp.UserEmail, // User's email, unique in DB
		CreatedAt:  db.TimeToString(resp.CreatedAt.Time),
		DeletedAt:  db.TimeToString(resp.DeletedAt.Time),
	}
}

func (a *AccountsAdapter) UserFindRowGrpcFromSql(resp *db.UserFindRow) *talv1.UserFindRow {
	record := &talv1.UserFindRow{
		UserId:            int32(resp.UserID),
		UserImage:         resp.UserImage,
		UserName:          resp.UserName,
		UserTypeId:        resp.UserTypeID,
		UserTypeName:      resp.UserTypeName,
		UserSecurityLevel: resp.UserSecurityLevel,
		PermissionCount:   int32(resp.PermissionCount),
		TenantName:        resp.TenantName,
		TenantId:          resp.TenantID,
		UserPhone:         resp.UserPhone,
		UserEmail:         resp.UserEmail,
		CreatedAt:         dateutils.DateTimeToStringDigit(resp.CreatedAt.Time),
		UpdatedAt:         dateutils.DateTimeToStringDigit(resp.UpdatedAt.Time),
		DeletedAt:         dateutils.DateTimeToStringDigit(resp.DeletedAt.Time),
	}
	if len(resp.Roles) > 0 {
		json.Unmarshal(resp.Roles, &record.Roles)
	}
	if len(resp.Logs) > 0 {
		json.Unmarshal(resp.Logs, &record.Logs)
	}
	return record
}
func (a *AccountsAdapter) UserSessionGrpcFropmSql(resp *redisclient.AuthSession) *talv1.UserSession {
	return &talv1.UserSession{
		UserId:                        int32(resp.UserID),
		SessionKey:                    resp.SessionKey,
		IpAddress:                     resp.IPAddress,
		IsBlocked:                     resp.IsBlocked,
		CreatedAt:                     dateutils.DateTimeToStringDigit(resp.CreatedAt),
		AccessTokenExpiresAt:          dateutils.DateTimeToStringDigit(resp.AccessTokenExpiresAt),
		RefreshTokenExpiresAt:         dateutils.DateTimeToStringDigit(resp.RefreshTokenExpiresAt),
		SupabaseAccessTokenExpiresAt:  dateutils.DateTimeToStringDigit(resp.SupabaseAccessTokenExpiresAt),
		SupabaseRefreshTokenExpiresAt: dateutils.DateTimeToStringDigit(resp.SupabaseRefreshTokenExpiresAt),
	}
}
func (a *AccountsAdapter) UserSessionsGrpcFropmSql(sessions []*redisclient.AuthSession) []*talv1.UserSession {
	response := make([]*talv1.UserSession, len(sessions))

	for index, session := range sessions {
		response[index] = a.UserSessionGrpcFropmSql(session)
	}
	return response
}
func (a *AccountsAdapter) UserViewEntityGrpcFromSql(resp *db.AccountsSchemaUserView) *talv1.AccountsSchemaUserView {
	record := &talv1.AccountsSchemaUserView{
		UserId:       int32(resp.UserID),
		UserImage:    resp.UserImage,
		UserName:     resp.UserName,
		UserTypeId:   resp.UserTypeID,
		UserTypeName: resp.UserTypeName,

		UserSecurityLevel: resp.UserSecurityLevel,
		TenantName:        resp.TenantName,
		TenantId:          resp.TenantID,
		UserPhone:         resp.UserPhone,
		UserEmail:         resp.UserEmail,
		CreatedAt:         dateutils.DateTimeToStringDigit(resp.CreatedAt.Time),
		UpdatedAt:         dateutils.DateTimeToStringDigit(resp.UpdatedAt.Time),
		DeletedAt:         dateutils.DateTimeToStringDigit(resp.DeletedAt.Time),
	}
	if len(resp.Roles) > 0 {
		json.Unmarshal(resp.Roles, &record.Roles)
	}
	return record
}
func (a *AccountsAdapter) UserCreateUpdateSqlFromGrpc(req *talv1.UserCreateUpdateRequest) *db.UserCreateUpdateParams {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.UserPassword), bcrypt.DefaultCost)
	resp := &db.UserCreateUpdateParams{
		UserID:       req.UserId,
		UserName:     req.UserName,
		UserImage:    req.UserImage,
		UserTypeID:   req.UserTypeId,
		UserPhone:    req.UserPhone,
		UserEmail:    req.UserEmail,
		UserPassword: string(hashedPassword),
		Roles:        req.Roles,
	}
	return resp
}
func (a *AccountsAdapter) UserFindForUpdateUpdateGrpcFromSql(resp *db.UserFindForUpdateRow) *talv1.UserCreateUpdateRequest {
	return &talv1.UserCreateUpdateRequest{
		UserId:     resp.UserID,
		TenantId:   resp.TenantID.Int32,
		UserName:   resp.UserName,
		UserImage:  resp.UserImage.String,
		UserTypeId: resp.UserTypeID,
		UserPhone:  resp.UserPhone.String,
		UserEmail:  resp.UserEmail,
		Roles:      resp.Roles,
	}
}
func (a *AccountsAdapter) UserTypeListInputGrpcFromSql(resp *[]db.UserTypeListInputRow) *talv1.UserTypeListInputResponse {
	records := make([]*talv1.SelectInputOption, 0)
	for _, v := range *resp {
		records = append(records, &talv1.SelectInputOption{
			Value: v.Value,
			Label: v.Label,
		})
	}
	return &talv1.UserTypeListInputResponse{
		Options: records,
	}
}

func (a *AccountsAdapter) UserListInputGrpcFromSql(resp *[]db.UserListInputRow) *talv1.UserListInputResponse {
	records := make([]*talv1.SelectInputOption, 0)
	for _, v := range *resp {
		records = append(records, &talv1.SelectInputOption{
			Value: v.Value,
			Note:  v.Note,
			Label: v.Label,
		})
	}
	return &talv1.UserListInputResponse{
		Options: records,
	}
}
func (a *AccountsAdapter) UserPermissionListInputGrpcFromSql(resp *[]db.UserPermissionListInputRow) *talv1.UserPermissionListInputResponse {
	groupedOptions := make([]*talv1.SelectInputOptionWithGroup, len(*resp))

	for groupIndex, v := range *resp {
		items := make([]*talv1.SelectInputOption, len(v.Options))
		if err := json.Unmarshal(v.Options, &items); err != nil {
			continue
		}
		groupedOptions[groupIndex] = &talv1.SelectInputOptionWithGroup{
			GroupName: v.PermissionGroup,
			Items:     items,
		}
	}
	return &talv1.UserPermissionListInputResponse{
		Options: groupedOptions,
	}
}

func (a *AccountsAdapter) UserListRowGrpcFromSql(resp *db.AccountsSchemaUserView) *talv1.UserListRow {
	record := &talv1.UserListRow{
		UserId:            int32(resp.UserID),
		UserImage:         resp.UserImage,
		UserName:          resp.UserName,
		UserTypeId:        resp.UserTypeID,
		UserTypeName:      resp.UserTypeName,
		UserSecurityLevel: resp.UserSecurityLevel,
		TenantName:        resp.TenantName,
		TenantId:          resp.TenantID,
		UserPhone:         resp.UserPhone,
		UserEmail:         resp.UserEmail,
		CreatedAt:         dateutils.DateTimeToStringDigit(resp.CreatedAt.Time),
		UpdatedAt:         dateutils.DateTimeToStringDigit(resp.UpdatedAt.Time),
		DeletedAt:         dateutils.DateTimeToStringDigit(resp.DeletedAt.Time),
	}
	if len(resp.Roles) > 0 {
		err := json.Unmarshal(resp.Roles, &record.Roles)
		if err != nil {
			log.Error().Err(err).Msg("error parsinng user roles into json")
		}
	}
	roleIds := make([]int32, len(record.Roles))
	for index, v := range record.Roles {
		roleIds[index] = v.RoleId
	}
	record.RoleIds = roleIds
	return record
}
func (a *AccountsAdapter) UserListGrpcFromSql(resp *[]db.AccountsSchemaUserView) *talv1.UserListResponse {
	records := make([]*talv1.UserListRow, 0)
	deletedRecords := make([]*talv1.UserListRow, 0)
	for _, v := range *resp {
		record := a.UserListRowGrpcFromSql(&v)
		if v.DeletedAt.Valid {
			deletedRecords = append(deletedRecords, record)
		} else {
			records = append(records, record)
		}
	}
	return &talv1.UserListResponse{
		DeletedRecords: deletedRecords,
		Records:        records,
	}
}
func (a *AccountsAdapter) UserCreateUpdateGrpcFromSql(resp *db.AccountsSchemaUser) *talv1.UserCreateUpdateResponse {
	return &talv1.UserCreateUpdateResponse{
		User: a.UserEntityGrpcFromSql(resp),
	}
}
