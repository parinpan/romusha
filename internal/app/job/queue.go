package job

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	queueKey = "job:queue"
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

func NewQueue(qc queueClient) *Queue {
	return &Queue{
		key: queueKey,
		qc:  qc,
	}
}

func (q *Queue) Push(ctx context.Context, jobEnvelope []byte) error {
	return q.qc.RPush(ctx, q.key, jobEnvelope).Err()
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
