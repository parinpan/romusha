package master

import (
	"context"
	"errors"

	"github.com/parinpan/romusha/definition"
	"github.com/parinpan/romusha/internal/app/job"
)

type jobStatusWatch interface {
	Set(ctx context.Context, jobID string, status string) error
}

func jobStatusWatcher(jobStatus jobStatusWatch) definition.Watcher {
	var jobEnvelope *definition.JobEnvelope
	var convertible bool

	conversionMap := map[definition.Topic]string{
		definition.Topic_BroadcastFailure: job.FailedStatus,
		definition.Topic_BroadcastSuccess: job.SuccessStatus,
	}

	return func(ctx context.Context, state definition.StateBody) error {
		allowedTopic := state.Topic == definition.Topic_BroadcastFailure ||
			state.Topic == definition.Topic_BroadcastSuccess

		if !allowedTopic {
			return nil
		}

		if jobEnvelope, convertible = state.Data.(*definition.JobEnvelope); !convertible {
			return errors.New("data type is not convertible")
		}

		return jobStatus.Set(ctx, jobEnvelope.GetID(), conversionMap[state.Topic])
	}
}
