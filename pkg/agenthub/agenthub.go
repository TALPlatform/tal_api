package agenthub

import (
	"context"
	"fmt"

	"github.com/darwishdev/genaiclient"
	"github.com/redis/go-redis/v9"
)

type SourcingSessionFiltersInfered struct {
	Locations  bool `json:"locations"`
	Skills     bool `json:"skills"`
	Companies  bool `json:"companies"`
	Projects   bool `json:"projects"`
	YearsOfExp bool `json:"years_of_exp"`
}

// 2. Actual filter values
type SourcingSessionFilters struct {
	Locations     []string `json:"locations"`
	Skills        []string `json:"skills"`
	Companies     []string `json:"companies"`
	Projects      []string `json:"projects"`
	YearsOfExpMin int32    `json:"years_of_exp_min"`
	YearsOfExpMax int32    `json:"years_of_exp_max"`
}

type StructuredAgents struct {
	SourcingSessionFiltersGuess genaiclient.GenAIStructuredAgentInterface[string, SourcingSessionFiltersInfered]
	SourcingSessionFiltersBuild genaiclient.GenAIStructuredAgentInterface[struct {
		Prompt  string
		Infered SourcingSessionFiltersInfered
	}, SourcingSessionFilters]
	SessionProfileJustifer genaiclient.GenAIStructuredAgentInterface[struct {
		SessionInfo map[string]any
		RawProfile  map[string]any
	}, string]
}
type AgentHub struct {
	agents          map[string]genaiclient.GenAIAgentInterface
	structredAgents StructuredAgents
	rdb             *redis.Client
}

func NewAgentHub(rdb *redis.Client) *AgentHub {
	return &AgentHub{
		agents: make(map[string]genaiclient.GenAIAgentInterface),
		rdb:    rdb,
	}
}

func (h *AgentHub) RegisterAgent(name string, agent genaiclient.GenAIAgentInterface) {
	h.agents[name] = agent
}

func (h *AgentHub) GetStructuredAgents() StructuredAgents {
	return h.structredAgents
}

func (h *AgentHub) GetAgent(name string) (genaiclient.GenAIAgentInterface, bool) {
	agent, ok := h.agents[name]
	return agent, ok
}

// Example: create an in-memory session
func (h *AgentHub) NewRedisSession(ctx context.Context, agentName, userID, sessionID string) (genaiclient.GenAISessionInterface, error) {
	agent, ok := h.agents[agentName]
	if !ok {
		return nil, fmt.Errorf("agent %s not registered", agentName)
	}
	return agent.NewRedisSession(ctx, userID, sessionID, h.rdb)
}
func (h *AgentHub) RegisterStructuredAgents(apiKey string, modelName string, enableTracer bool) error {
	guessAgent, err := genaiclient.NewStructuredAgent[string, SourcingSessionFiltersInfered](
		"tal",
		apiKey,
		modelName,
		"SourcingSessionFiltersGuess",
		"Infers which filters can be inferred from a user prompt",
		"Given a user prompt, return true for only the filters that can be inferred",
		enableTracer,
	)
	if err != nil {
		return fmt.Errorf("failed to create guess agent: %w", err)
	}
	h.structredAgents.SourcingSessionFiltersGuess = guessAgent

	// 2. SourcingSessionFiltersBuild
	filtersAgent, err := genaiclient.NewStructuredAgent[struct {
		Prompt  string
		Infered SourcingSessionFiltersInfered
	}, SourcingSessionFilters](
		"tal",
		apiKey,
		modelName,
		"SourcingSessionFiltersBuild",
		"Returns the actual filter values given the prompt and inferred filters",
		"Return arrays of matching values only for the true keys in the inferred object",
		enableTracer,
	)
	if err != nil {
		return fmt.Errorf("failed to create filters agent: %w", err)
	}
	h.structredAgents.SourcingSessionFiltersBuild = filtersAgent

	// 3. AgentSourcingSessionRawProfileJustifier
	justifierAgent, err := genaiclient.NewStructuredAgent[struct {
		SessionInfo map[string]any
		RawProfile  map[string]any
	}, string](
		"tal",
		apiKey,
		"gemini-2.5-pro",
		"AgentSourcingSessionRawProfileJustifier",
		"Provides justification why a raw profile matches or not",
		"Return a short justification string based on the session info and raw profile",
		enableTracer,
	)
	if err != nil {
		return fmt.Errorf("failed to create justifier agent: %w", err)
	}
	h.structredAgents.SessionProfileJustifer = justifierAgent

	return nil
}
