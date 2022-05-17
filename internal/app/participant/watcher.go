package participant

import (
	"context"
	"encoding/json"

	"github.com/parinpan/romusha/definition"
)

type Watcher func(ctx context.Context, state StateBody) error

type participator interface {
	List(ctx context.Context) List
}

func AddParticipant(participator participator) Watcher {
	var member *definition.Member

	return func(ctx context.Context, state StateBody) error {
		if state.Topic != definition.Topic_Join {
			return nil
		}

		if err := json.Unmarshal(state.Data, &member); err != nil {
			return err
		}

		return participator.List(ctx).Add(ctx, member)
	}
}

func RemoveParticipant(participator participator) Watcher {
	var member *definition.Member

	return func(ctx context.Context, state StateBody) error {
		if state.Topic != definition.Topic_Busy {
			return nil
		}

		if err := json.Unmarshal(state.Data, &member); err != nil {
			return err
		}

		return participator.List(ctx).Remove(ctx, member)
	}
}
