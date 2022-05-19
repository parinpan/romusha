package job

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/parinpan/romusha/definition"
)

const (
	queueKey = "job:queue"

	nJobCount = 100
)

type QueueProcessor func(ctx context.Context, envelope *definition.JobEnvelope) error

type queueClient interface {
	LPopCount(ctx context.Context, key string, count int) *redis.StringSliceCmd
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

func (q *Queue) Push(ctx context.Context, jobEnvelope interface{}) error {
	return q.qc.RPush(ctx, q.key, jobEnvelope).Err()
}

func (q *Queue) Poll(ctx context.Context, timeout time.Duration, processor QueueProcessor) error {
	var jobEnvelopes []*definition.JobEnvelope

	for {
		err := q.qc.LPopCount(ctx, q.key, nJobCount).ScanSlice(&jobEnvelopes)
		if err != nil && err != redis.Nil {
			return err
		}

		for _, envelope := range jobEnvelopes {
			if err := processor(ctx, envelope); err != nil {
				return err
			}
		}

		time.Sleep(timeout)
	}
}
