package participant

import (
	"context"
	"errors"

	"github.com/parinpan/romusha/definition"
)

type Watcher func(ctx context.Context, state StateBody) error

type participator interface {
	List(ctx context.Context) List
}

func AddParticipant(participator participator) Watcher {
	var member *definition.Member
	var convertible bool

	return func(ctx context.Context, state StateBody) error {
		if state.Topic != definition.Topic_Join {
			return nil
		}

		if member, convertible = state.Data.(*definition.Member); !convertible {
			return errors.New("data type is not convertible")
		}

		return participator.List(ctx).Add(ctx, member, definition.Status_Available)
	}
}

func RemoveParticipant(participator participator) Watcher {
	var member *definition.Member
	var convertible bool

	return func(ctx context.Context, state StateBody) error {
		if state.Topic != definition.Topic_Busy {
			return nil
		}

		if member, convertible = state.Data.(*definition.Member); !convertible {
			return errors.New("data type is not convertible")
		}

		return participator.List(ctx).Remove(ctx, member)
	}
}
