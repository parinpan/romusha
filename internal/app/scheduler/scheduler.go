package scheduler

import (
	"context"
	"log"
	"time"
)

type Exec func(ctx context.Context) error

func Schedule(ctx context.Context, doEvery time.Duration, exec Exec) {
	ticker := time.NewTicker(doEvery)

	for range ticker.C {
		if err := exec(ctx); err != nil {
			log.Printf("scheduler with an error: %#v", err)
		}
	}
}
