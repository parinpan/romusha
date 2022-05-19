package assignor

import (
	"context"

	"github.com/parinpan/romusha/definition"
	"github.com/parinpan/romusha/internal/app/participant"
)

type participator interface {
	List(ctx context.Context) participant.List
	Notify(ctx context.Context, state participant.StateBody) error
}

type jobQueuer interface {
	Push(ctx context.Context, jobEnvelope interface{}) error
}

type bridgeManager interface {
	AssignByHost(ctx context.Context, host string, envelope *definition.JobEnvelope) (resp *definition.Response, err error)
}

type Assignor struct {
	bridger      bridgeManager
	jobQueuer    jobQueuer
	participator participator
}

func (a *Assignor) Assign(ctx context.Context, jobID string, source, callbackUrl string, executor []byte) error {
	participants := a.participator.List(ctx)
	pickedMember := firstPick(participants)

	envelope := &definition.JobEnvelope{
		ID: jobID,
		Request: &definition.JobRequest{
			Source:      source,
			CallbackUrl: callbackUrl,
		},
		Executor: executor,
	}

	if pickedMember == nil {
		return a.pushBack(ctx, envelope)
	}

	return a.assign(ctx, pickedMember, envelope)
}

func (a *Assignor) assign(ctx context.Context, member *definition.Member, envelope *definition.JobEnvelope) error {
	defer a.participator.List(ctx).Add(ctx, member, definition.Status_Occupied)

	response, err := a.bridger.AssignByHost(ctx, member.Host, envelope)
	if err != nil {
		return err
	}

	if response.Status != definition.BridgeStatus_Success {
		return a.pushBack(ctx, envelope)
	}

	return nil
}

func (a *Assignor) pushBack(ctx context.Context, envelope *definition.JobEnvelope) error {
	return a.jobQueuer.Push(ctx, envelope)
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
