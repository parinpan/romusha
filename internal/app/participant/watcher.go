package participant

import (
	"context"
	"errors"

	"google.golang.org/grpc"

	"github.com/parinpan/romusha/definition"
)

type participator interface {
	List(ctx context.Context) List
}

type bridgeManager interface {
	Add(host string, bridger definition.BridgeClient)
}

func AddParticipant(participator participator, bridgeManager bridgeManager) definition.Watcher {
	var member *definition.Member
	var convertible bool

	var assignBridge = func(member *definition.Member) (err error) {
		if conn, err := grpc.Dial(member.Host); err == nil {
			bridgeManager.Add(member.Host, definition.NewBridgeClient(conn))
		}
		return
	}

	return func(ctx context.Context, state definition.StateBody) error {
		if state.Topic != definition.Topic_Join {
			return nil
		}

		if member, convertible = state.Data.(*definition.Member); !convertible {
			return errors.New("data type is not convertible")
		}

		if err := assignBridge(member); err != nil {
			return err
		}

		return participator.List(ctx).Add(ctx, member, definition.Status_Available)
	}
}

func RemoveParticipant(participator participator) definition.Watcher {
	var member *definition.Member
	var convertible bool

	return func(ctx context.Context, state definition.StateBody) error {
		if state.Topic != definition.Topic_Busy {
			return nil
		}

		if member, convertible = state.Data.(*definition.Member); !convertible {
			return errors.New("data type is not convertible")
		}

		return participator.List(ctx).Add(ctx, member, definition.Status_Occupied)
	}
}
