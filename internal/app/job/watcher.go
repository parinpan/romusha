package job

import (
	"context"
	"errors"

	"github.com/parinpan/romusha/definition"
)

type queuer interface {
	Push(ctx context.Context, jobEnvelope interface{}) error
}

type statusWatcher interface {
	Get(ctx context.Context, jobID string) string
	Set(ctx context.Context, jobID string, status string) error
}

func RequeueJob(queue queuer, sw statusWatcher) definition.Watcher {
	var jobEnvelope *definition.JobEnvelope
	var convertible bool

	return func(ctx context.Context, state definition.StateBody) error {
		if state.Topic != definition.Topic_BroadcastFailure {
			return nil
		}

		if jobEnvelope, convertible = state.Data.(*definition.JobEnvelope); !convertible {
			return errors.New("data type is not convertible")
		}

		if err := sw.Set(ctx, jobEnvelope.GetID(), FailedStatus); err != nil {
			return err
		}

		return queue.Push(ctx, state.Data)
	}
}
