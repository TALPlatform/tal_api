package usecase

import (
	"context"
	"time"

	"connectrpc.com/connect"
	"github.com/TALPlatform/tal_api/app/accounts/adapter"
	"github.com/TALPlatform/tal_api/app/accounts/repo"
	"github.com/TALPlatform/tal_api/db"
	"github.com/TALPlatform/tal_api/pkg/auth"
	"github.com/TALPlatform/tal_api/pkg/redisclient"
	talv1 "github.com/TALPlatform/tal_api/proto_gen/tal/v1"
	supaapigo "github.com/darwishdev/supaapi-go"
)

type AccountsUsecaseInterface interface {
	// CheckForAccess(ctx context.Context, header http.Header, functionName string, isCreateUpdate bool) (*talv1.AvailableOptions, error)
	AuthLogin(ctx context.Context, req *connect.Request[talv1.AuthLoginRequest]) (*talv1.AuthLoginResponse, error)
	AuthSessionSetBlocked(
		ctx context.Context,
		req *connect.Request[talv1.AuthSessionSetBlockedRequest],
	) (*talv1.AuthSessionSetBlockedResponse, error)
	AuthSessionDelete(
		ctx context.Context,
		req *connect.Request[talv1.AuthSessionDeleteRequest],
	) (*talv1.AuthSessionDeleteResponse, error)
	AuthSessionList(ctx context.Context, req *connect.Request[talv1.AuthSessionListRequest]) (*talv1.AuthSessionListResponse, error)
	AuthRegister(ctx context.Context, req *connect.Request[talv1.AuthRegisterRequest]) (*talv1.AuthRegisterResponse, error)
	UserDeleteRestore(ctx context.Context, req *connect.Request[talv1.UserDeleteRestoreRequest]) (*talv1.UserDeleteRestoreResponse, error)
	UserPermissionListInput(ctx context.Context, req *connect.Request[talv1.UserPermissionListInputRequest]) (*talv1.UserPermissionListInputResponse, error)
	UserList(ctx context.Context) (*talv1.UserListResponse, error)
	UserTypeListInput(ctx context.Context) (*talv1.UserTypeListInputResponse, error)
	UserListInput(ctx context.Context) (*talv1.UserListInputResponse, error)
	UserFindForUpdate(ctx context.Context, req *connect.Request[talv1.UserFindForUpdateRequest]) (*talv1.UserFindForUpdateResponse, error)
	UserCreateUpdate(ctx context.Context, req *connect.Request[talv1.UserCreateUpdateRequest]) (*talv1.UserCreateUpdateResponse, error)
	AuthLoginProviderCallback(ctx context.Context, req *connect.Request[talv1.AuthLoginProviderCallbackRequest]) (*talv1.AuthLoginProviderCallbackResponse, error)
	AuthLoginProvider(ctx context.Context, req *connect.Request[talv1.AuthLoginProviderRequest]) (*talv1.AuthLoginProviderResponse, error)
	UserDelete(ctx context.Context, req *connect.Request[talv1.UserDeleteRequest]) (*talv1.UserDeleteResponse, error)
	AuthInvite(ctx context.Context, req *connect.Request[talv1.AuthInviteRequest]) (*talv1.AuthInviteResponse, error)
	RoleDelete(ctx context.Context, req *connect.Request[talv1.RoleDeleteRequest]) (*talv1.RoleDeleteResponse, error)
	RoleDeleteRestore(ctx context.Context, req *connect.Request[talv1.RoleDeleteRestoreRequest]) (*talv1.RoleDeleteRestoreResponse, error)
	RoleFindForUpdate(ctx context.Context, req *connect.Request[talv1.RoleFindForUpdateRequest]) (*talv1.RoleFindForUpdateResponse, error)
	RoleListInput(ctx context.Context) (*talv1.RoleListInputResponse, error)
	RoleList(ctx context.Context, req *connect.Request[talv1.RoleListRequest]) (*talv1.RoleListResponse, error)
	UserGenerateTokens(username string, userId int32, tenantId int32, userSecurityLevel int32) (*talv1.LoginInfo, string, string, error)
	AuthLogout(
		ctx context.Context,
		req *connect.Request[talv1.AuthLogoutRequest],
	) (*talv1.AuthLogoutResponse, error)
	AppLogin(ctx context.Context, loginCode string, userId int32) (*talv1.AuthLoginResponse, error)

	UserFind(ctx context.Context, req *connect.Request[talv1.UserFindRequest]) (*talv1.UserFindResponse, error)
	AuthResetPassword(ctx context.Context, req *connect.Request[talv1.AuthResetPasswordRequest]) (*talv1.AuthResetPasswordResponse, error)
	AuthResetPasswordEmail(ctx context.Context, req *connect.Request[talv1.AuthResetPasswordEmailRequest]) (*talv1.AuthResetPasswordEmailResponse, error)
	AuthRefreshToken(ctx context.Context, req *connect.Request[talv1.AuthRefreshTokenRequest]) (*talv1.AuthRefreshTokenResponse, error)
	RoleCreateUpdate(ctx context.Context, req *connect.Request[talv1.RoleCreateUpdateRequest]) (*talv1.RoleCreateUpdateResponse, error)
}

type AccountsUsecase struct {
	store                db.Store
	adapter              adapter.AccountsAdapterInterface
	tokenMaker           auth.Maker
	tokenDuration        time.Duration
	refreshTokenDuration time.Duration
	supaapi              supaapigo.Supaapi
	redisClient          redisclient.RedisClientInterface
	repo                 repo.AccountsRepoInterface
}

func NewAccountsUsecase(store db.Store, supaapi supaapigo.Supaapi, redisClient redisclient.RedisClientInterface, tokenMaker auth.Maker, tokenDuration time.Duration, refreshTokenDuration time.Duration) AccountsUsecaseInterface {
	return &AccountsUsecase{
		supaapi:              supaapi,
		tokenMaker:           tokenMaker,
		redisClient:          redisClient,
		tokenDuration:        tokenDuration,
		refreshTokenDuration: refreshTokenDuration,
		store:                store,
		adapter:              adapter.NewAccountsAdapter(),
		repo:                 repo.NewAccountsRepo(store),
	}
}
