package redisclient

import (
	"context"
	"fmt"
	"time"

	"github.com/TALPlatform/tal_api/db"
	"github.com/redis/go-redis/v9"
)

type RedisClientInterface interface {
	UserPermissionCreate(ctx context.Context, userID int32, perms *PermissionsMap) error
	UserPermissionDelete(ctx context.Context, userID int32) error
	UserPermissionFind(ctx context.Context, userID int32) (PermissionsMap, error)
	GenerateAuthSessionKey(userID int32, tokenID string) string
	AuthSessionCreate(ctx context.Context, session *AuthSession, refreshTokenID string, ttl time.Duration) error
	AuthSessionCreateByKey(ctx context.Context, session *AuthSession, key string, ttl time.Duration) error
	AuthSessionFindByAccessTokenID(ctx context.Context, accessTokenID string) (*AuthSession, error)
	AuthSessionUpdateTokens(
		ctx context.Context,
		userID int32,
		refreshTokenID string,
		accessToken string,
		accessTokenExpiresAt time.Time,
		refreshToken string,
		refreshTokenExpiresAt time.Time,
		supabaseAccessToken string,
		supabaseAccessTokenExpiresAt time.Time,
		supabaseRefreshToken string,
		supabaseRefreshTokenExpiresAt time.Time,
	) error

	AuthSessionFind(ctx context.Context, userID int32, refreshTokenID string) (*AuthSession, error)
	AuthSessionFindByKey(ctx context.Context, key string) (*AuthSession, error)
	AuthSessionListAll(ctx context.Context) ([]*AuthSession, error)
	AuthSessionListByUser(ctx context.Context, userID int32) ([]*AuthSession, error)

	AuthSessionSetBlocked(ctx context.Context, userID int32, refreshTokenID string, blocked bool) error
	AuthSessionSetBlockedByKey(ctx context.Context, session_id string, blocked bool) error

	AuthSessionDelete(ctx context.Context, userID int32, refreshTokenID string) error
	AuthSessionDeleteByKey(ctx context.Context, key string) error
	AuthSessionDeleteByUserAgent(ctx context.Context, userID int32, userAgent string) error
	AuthSessionClearAll(ctx context.Context, userID int32) error

	GetClient() *redis.Client
	TenantCreate(ctx context.Context, tenantId int32, tenant *db.TenantFindRow) (*db.TenantFindRow, error)
	DeleteAllTenants(ctx context.Context) error
	TenantDelete(ctx context.Context, tenantId int32) error
	TenantDeleteBulk(ctx context.Context, tenantIds []int32) error
	TenantFind(ctx context.Context, tenantId int32) (*db.TenantFindRow, error)
}

type RedisClient struct {
	client     *redis.Client
	isDisabled bool
}

func NewRedisClient(host string, port string, password string, db int, isDisabled bool) RedisClientInterface {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password, // no password set
		DB:       db,       // use default DB

	})

	return &RedisClient{
		client:     client,
		isDisabled: isDisabled,
	}
}
func (r *RedisClient) GetClient() *redis.Client {
	return r.client
}
