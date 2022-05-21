package master

import (
	"context"

	"github.com/parinpan/romusha/definition"
	"github.com/parinpan/romusha/internal/app/scheduler"
)

type participator interface {
	Notify(ctx context.Context, state definition.StateBody) error
}

func periodicCallToJoin(participator participator) scheduler.Exec {
	return func(ctx context.Context) error {
		return callToJoin(ctx, participator)
	}
}

func callToJoin(ctx context.Context, participator participator) error {
	return participator.Notify(ctx, definition.StateBody{
		Topic: definition.Topic_Call,
		Data:  nil,
	})
}
