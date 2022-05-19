package job

import (
	"context"

	"github.com/parinpan/romusha/definition"
)

type assignor interface {
	Assign(ctx context.Context, jobID string, source, callbackUrl string, executor []byte) error
}

func AssignJob(assignor assignor) QueueProcessor {
	return func(ctx context.Context, envelope *definition.JobEnvelope) error {
		return assignor.Assign(
			ctx,
			envelope.GetID(),
			envelope.GetRequest().GetSource(),
			envelope.GetRequest().GetCallbackUrl(),
			envelope.GetExecutor())
	}
}
