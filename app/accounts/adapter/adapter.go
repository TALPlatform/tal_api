package adapter

import (
	"github.com/TALPlatform/tal_api/db"
	"github.com/TALPlatform/tal_api/pkg/redisclient"
	talv1 "github.com/TALPlatform/tal_api/proto_gen/tal/v1"
	"github.com/supabase-community/auth-go/types"
)

type AccountsAdapterInterface interface {
	UserPermissionsMapRedisFromSql(resp *[]db.UserPermissionsMapRow) (*redisclient.PermissionsMap, error)
	AuttSessionRedisFromGrpc(response *talv1.AuthLoginResponse, ipAddress string, userAgent string) (*redisclient.AuthSession, error)
	AuthLoginSqlFromGrpc(req *talv1.AuthLoginRequest) (*db.UserFindForAuthParams, *types.TokenRequest)
	UserCreateUpdateRequestFromAuthRegister(req *talv1.AuthRegisterRequest) *talv1.UserCreateUpdateRequest
	AuthSessionListGrpcFromRedis(resp []*redisclient.AuthSession) *talv1.AuthSessionListResponse
	AuthResetPasswordSupaFromGrpc(req *talv1.AuthResetPasswordRequest) *types.VerifyForUserRequest
	AuthLoginGrpcFromSql(resp *db.AccountsSchemaUserView) *talv1.AuthLoginResponse
	UserNavigationBarFindGrpcFromSql(dbResponse []db.UserNavigationBarFindRow, sourcingSessions *[]db.SourcingSessionListRow) ([]*talv1.NavigationBarItem, error)
	UserCreateUpdateGrpcFromSql(resp *db.AccountsSchemaUser) *talv1.UserCreateUpdateResponse
	UserFindRowGrpcFromSql(resp *db.UserFindRow) *talv1.UserFindRow
	UserSessionsGrpcFropmSql(sessions []*redisclient.AuthSession) []*talv1.UserSession
	UserFindForUpdateUpdateGrpcFromSql(resp *db.UserFindForUpdateRow) *talv1.UserCreateUpdateRequest
	UserTypeListInputGrpcFromSql(resp *[]db.UserTypeListInputRow) *talv1.UserTypeListInputResponse
	UserPermissionListInputGrpcFromSql(resp *[]db.UserPermissionListInputRow) *talv1.UserPermissionListInputResponse
	UserListInputGrpcFromSql(resp *[]db.UserListInputRow) *talv1.UserListInputResponse
	NavigationBarItemGrpcFromSql(resp *db.UserNavigationBarFindRow) *talv1.NavigationBarItem
	UserListGrpcFromSql(resp *[]db.AccountsSchemaUserView) *talv1.UserListResponse
	UserCreateUpdateSqlFromGrpc(req *talv1.UserCreateUpdateRequest) *db.UserCreateUpdateParams
	UserViewEntityGrpcFromSql(resp *db.AccountsSchemaUserView) *talv1.AccountsSchemaUserView
	UserEntityGrpcFromSql(resp *db.AccountsSchemaUser) *talv1.AccountsSchemaUser
	RoleFindForUpdateUpdateGrpcFromSql(resp *db.RoleFindForUpdateRow) *talv1.RoleCreateUpdateRequest
	RoleListInputGrpcFromSql(resp *[]db.RoleListInputRow) *talv1.RoleListInputResponse
	RoleListGrpcFromSql(resp *[]db.RoleListRow) *talv1.RoleListResponse
	RoleEntityGrpcFromSql(resp *db.AccountsSchemaRole) *talv1.AccountsSchemaRole
	RoleCreateUpdateSqlFromGrpc(req *talv1.RoleCreateUpdateRequest) *db.RoleCreateUpdateParams
	RoleCreateUpdateGrpcFromSql(resp *db.AccountsSchemaRole) *talv1.RoleCreateUpdateResponse
}

type AccountsAdapter struct {
}

func NewAccountsAdapter() AccountsAdapterInterface {
	return &AccountsAdapter{}
}
