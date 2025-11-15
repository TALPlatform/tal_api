package asynqclient

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/hibiken/asynq"
)

type Enqueuer interface {
	Enqueue(ctx context.Context, taskType string, payload interface{}) (*asynq.TaskInfo, error)
	Stop() error
}
type Client struct {
	client *asynq.Client
}

func New(redisAddr string) Enqueuer {
	return &Client{
		client: asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr}),
	}
}

func (c *Client) Close() error {
	return c.client.Close()
}

func (c *Client) Enqueue(ctx context.Context, taskType string, payload interface{}) (*asynq.TaskInfo, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	task := asynq.NewTask(taskType, data)
	return c.client.EnqueueContext(ctx, task, asynq.MaxRetry(3), asynq.Timeout(5*time.Minute))
}
func (c *Client) Stop() error {
	log.Println("[Asynq] Client stopping gracefully...")
	return c.client.Close()
}
