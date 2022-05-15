package job

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	jobKeyFormat = "job:queue:%s"
)

type QueueProcessor func(ctx context.Context, job []byte) error

type queueClient interface {
	LPop(ctx context.Context, key string) *redis.StringCmd
	RPush(ctx context.Context, key string, values ...interface{}) *redis.IntCmd
}

type Queue struct {
	key string
	qc  queueClient
}

func NewQueue(key string, qc queueClient) *Queue {
	return &Queue{
		key: jobKey(key),
		qc:  qc,
	}
}

func (q *Queue) Push(ctx context.Context, goBytes []byte) error {
	return q.qc.RPush(ctx, q.key, goBytes).Err()
}

func (q *Queue) Poll(ctx context.Context, timeout time.Duration, processor QueueProcessor) error {
	for {
		b, err := q.qc.LPop(ctx, q.key).Bytes()
		if err != nil && err != redis.Nil {
			return err
		}

		if err := processor(ctx, b); err != nil {
			return err
		}

		time.Sleep(timeout)
	}
}

func jobKey(queueKey string) string {
	return fmt.Sprintf(jobKeyFormat, queueKey)
}
