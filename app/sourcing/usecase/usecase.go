package usecase

import (
	"context"

	"connectrpc.com/connect"
	"github.com/TALPlatform/tal_api/app/sourcing/adapter"
	"github.com/TALPlatform/tal_api/app/sourcing/repo"
	"github.com/TALPlatform/tal_api/db"
	talv1 "github.com/TALPlatform/tal_api/proto_gen/tal/v1"
	"github.com/darwishdev/genaiclient"
)

type SourcingUsecaseInterface interface {
	// INJECT INTERFACE

SourcingSessionProfileCreateUpdate(ctx context.Context, req *connect.Request[talv1.SourcingSessionProfileCreateUpdateRequest]) (*talv1.SourcingSessionProfileCreateUpdateResponse, error) 


SourcingSessionFind(ctx context.Context, req *connect.Request[talv1.SourcingSessionFindRequest]) (*talv1.SourcingSessionFindResponse, error) 


SourcingSessionCreateUpdate(ctx context.Context, req *connect.Request[talv1.SourcingSessionCreateUpdateRequest]) (*talv1.SourcingSessionCreateUpdateResponse, error) 


	ProjectInputList(ctx context.Context, req *connect.Request[talv1.ProjectInputListRequest]) (*talv1.ProjectInputListResponse, error)
}

type SourcingUsecase struct {
	store     db.Store
	adapter   adapter.SourcingAdapterInterface
	llmClient genaiclient.GenaiClientInterface
	repo      repo.SourcingRepoInterface
}

func NewSourcingUsecase(store db.Store, llmClient genaiclient.GenaiClientInterface) SourcingUsecaseInterface {
	return &SourcingUsecase{
		store:     store,
		llmClient: llmClient,
		adapter:   adapter.NewSourcingAdapter(),
		repo:      repo.NewSourcingRepo(store),
	}
}
