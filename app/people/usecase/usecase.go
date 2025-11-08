package usecase

import (
	"context"

	"connectrpc.com/connect"
	"github.com/TALPlatform/tal_api/app/people/adapter"
	"github.com/TALPlatform/tal_api/app/people/repo"
	"github.com/TALPlatform/tal_api/db"
	"github.com/TALPlatform/tal_api/pkg/crustdata"
	talv1 "github.com/TALPlatform/tal_api/proto_gen/tal/v1"
	"github.com/darwishdev/genaiclient"
)

type PeopleUsecaseInterface interface {
	// INJECT INTERFACE
	RawProfileList(ctx context.Context, req *connect.Request[talv1.RawProfileListRequest]) (*talv1.RawProfileListResponse, error)
}

type PeopleUsecase struct {
	store           db.Store
	adapter         adapter.PeopleAdapterInterface
	llmClient       genaiclient.GenaiClientInterface
	crustDataClient crustdata.CrustdataServiceInterface
	repo            repo.PeopleRepoInterface
}

func NewPeopleUsecase(store db.Store, crustDataClient crustdata.CrustdataServiceInterface, llmClient genaiclient.GenaiClientInterface) PeopleUsecaseInterface {
	return &PeopleUsecase{
		store:           store,
		crustDataClient: crustDataClient,
		llmClient:       llmClient,
		adapter:         adapter.NewPeopleAdapter(llmClient.Embed, llmClient.EmbedBulk),
		repo:            repo.NewPeopleRepo(store),
	}
}
