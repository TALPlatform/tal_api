package asynqworker

import (
	"context"
	"log"

	"github.com/hibiken/asynq"
)

type HandlerFunc func(context.Context, *asynq.Task) error

type Worker struct {
	srv *asynq.Server
	mux *asynq.ServeMux
}

func New(redisAddr string) *Worker {
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr},
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"default": 10,
			},
		},
	)

	mux := asynq.NewServeMux()
	return &Worker{srv: srv, mux: mux}
}

func (w *Worker) Register(taskType string, handler HandlerFunc) {
	w.mux.HandleFunc(taskType, handler)
}

func (w *Worker) Run() error {
	log.Println("[Asynq] Worker started...")
	return w.srv.Run(w.mux)
}
func (w *Worker) Stop() error {
	log.Println("[Asynq] Worker stopping gracefully...")
	w.srv.Shutdown()
	return nil
}
