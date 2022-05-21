package worker

import (
	"context"

	"github.com/parinpan/romusha/definition"
	"github.com/parinpan/romusha/internal/app/scheduler"
)

type bridger interface {
	Status() definition.Status
}

type participator interface {
	Notify(ctx context.Context, state definition.StateBody) error
}

func periodicJoinMasterPing(bridger bridger, participator participator) scheduler.Exec {
	return func(ctx context.Context) error {
		if bridger.Status() != definition.Status_Occupied {
			return nil
		}

		return pingMaster(ctx, participator)
	}
}

func pingMaster(ctx context.Context, participator participator) error {
	return participator.Notify(ctx, definition.StateBody{
		Topic: definition.Topic_Join,
		Data:  definition.Member{},
	})
}
