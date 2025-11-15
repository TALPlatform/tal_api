package usecase

import (
	"context"
	"os"

	"connectrpc.com/connect"
	"github.com/TALPlatform/tal_api/app/sourcing/adapter"
	"github.com/TALPlatform/tal_api/app/sourcing/repo"
	"github.com/TALPlatform/tal_api/db"
	"github.com/TALPlatform/tal_api/pkg/agenthub"
	"github.com/TALPlatform/tal_api/pkg/asynqclient"
	"github.com/TALPlatform/tal_api/pkg/crustdata"
	talv1 "github.com/TALPlatform/tal_api/proto_gen/tal/v1"
	"github.com/darwishdev/genaiclient"

	"github.com/redis/go-redis/v9"
)

type SourcingUsecaseInterface interface {
	// INJECT INTERFACE

	SourcingSessionFiltersInfered(ctx context.Context, req *connect.Request[talv1.SourcingSessionFiltersInferedRequest]) (*talv1.SourcingSessionFiltersInferedResponse, error)

	SourcingSessionFiltersBuilder(ctx context.Context, req *connect.Request[talv1.SourcingSessionFiltersBuilderRequest]) (*talv1.SourcingSessionFiltersBuilderResponse, error)
	SourcingSessionApply(ctx context.Context, req *connect.Request[talv1.SourcingSessionApplyRequest], stream *connect.ServerStream[talv1.SourcingSessionApplyResponse]) error
	SyncSourcingSessionProfiles(
		ctx context.Context,
		sessionID int32,
		crustDataProfiles *[]byte,
		DBProfiles *[]byte,
	) error
	SourcingSessionCriteriaCreate(ctx context.Context, req *connect.Request[talv1.SourcingSessionCriteriaCreateRequest], stream *connect.ServerStream[talv1.SourcingSessionCriteriaCreateResponse]) error

	SourcingSessionProfileJustify(ctx context.Context, req *connect.Request[talv1.SourcingSessionProfileJustifyRequest], stream *connect.ServerStream[talv1.SourcingSessionProfileJustifyResponse]) error
	SourcingSessionList(ctx context.Context, req *connect.Request[talv1.SourcingSessionListRequest]) (*talv1.SourcingSessionListResponse, error)
	SourcingSessionProfileCreateUpdate(ctx context.Context, req *connect.Request[talv1.SourcingSessionProfileCreateUpdateRequest]) (*talv1.SourcingSessionProfileCreateUpdateResponse, error)
	SourcingSessionFind(ctx context.Context, req *connect.Request[talv1.SourcingSessionFindRequest]) (*talv1.SourcingSessionFindResponse, error)
	SourcingSessionCreateUpdate(ctx context.Context, req *connect.Request[talv1.SourcingSessionCreateUpdateRequest]) (*talv1.SourcingSessionCreateUpdateResponse, error)
	ProjectInputList(ctx context.Context, req *connect.Request[talv1.ProjectInputListRequest]) (*talv1.ProjectInputListResponse, error)
}

// type SourcingSessionFilters struct {
// 	Locations []string `protobuf:"bytes,4,rep,name=locations,proto3" json:"locations,omitempty"`
// 	Skills    []string `protobuf:"bytes,5,rep,name=skills,proto3" json:"skills,omitempty"`
// 	Companies []string `protobuf:"bytes,6,rep,name=companies,proto3" json:"companies,omitempty"`
// 	Projects  []string `protobuf:"bytes,7,rep,name=projects,proto3" json:"projects,omitempty"`
// }

// type SourcingSessionFilters struct {
// 	Locations []string `protobuf:"bytes,4,rep,name=locations,proto3" json:"locations,omitempty"`
// 	Skills    []string `protobuf:"bytes,5,rep,name=skills,proto3" json:"skills,omitempty"`
// 	Companies []string `protobuf:"bytes,6,rep,name=companies,proto3" json:"companies,omitempty"`
// 	Projects  []string `protobuf:"bytes,7,rep,name=projects,proto3" json:"projects,omitempty"`
// }

// func initSourcingSessionAgent(ctx context.Context, client genaiclient.GenaiClientInterface) (llm.StructuredAgentInterface[SourcingSessionFilters], error) {
// 	agentID := "raw_profile_list_agent"
// 	persona := "You are an AI that builds RawProfileListRequest from user prompts."
// 	systemInstruction := loadTemplate("app/sourcing/agents/instruction.tmpl")
//
// 	rawAgent, err := llm.NewStructuredAgent[SourcingSessionFilters](
// 		ctx,
// 		client,
// 		agentID,
// 		persona,
// 		systemInstruction,
// 		"gemini-2.5-flash",
// 	)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to create raw profile structured agent: %w", err)
// 	}
//
// 	return rawAgent, nil
// }

// loadTemplate loads the agent template from disk
func loadTemplate(path string) string {
	content, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(content)
}

type SourcingUsecase struct {
	store   db.Store
	adapter adapter.SourcingAdapterInterface
	// llmClient        genaiclient.GenaiClientInterface
	crustDataClient crustdata.CrustdataServiceInterface
	// filtersAgent     llm.StructuredAgentInterface[SourcingSessionFilters]
	asynqClient                      asynqclient.Enqueuer // 👈 new dependency
	rawProfileFinder                 func(ctx context.Context, req *connect.Request[talv1.RawProfileFindRequest]) (*talv1.RawProfileFindResponse, error)
	sourcingSessionFiltersGuessAgent genaiclient.GenAIStructuredAgentInterface[string, agenthub.SourcingSessionFiltersInfered]
	sourcingSessionFiltersBuildAgent genaiclient.GenAIStructuredAgentInterface[struct {
		Prompt  string
		Infered agenthub.SourcingSessionFiltersInfered
	}, agenthub.SourcingSessionFilters]
	sessionProfileJustiferAgent genaiclient.GenAIStructuredAgentInterface[struct {
		SessionInfo map[string]any
		RawProfile  map[string]any
	}, string]
	agentsStore *redis.Client
	repo        repo.SourcingRepoInterface
}

func NewSourcingUsecase(
	store db.Store,
	crustDataClient crustdata.CrustdataServiceInterface,
	asynqClient asynqclient.Enqueuer, // 👈 new dependency
	rawProfileFinder func(ctx context.Context, req *connect.Request[talv1.RawProfileFindRequest]) (*talv1.RawProfileFindResponse, error),
	agentsStore *redis.Client,
	sourcingSessionFiltersGuessAgent genaiclient.GenAIStructuredAgentInterface[string, agenthub.SourcingSessionFiltersInfered],
	sourcingSessionFiltersBuildAgent genaiclient.GenAIStructuredAgentInterface[struct {
		Prompt  string
		Infered agenthub.SourcingSessionFiltersInfered
	}, agenthub.SourcingSessionFilters],
	sessionProfileJustiferAgent genaiclient.GenAIStructuredAgentInterface[struct {
		SessionInfo map[string]any
		RawProfile  map[string]any
	}, string],

) SourcingUsecaseInterface {
	return &SourcingUsecase{
		store:                            store,
		crustDataClient:                  crustDataClient,
		asynqClient:                      asynqClient,
		rawProfileFinder:                 rawProfileFinder,
		agentsStore:                      agentsStore,
		sourcingSessionFiltersGuessAgent: sourcingSessionFiltersGuessAgent,
		sourcingSessionFiltersBuildAgent: sourcingSessionFiltersBuildAgent,
		sessionProfileJustiferAgent:      sessionProfileJustiferAgent,
		adapter:                          adapter.NewSourcingAdapter(),
		repo:                             repo.NewSourcingRepo(store),
	}
}
