package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type StateConfig struct {
	State string `mapstructure:"STATE"`
}

// Config stores all configuration of the application.
// The values are read by viper from a config file or environment variable.
type Config struct {
	State string `mapstructure:"STATE"`

	CrustdataAPIKey       string   `mapstructure:"CRUSTDATA_API_KEY"`
	CrustdataAPIURL       string   `mapstructure:"CRUSTDATA_API_URL"`
	GeminiAPIKey          string   `mapstructure:"GEMINI_API_KEY"`
	DefaultModel          string   `mapstructure:"DEFAULT_MODEL"`
	DefaultEmbeddingModel string   `mapstructure:"DEFAULT_EMBEDDING_MODEL"`
	AllowedOrigins        []string `mapstructure:"ALLOWED_ORIGINS"`
	ApiVersion            string   `mapstructure:"API_VERSION"`
	ApiServiceName        string   `mapstructure:"API_SERVICE_NAME"`
	DefaultBucket         string   `mapstructure:"DEFAULT_BUCKET"`

	DBProjectREF string `mapstructure:"DB_PROJECT_REF"`
	DBSource     string `mapstructure:"DB_SOURCE"`

	SupabaseTokenCookieName        string        `mapstructure:"SUPABASE_TOKEN_COOKIE_NAME"`
	SupabaseRefreshTokenCookieName string        `mapstructure:"SUPABASE_REFRESH_TOKEN_COOKIE_NAME"`
	RefreshTokenCookieName         string        `mapstructure:"REFRESH_TOKEN_COOKIE_NAME"`
	GRPCServerAddress              string        `mapstructure:"GRPC_SERVER_ADDRESS"`
	AccessTokenDuration            time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration           time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	TokenSymmetricKey              string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	ClientBaseUrl                  string        `mapstructure:"CLIENT_BASE_URL"`
	DockerHubUser                  string        `mapstructure:"DOCKERHUB_USER"`
	AppName                        string        `mapstructure:"APP_NAME"`
	GitUser                        string        `mapstructure:"GIT_USER"`

	ResendApiKey      string `mapstructure:"RESEND_API_KEY"`
	RedisAddress      string `mapstructure:"REDIS_ADDRESS"`
	GenAIRedisAddress string `mapstructure:"GENAI_REDIS_DATABASE"`

	RedisPort          string `mapstructure:"REDIS_PORT"`
	RedisHost          string `mapstructure:"REDIS_HOST"`
	RedisPassword      string `mapstructure:"REDIS_PASSWORD"`
	RedisDatabase      int    `mapstructure:"REDIS_DATABASE"`
	GenAIRedisDatabase int    `mapstructure:"GENAI_REDIS_DATABASE"`

	IsCacheDisabled        bool   `mapstructure:"IS_CACHE_DISABLED"`
	SupabaseAPIPort        uint32 `mapstructure:"SUPABASE_API_PORT"`
	SupabaseServiceRoleKey string `mapstructure:"SUPABASE_SERVICE_ROLE_KEY"`
	SupabaseApiKey         string `mapstructure:"SUPABASE_API_KEY"`

	// 👉 Typesense fields
	TypesenseHost     string `mapstructure:"TYPESENSE_HOST"`
	TypesensePort     string `mapstructure:"TYPESENSE_PORT"`
	TypesenseProtocol string `mapstructure:"TYPESENSE_PROTOCOL"`
	TypesenseApiKey   string `mapstructure:"TYPESENSE_API_KEY"`

	// Weaviate
	WeaviateHost       string `mapstructure:"WEAVIATE_HOST"`
	IsWeaviateDisabled string `mapstructure:"IS_WEAVIATE_DISABLED"`
	WeaviateScheme     string `mapstructure:"WEAVIATE_SCHEME"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadState(path string) (config StateConfig, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("state")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}

func loadAncScan(config *Config) (err error) {
	err = viper.ReadInConfig()
	if err != nil {
		return err
	}
	err = viper.Unmarshal(config)
	if err != nil {
		return err
	}
	return nil
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string, state string) (config Config, err error) {
	stateEnvFilePath := fmt.Sprintf("%s.env", state)
	viper.SetConfigName(stateEnvFilePath)

	err = loadAncScan(&config)
	if err != nil {
		return
	}
	viper.SetConfigName("shared.env")
	err = loadAncScan(&config)
	config.State = state
	if err != nil {
		return
	}
	return
}
