package job

import (
	"context"

	"github.com/parinpan/romusha/internal/app/participant"
)

type queuer interface {
	Push(ctx context.Context, goBytes []byte) error
}

func RequeueJob(queue queuer) participant.Watcher {
	return func(ctx context.Context, state participant.StateBody) error {
		if state.Topic != participant.Fault {
			return nil
		}

		return queue.Push(ctx, state.Data)
	}
}
