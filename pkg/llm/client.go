package llm

import (
	"context"
	"fmt"

	"github.com/TALPlatform/tal_api/config"
	"github.com/darwishdev/genaiclient"
	"github.com/redis/go-redis/v9"
	"google.golang.org/genai"
)

func NewLLM(config config.Config, redisClient *redis.Client) (genaiclient.GenaiClientInterface, error) {
	ctx := context.Background()
	geminiClient, err := genai.NewClient(ctx, &genai.ClientConfig{APIKey: config.GeminiAPIKey})
	if err != nil {
		return nil, fmt.Errorf("Error create Gemini client :%w.", err)
	}
	genaiClient, err := genaiclient.NewGenaiClient(ctx, geminiClient, redisClient, config.DefaultModel, config.DefaultEmbeddingModel)
	if err != nil {
		return nil, err
	}
	return genaiClient, nil
}
