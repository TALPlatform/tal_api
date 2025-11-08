package db

import (
	"context"
	"os"
	"testing"

	"github.com/TALPlatform/tal_api/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

var store Store
var connPool *pgxpool.Pool

func intInSlice(target int32, slice []int32) bool {
	for _, item := range slice {
		if item == target {
			return true
		}
	}
	return false
}

func TestMain(m *testing.M) {
	ctx := context.Background()
	_, err := config.LoadState("../config")
	if err != nil {
		panic(err)
	}
	config, err := config.LoadConfig("../config", "test")

	if err != nil {
		log.Fatal().Err(err).Msg("failed load config")
	}
	store, _, err = InitDB(ctx, config.DBSource, false)

	if err != nil {
		log.Fatal().Str("DBSource", config.DBSource).Err(err).Msg("db failed to connect")
	}

	os.Exit(m.Run())
}
