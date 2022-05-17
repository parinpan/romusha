package job

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	statusFormatKey = "job:status:%s"

	UnknownStatus    = ""
	FailedStatus     = "failed"
	SuccessStatus    = "success"
	ProcessingStatus = "processing"
)

type statusClient interface {
	Get(ctx context.Context, key string) *redis.StringCmd
	SetEX(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
}

type Status struct {
	sc  statusClient
	ttl time.Duration
}

func NewStatus(sc statusClient, ttl time.Duration) *Status {
	return &Status{
		sc:  sc,
		ttl: ttl,
	}
}

func (s *Status) Set(ctx context.Context, jobID string, status string) error {
	return s.sc.SetEX(ctx, fmt.Sprintf(statusFormatKey, jobID), status, s.ttl).Err()
}

func (s *Status) Get(ctx context.Context, jobID string) string {
	return s.sc.Get(ctx, fmt.Sprintf(statusFormatKey, jobID)).String()
}
