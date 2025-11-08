package usecase

import (
	"context"
	"fmt"
	"os"

	"connectrpc.com/connect"
	"github.com/TALPlatform/tal_api/app/people/adapter"
	"github.com/TALPlatform/tal_api/app/people/repo"
	"github.com/TALPlatform/tal_api/db"
	"github.com/TALPlatform/tal_api/pkg/crustdata"
	"github.com/TALPlatform/tal_api/pkg/llm"
	talv1 "github.com/TALPlatform/tal_api/proto_gen/tal/v1"
	"github.com/darwishdev/genaiclient"
)

type PeopleUsecaseInterface interface {
	// INJECT INTERFACE

	RawProfileListRequestBuild(ctx context.Context, req *connect.Request[talv1.RawProfileListRequestBuildRequest]) (*talv1.RawProfileListRequestBuildResponse, error)

	RawProfileFind(ctx context.Context, req *connect.Request[talv1.RawProfileFindRequest]) (*talv1.RawProfileFindResponse, error)

	RawProfileList(ctx context.Context, req *connect.Request[talv1.RawProfileListRequest]) (*talv1.RawProfileListResponse, error)
}

type PeopleUsecase struct {
	store           db.Store
	adapter         adapter.PeopleAdapterInterface
	llmClient       genaiclient.GenaiClientInterface
	rawProfileAgent llm.StructuredAgentInterface[RawProfileListRequestRaw] // new
	crustDataClient crustdata.CrustdataServiceInterface
	repo            repo.PeopleRepoInterface
}

type RawProfileListRequestRaw struct {
	Industries []string `protobuf:"bytes,3,rep,name=industries,proto3" json:"industries,omitempty"`
	Locations  []string `protobuf:"bytes,4,rep,name=locations,proto3" json:"locations,omitempty"`
	Skills     []string `protobuf:"bytes,5,rep,name=skills,proto3" json:"skills,omitempty"`
	Companies  []string `protobuf:"bytes,6,rep,name=companies,proto3" json:"companies,omitempty"`
	Projects   []string `protobuf:"bytes,7,rep,name=projects,proto3" json:"projects,omitempty"`
}

func initRawProfileAgent(ctx context.Context, client genaiclient.GenaiClientInterface) (llm.StructuredAgentInterface[RawProfileListRequestRaw], error) {
	agentID := "raw_profile_list_agent"
	persona := "You are an AI that builds RawProfileListRequest from user prompts."
	systemInstruction := loadTemplate("app/people/agents/raw_profile_list.tmpl")

	rawAgent, err := llm.NewStructuredAgent[RawProfileListRequestRaw](
		ctx,
		client,
		agentID,
		persona,
		systemInstruction,
		"gemini-2.5-flash", // or your default model
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create raw profile structured agent: %w", err)
	}

	return rawAgent, nil
}

// loadTemplate loads the agent template from disk
func loadTemplate(path string) string {
	content, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(content)
}
func NewPeopleUsecase(store db.Store, crustDataClient crustdata.CrustdataServiceInterface, llmClient genaiclient.GenaiClientInterface) (PeopleUsecaseInterface, error) {
	ctx := context.Background()
	rawProfileAgent, err := initRawProfileAgent(
		ctx,
		llmClient,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create structured agent for raw profile requests: %w", err)
	}
	return &PeopleUsecase{
		store:           store,
		crustDataClient: crustDataClient,
		llmClient:       llmClient,
		rawProfileAgent: rawProfileAgent,
		adapter:         adapter.NewPeopleAdapter(llmClient.Embed, llmClient.EmbedBulk),
		repo:            repo.NewPeopleRepo(store),
	}, nil
}
