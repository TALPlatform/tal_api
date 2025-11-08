package usecase

import (
	"context"

	"connectrpc.com/connect"
	"github.com/TALPlatform/tal_api/app/public/adapter"
	"github.com/TALPlatform/tal_api/app/public/repo"
	"github.com/TALPlatform/tal_api/db"
	"github.com/TALPlatform/tal_api/pkg/redisclient"
	"github.com/TALPlatform/tal_api/pkg/resend"
	"github.com/TALPlatform/tal_api/pkg/weaviateclient"
	talv1 "github.com/TALPlatform/tal_api/proto_gen/tal/v1"
	supaapigo "github.com/darwishdev/supaapi-go"
)

type PublicUsecaseInterface interface {
	TranslationList(ctx context.Context) (*talv1.TranslationListResponse, error)
	TranslationCreateUpdateBulk(ctx context.Context, req *talv1.TranslationCreateUpdateBulkRequest) (*talv1.TranslationCreateUpdateBulkResponse, error)
	TranslationFindLocale(ctx context.Context, req *talv1.TranslationFindLocaleRequest) (*talv1.TranslationFindLocaleResponse, error)
	TranslationDelete(ctx context.Context, req *talv1.TranslationDeleteRequest) (*talv1.TranslationDeleteResponse, error)
	GalleryList(ctx context.Context, req *talv1.GalleryListRequest) (*talv1.GalleryListResponse, error)
	FileDeleteByBucket(ctx context.Context, req *talv1.FileDeleteByBucketRequest) (*talv1.FileDeleteByBucketResponse, error)
	FileDelete(ctx context.Context, req *talv1.FileDeleteRequest) (*talv1.FileDeleteResponse, error)
	FileList(ctx context.Context, req *talv1.FileListRequest) (*talv1.FileListResponse, error)
	EmailSend(ctx context.Context, req *talv1.EmailSendRequest) (*talv1.EmailSendResponse, error)
	BucketList(ctx context.Context, req *talv1.BucketListRequest) (*talv1.BucketListResponse, error)
	SettingUpdate(ctx context.Context, req *talv1.SettingUpdateRequest) error
	SettingFindForUpdate(ctx context.Context, req *talv1.SettingFindForUpdateRequest) (*talv1.SettingFindForUpdateResponse, error)
	FileCreate(ctx context.Context, req *talv1.FileCreateRequest) (*talv1.FileCreateResponse, error)
	BucketCreateUpdate(ctx context.Context, req *talv1.BucketCreateUpdateRequest) (*talv1.BucketCreateUpdateResponse, error)

	IconFind(ctx context.Context, req *talv1.IconFindRequest) (*talv1.IconFindResponse, error)
	IconCreateUpdateBulk(ctx context.Context, req *talv1.IconCreateUpdateBulkRequest) (*talv1.IconListResponse, error)
	IconList(ctx context.Context) (*talv1.IconListResponse, error)
	FileUploadUrlFind(ctx context.Context, req *talv1.FileUploadUrlFindRequest) (*talv1.FileUploadUrlFindResponse, error)
	FileCreateBulk(ctx context.Context, req *talv1.FileCreateBulkRequest) (*talv1.FileCreateBulkResponse, error)

	CommandPalleteSync(ctx context.Context, req *connect.Request[talv1.CommandPalleteSyncRequest]) (*talv1.CommandPalleteSyncResponse, error)
 CommandPalleteSearch(ctx context.Context, req *connect.Request[talv1.CommandPalleteSearchRequest]) (*talv1.CommandPalleteSearchResponse, error) 
}

type PublicUsecase struct {
	store          db.Store
	repo           repo.PublicRepoInterface
	adapter        adapter.PublicAdapterInterface
	supaapi        supaapigo.Supaapi
	supaAnonApiKey string
	resendClient   resend.ResendServiceInterface
	weaviateClient weaviateclient.WeaviateClientInterface
	redisClient    redisclient.RedisClientInterface
}

func NewPublicUsecase(store db.Store, supaAnonApiKey string, supaapi supaapigo.Supaapi, redisClient redisclient.RedisClientInterface, resendClient resend.ResendServiceInterface, weaviateClient weaviateclient.WeaviateClientInterface) PublicUsecaseInterface {
	return &PublicUsecase{
		resendClient:   resendClient,
		supaAnonApiKey: supaAnonApiKey,
		supaapi:        supaapi,
		redisClient:    redisClient,
		weaviateClient: weaviateClient,
		adapter:        adapter.NewPublicAdapter(),
		repo:           repo.NewPublicRepo(store),
		store:          store,
	}
}
