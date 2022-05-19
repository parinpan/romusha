package job

import (
	"context"

	"github.com/parinpan/romusha/definition"
	"github.com/parinpan/romusha/internal/app/participant"
)

type queuer interface {
	Push(ctx context.Context, jobEnvelope interface{}) error
}

func RequeueJob(queue queuer) participant.Watcher {
	return func(ctx context.Context, state participant.StateBody) error {
		if state.Topic != definition.Topic_BroadcastFailure {
			return nil
		}

		return queue.Push(ctx, state.Data)
	}
}
