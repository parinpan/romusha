package participant

import (
	"context"
	"encoding/json"
)

type Watcher func(ctx context.Context, state StateBody) error

type participator interface {
	List(ctx context.Context) List
}

func AddParticipant(participator participator) Watcher {
	var member Member

	return func(ctx context.Context, state StateBody) error {
		if state.Topic != Join {
			return nil
		}

		if err := json.Unmarshal(state.Data, &member); err != nil {
			return err
		}

		return participator.List(ctx).Add(ctx, member)
	}
}

func RemoveParticipant(participator participator) Watcher {
	var member Member

	return func(ctx context.Context, state StateBody) error {
		if state.Topic != Busy {
			return nil
		}

		if err := json.Unmarshal(state.Data, &member); err != nil {
			return err
		}

		return participator.List(ctx).Remove(ctx, member)
	}
}
