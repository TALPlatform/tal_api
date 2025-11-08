package repo

import (
	"context"

	"github.com/TALPlatform/tal_api/db"
)

type AccountsRepoInterface interface {
	UserNavigationBarFind(ctx context.Context, req db.UserNavigationBarFindParams) (*[]db.UserNavigationBarFindRow, error)
	UserTypeListInput(ctx context.Context) (*[]db.UserTypeListInputRow, error)
	UserListInput(ctx context.Context) (*[]db.UserListInputRow, error)
	UserPermissionListInput(ctx context.Context) (*[]db.UserPermissionListInputRow, error)
	UseriFindForUpdate(ctx context.Context, req int32) (*db.UserFindForUpdateRow, error)
	UserFind(ctx context.Context, userID int32) (*db.UserFindRow, error)
	UserDeleteRestore(ctx context.Context, req db.UserDeleteRestoreParams) (*db.AccountsSchemaUser, error)
	UserDelete(ctx context.Context, req db.UserDeleteParams) (*db.AccountsSchemaUser, error)
	UserPermissionsMap(ctx context.Context, userID int32) (*[]db.UserPermissionsMapRow, error)
	UserCreateUpdate(ctx context.Context, req db.UserCreateUpdateParams) (*db.AccountsSchemaUser, error)
	UserList(ctx context.Context) (*[]db.AccountsSchemaUserView, error)
	RoleListInput(ctx context.Context) (*[]db.RoleListInputRow, error)
	RoleFindForUpdate(ctx context.Context, req int32) (*db.RoleFindForUpdateRow, error)
	UserFindForToken(ctx context.Context, req *db.UserFindForTokenParams) (*db.UserFindForTokenRow, error)
	UserFindForAuth(ctx context.Context, params db.UserFindForAuthParams) (*db.AccountsSchemaUserView, error)
	RoleList(ctx context.Context) (*[]db.RoleListRow, error)
	AuthUserIDFindByEmail(ctx context.Context, req string) (*string, error)
	RoleCreateUpdate(ctx context.Context, req db.RoleCreateUpdateParams) (*db.AccountsSchemaRole, error)
	RoleDelete(ctx context.Context, req db.RoleDeleteParams) (*db.AccountsSchemaRole, error)
	RoleDeleteRestore(ctx context.Context, req db.RoleDeleteRestoreParams) (*db.AccountsSchemaRole, error)
}

type AccountsRepo struct {
	store        db.Store
	errorHandler map[string]string
}

func NewAccountsRepo(store db.Store) AccountsRepoInterface {
	errorHandler := map[string]string{
		"role_role_name_key":  "roleName",
		"user_user_name_key":  "userName",
		"user_user_email_key": "userEmail",
		"user_user_phone_key": "userPhone",
	}
	return &AccountsRepo{
		store:        store,
		errorHandler: errorHandler,
	}
}
