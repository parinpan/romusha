package worker

import (
	"context"
	"log"

	"github.com/parinpan/romusha/definition"
	"github.com/parinpan/romusha/internal/app/participant"
)

type participantNotifier interface {
	Notify(ctx context.Context, state participant.StateBody) error
}

type JobTracker struct {
	member *definition.Member
	pn     participantNotifier
}

func NewJobTracker(pn participantNotifier, member *definition.Member) *JobTracker {
	return &JobTracker{
		member: member,
		pn:     pn,
	}
}

func (t *JobTracker) Track(
	ctx context.Context,
	envelope *definition.JobEnvelope,
	executor definition.Executor, status *definition.Status) {

	t.notify(ctx, definition.Topic_Busy, t.member)

	defer func() {
		*status = definition.Status_Available
		t.notify(ctx, definition.Topic_Join, t.member)
	}()

	if err := executor(ctx, envelope.Request.Source); err != nil {
		t.notify(ctx, definition.Topic_BroadcastFailure, envelope)
		return
	}

	t.notify(ctx, definition.Topic_BroadcastSuccess, envelope)
}

func (t *JobTracker) notify(ctx context.Context, topic definition.Topic, data interface{}) {
	go func() {
		if err := t.pn.Notify(ctx, participant.StateBody{
			Topic: topic,
			Data:  data,
		}); err != nil {
			log.Printf("notify %s with an error: %#v", topic.String(), err)
		}
	}()
}
