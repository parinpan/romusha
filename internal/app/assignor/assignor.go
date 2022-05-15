package assignor

import (
	"context"
	"errors"

	"github.com/parinpan/romusha/bridge"
	"github.com/parinpan/romusha/internal/app/participant"
)

var (
	ErrAssignError = errors.New("failed to assign worker")
)

type participantClient interface {
	List(ctx context.Context) participant.List
}

type queueClient interface {
	Push(ctx context.Context, job []byte) error
}

type bridgeClient interface {
	Assign(ctx context.Context, job *bridge.Job) (*bridge.Response, error)
}

type Assignor struct {
	bridgeClient bridgeClient
	queueClient  queueClient
	participant  participantClient
}

func (a *Assignor) Assign(ctx context.Context, job []byte) error {
	participants := a.participant.List(ctx)
	pickedMember := firstPick(participants)

	if pickedMember == nil {
		return a.pushBack(ctx, job)
	}

	return a.assign(ctx, job)
}

func (a *Assignor) assign(ctx context.Context, job []byte) error {
	response, err := a.bridgeClient.Assign(ctx, &bridge.Job{Gob: job})
	if err != nil {
		return err
	}

	if response.State != bridge.State_SUCCESS {
		return ErrAssignError
	}

	return nil
}

func (a *Assignor) pushBack(ctx context.Context, job []byte) error {
	return a.queueClient.Push(ctx, job)
}

func firstPick(participants participant.List) *participant.Member {
	if len(participants) == 0 {
		return nil
	}

	for host, endpoint := range participants {
		return &participant.Member{
			Host:     string(host),
			Endpoint: string(endpoint),
		}
	}

	return nil
}
