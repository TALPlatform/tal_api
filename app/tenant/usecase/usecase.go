package usecase

import (
	"context"

	"connectrpc.com/connect"
	"github.com/TALPlatform/tal_api/app/tenant/adapter"
	"github.com/TALPlatform/tal_api/app/tenant/repo"
	"github.com/TALPlatform/tal_api/db"
	"github.com/TALPlatform/tal_api/pkg/redisclient"
	talv1 "github.com/TALPlatform/tal_api/proto_gen/tal/v1"
)

type TenantUsecaseInterface interface {
	// INJECT INTERFACE

	SectionListInpt(ctx context.Context, req *connect.Request[talv1.SectionListInptRequest]) (*talv1.SectionListInptResponse, error)

	PartialTypeListInput(ctx context.Context, req *connect.Request[talv1.PartialTypeListInputRequest]) (*talv1.PartialTypeListInputResponse, error)

	PartialDeleteRestore(ctx context.Context, req *connect.Request[talv1.PartialDeleteRestoreRequest]) (*talv1.PartialDeleteRestoreResponse, error)
	PartialCreateUpdate(ctx context.Context, req *connect.Request[talv1.PartialCreateUpdateRequest]) (*talv1.PartialCreateUpdateResponse, error)
	PartialList(ctx context.Context, req *connect.Request[talv1.PartialListRequest]) (*talv1.PartialListResponse, error)

	PartialFindForUpdate(ctx context.Context, req *connect.Request[talv1.PartialFindForUpdateRequest]) (*talv1.PartialFindForUpdateResponse, error)

	PageFindForUpdate(ctx context.Context, req *connect.Request[talv1.PageFindForUpdateRequest]) (*talv1.PageFindForUpdateResponse, error)

	SectionFindForUpdate(ctx context.Context, req *connect.Request[talv1.SectionFindForUpdateRequest]) (*talv1.SectionFindForUpdateResponse, error)
	SectionDeleteRestore(ctx context.Context, req *connect.Request[talv1.SectionDeleteRestoreRequest]) (*talv1.SectionDeleteRestoreResponse, error)
	SectionCreateUpdate(ctx context.Context, req *connect.Request[talv1.SectionCreateUpdateRequest]) (*talv1.SectionCreateUpdateResponse, error)
	SectionList(ctx context.Context, req *connect.Request[talv1.SectionListRequest]) (*talv1.SectionListResponse, error)

	PageDeleteRestore(ctx context.Context, req *connect.Request[talv1.PageDeleteRestoreRequest]) (*talv1.PageDeleteRestoreResponse, error)

	PageCreateUpdate(ctx context.Context, req *connect.Request[talv1.PageCreateUpdateRequest]) (*talv1.PageCreateUpdateResponse, error)

	PageList(ctx context.Context, req *connect.Request[talv1.PageListRequest]) (*talv1.PageListResponse, error)

	TenantDeleteRestore(ctx context.Context, req *connect.Request[talv1.TenantDeleteRestoreRequest]) (*talv1.TenantDeleteRestoreResponse, error)
	TenantList(ctx context.Context, req *connect.Request[talv1.TenantListRequest]) (*talv1.TenantListResponse, error)
	TenantListInput(ctx context.Context, req *connect.Request[talv1.TenantListInputRequest]) (*talv1.TenantListInputResponse, error)
	TenantFind(ctx context.Context, req *connect.Request[talv1.TenantFindRequest]) (*talv1.TenantFindResponse, error)
	TenantCreateUpdate(ctx context.Context, req *connect.Request[talv1.TenantCreateUpdateRequest]) (*talv1.TenantCreateUpdateResponse, error)
	TenantDashboard(ctx context.Context, req *connect.Request[talv1.TenantDashboardRequest]) (*talv1.TenantDashboardResponse, error)
}
type TenantUsecase struct {
	store       db.Store
	adapter     adapter.TenantAdapterInterface
	redisClient redisclient.RedisClientInterface
	repo        repo.TenantRepoInterface
}

func NewTenantUsecase(store db.Store, redisClient redisclient.RedisClientInterface) TenantUsecaseInterface {
	return &TenantUsecase{
		store:       store,
		redisClient: redisClient,
		adapter:     adapter.NewTenantAdapter(),
		repo:        repo.NewTenantRepo(store),
	}
}
