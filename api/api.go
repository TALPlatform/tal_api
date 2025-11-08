package api

import (
	// USECASE_IMPORTS
	sourcingUsecase "github.com/TALPlatform/tal_api/app/sourcing/usecase"

	"fmt"

	peopleUsecase "github.com/TALPlatform/tal_api/app/people/usecase"

	accountsUsecase "github.com/TALPlatform/tal_api/app/accounts/usecase"
	publicUsecase "github.com/TALPlatform/tal_api/app/public/usecase"
	tenantUsecase "github.com/TALPlatform/tal_api/app/tenant/usecase"
	"github.com/TALPlatform/tal_api/config"
	"github.com/TALPlatform/tal_api/db"
	"github.com/TALPlatform/tal_api/pkg/auth"
	"github.com/TALPlatform/tal_api/pkg/crustdata"
	"github.com/TALPlatform/tal_api/pkg/llm"
	"github.com/TALPlatform/tal_api/pkg/redisclient"
	"github.com/bufbuild/protovalidate-go"

	"github.com/TALPlatform/tal_api/pkg/resend"
	weaviateclient "github.com/TALPlatform/tal_api/pkg/weaviateclient"
	"github.com/TALPlatform/tal_api/proto_gen/tal/v1/talv1connect"
	"github.com/darwishdev/sqlseeder"
	supaapigo "github.com/darwishdev/supaapi-go"
	"golang.org/x/crypto/bcrypt"
)

type Api struct {
	talv1connect.UnimplementedTalServiceHandler
	accountsUsecase accountsUsecase.AccountsUsecaseInterface
	config          config.Config
	validator       *protovalidate.Validator
	tokenMaker      auth.Maker
	sqlSeeder       sqlseeder.SeederInterface
	publicUsecase   publicUsecase.PublicUsecaseInterface

	weaviateClient weaviateclient.WeaviateClientInterface // 👈 NEW FIELD
	tenantUsecase  tenantUsecase.TenantUsecaseInterface
	// USECASE_FIELDS
	sourcingUsecase sourcingUsecase.SourcingUsecaseInterface

	peopleUsecase peopleUsecase.PeopleUsecaseInterface

	supaapi     supaapigo.Supaapi
	redisClient redisclient.RedisClientInterface
	store       db.Store
}

func HashFunc(req string) string {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req), bcrypt.DefaultCost)
	return string(hashedPassword)
}
func NewApi(config config.Config, store db.Store, tokenMaker auth.Maker, redisClient redisclient.RedisClientInterface, validator *protovalidate.Validator) (talv1connect.TalServiceHandler, error) {
	resendClient, err := resend.NewResendService(config.ResendApiKey, config.ClientBaseUrl)
	if err != nil {
		return nil, err
	}

	// typesenseClient := typesenseclient.NewTypesenseClient(config.TypesenseHost, config.TypesensePort, config.TypesenseProtocol , config.TypesenseApiKey , false)
	//
	// err = typesenseClient.CreateCommandPaletteCollectionIfNotExists(context.Background())
	// if err != nil {
	// 	return nil , err
	// }
	// 👇 INIT WEAVIATE CLIENT
	weaviateClient, err := weaviateclient.NewWeaviateClient(
		config.WeaviateHost,
		config.WeaviateScheme,
	)
	if err != nil {
		return nil, err
	}
	supEnv := supaapigo.DEV

	if config.State == "prod" {
		supEnv = supaapigo.PROD
	}
	supaapi := supaapigo.NewSupaapi(supaapigo.SupaapiConfig{
		ProjectRef:     config.DBProjectREF,
		Env:            supEnv,
		Port:           config.SupabaseAPIPort,
		ServiceRoleKey: config.SupabaseServiceRoleKey,
		ApiKey:         config.SupabaseApiKey,
	})
	llmClient, err := llm.NewLLM(config, redisClient.GetClient())
	if err != nil {
		return nil, fmt.Errorf("error creating llm api client : %w", err)
	}

	crustDataClient, err := crustdata.NewCrustdataService(config.CrustdataAPIKey, config.CrustdataAPIURL)
	if err != nil {
		return nil, fmt.Errorf("error creating crustdata api client : %w", err)
	}

	sqlSeeder := sqlseeder.NewSeeder(sqlseeder.SeederConfig{HashFunc: HashFunc})
	// USECASE_INSTANTIATIONS
	sourcingUsecase := sourcingUsecase.NewSourcingUsecase(store, llmClient)

	peopleUsecase := peopleUsecase.NewPeopleUsecase(store, crustDataClient, llmClient)

	tenantUsecase := tenantUsecase.NewTenantUsecase(store, redisClient)
	accountsUsecase := accountsUsecase.NewAccountsUsecase(store, supaapi, redisClient, tokenMaker, config.AccessTokenDuration, config.RefreshTokenDuration)
	publicUsecase := publicUsecase.NewPublicUsecase(store, config.SupabaseApiKey, supaapi, redisClient, resendClient, weaviateClient)
	return &Api{
		// USECASE_INJECTIONS
		sourcingUsecase: sourcingUsecase,

		peopleUsecase: peopleUsecase,

		weaviateClient: weaviateClient,
		// typesenseClient: typesenseClient,
		accountsUsecase: accountsUsecase,
		tenantUsecase:   tenantUsecase,
		store:           store,
		redisClient:     redisClient,
		tokenMaker:      tokenMaker,
		supaapi:         supaapi,
		config:          config,
		sqlSeeder:       sqlSeeder,
		publicUsecase:   publicUsecase,
		validator:       validator,
	}, nil
}
