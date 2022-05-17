package assignor

import (
	"context"

	"github.com/parinpan/romusha/definition"
	"github.com/parinpan/romusha/internal/app/participant"
)

type participator interface {
	List(ctx context.Context) participant.List
}

type queuer interface {
	Push(ctx context.Context, jobID string, job []byte) error
}

type bridgeManager interface {
	AssignByHost(ctx context.Context, host string, envelope *definition.Envelope) (resp *definition.Response, err error)
}

type Assignor struct {
	bridger      bridgeManager
	queuer       queuer
	participator participator
}

func (a *Assignor) Assign(ctx context.Context, jobID string, callbackEndpoint string, job []byte) error {
	participants := a.participator.List(ctx)
	pickedMember := firstPick(participants)

	envelope := &definition.Envelope{
		ID:               jobID,
		CallbackEndpoint: callbackEndpoint,
		Job:              job,
	}

	if pickedMember == nil {
		return a.pushBack(ctx, envelope)
	}

	return a.assign(ctx, pickedMember, envelope)
}

func (a *Assignor) assign(ctx context.Context, member *definition.Member, envelope *definition.Envelope) error {
	response, err := a.bridger.AssignByHost(ctx, member.Host, envelope)
	if err != nil {
		return err
	}

	if response.Status != definition.BridgeStatus_Success {
		return a.pushBack(ctx, envelope)
	}

	return a.participator.List(ctx).Remove(ctx, member)
}

func (a *Assignor) pushBack(ctx context.Context, envelope *definition.Envelope) error {
	return a.queuer.Push(ctx, envelope.ID, envelope.Job)
}

func firstPick(participants participant.List) *definition.Member {
	if len(participants) == 0 {
		return nil
	}

	for host := range participants {
		return &definition.Member{
			Host: host,
		}
	}

	return nil
}
