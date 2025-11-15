package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/TALPlatform/tal_api/api"
	"github.com/TALPlatform/tal_api/config"
	"github.com/TALPlatform/tal_api/db"
	"github.com/TALPlatform/tal_api/pkg/asynqclient"
	"github.com/TALPlatform/tal_api/pkg/asynqworker"
	"github.com/TALPlatform/tal_api/pkg/auth"
	"github.com/TALPlatform/tal_api/pkg/redisclient"
	"github.com/bufbuild/protovalidate-go"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// operation is a clean up function on shutting down
type operation func(ctx context.Context) error

// gracefulShutdown waits for termination syscalls and doing clean up operations after received it
func gracefulShutdown(ctx context.Context, timeout time.Duration, ops map[string]operation) <-chan struct{} {
	wait := make(chan struct{})
	go func() {
		s := make(chan os.Signal, 1)

		// add any other syscalls that you want to be notified with
		signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
		<-s

		log.Info().Msg("shutting down")

		// set timeout for the ops to be done to prevent system hang
		timeoutFunc := time.AfterFunc(timeout, func() {
			log.Printf("timeout %d ms has been elapsed, force exit", timeout.Milliseconds())
			os.Exit(0)
		})

		defer timeoutFunc.Stop()

		var wg sync.WaitGroup

		// Do the operations asynchronously to save time
		for key, op := range ops {
			wg.Add(1)
			innerOp := op
			innerKey := key
			go func() {
				defer wg.Done()

				log.Printf("cleaning up: %s", innerKey)
				if err := innerOp(ctx); err != nil {
					log.Printf("%s: clean up failed: %s", innerKey, err.Error())
					return
				}

				log.Printf("%s was shutdown gracefully", innerKey)
			}()
		}

		wg.Wait()

		close(wait)
	}()

	return wait
}
func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	ctx := context.Background()

	state, err := config.LoadState("./config")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load the state config")
	}
	config, err := config.LoadConfig("./config", state.State)

	store, connPool, err := db.InitDB(ctx, config.DBSource, config.State == "dev")
	if err != nil {
		log.Fatal().Str("DBSource", config.DBSource).Err(err).Msg("db failed to connect")
	}
	tokenMaker, err := auth.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		panic("cann't create paset maker in gapi/api.go")
	}

	redisClient := redisclient.NewRedisClient(config.RedisHost, config.RedisPort, config.RedisPassword, config.RedisDatabase, config.IsCacheDisabled)
	validator, err := protovalidate.New()
	if err != nil {
		log.Fatal().Err(err).Msg("can't get the validator")
	}

	genAiRedisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.RedisHost, config.RedisPort),
		Password: config.RedisPassword,      // no password set
		DB:       config.GenAIRedisDatabase, // use default DB
	})

	asynqAddr := fmt.Sprintf("%s:%s", config.RedisHost, config.RedisPort)
	asynqClient := asynqclient.New(asynqAddr)
	asynqWorker := asynqworker.New(asynqAddr)
	go func() {
		if err := asynqWorker.Run(); err != nil {
			log.Fatal().Err(err).Msg("asynq worker failed")
		}
	}()
	server, err := api.NewServer(config, store, tokenMaker, redisClient, genAiRedisClient, validator, asynqClient, asynqWorker) // Start the server in a goroutine
	if err != nil {
		log.Fatal().Err(err).Msg("server initialization failed")
	}
	httpServer := server.NewGrpcHttpServer()
	go func() {
		log.Info().Str("server address", config.GRPCServerAddress).Msg("GRPC server start")
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("HTTP listen and serve failed")
		}
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	wait := gracefulShutdown(ctx, 3*time.Second, map[string]operation{
		"database": func(ctx context.Context) error {
			connPool.Close()
			return nil
		},
		"http-server": func(ctx context.Context) error {
			return httpServer.Shutdown(ctx)

		},
		"asynq-client": func(ctx context.Context) error {
			return asynqClient.Stop()
		},
		"asynq-worker": func(ctx context.Context) error {
			asynqWorker.Stop() // Implement Stop in your worker (srv.Shutdown())
			return nil
		},
		// Add other cleanup operations here
	})
	<-wait
}
