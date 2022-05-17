package worker

import (
	"context"
	"encoding/json"
	"log"

	"github.com/parinpan/romusha/definition"
	"github.com/parinpan/romusha/internal/app/participant"
)

type participator interface {
	Notify(ctx context.Context, state participant.StateBody) error
}

type JobTracker struct {
	member      *definition.Member
	participant participator
}

func NewJobTracker(participant participator, member *definition.Member) *JobTracker {
	return &JobTracker{
		member:      member,
		participant: participant,
	}
}

func (t *JobTracker) Track(
	ctx context.Context,
	envelope *definition.JobEnvelope,
	executor definition.Executor, status *definition.Status) {

	notify[*definition.Member](ctx, t.participant, definition.Topic_Busy, t.member)

	defer func() {
		*status = definition.Status_Available
		notify[*definition.Member](ctx, t.participant, definition.Topic_Join, t.member)
	}()

	if err := executor(ctx, envelope.Request.Source); err != nil {
		notify[*definition.JobEnvelope](ctx, t.participant, definition.Topic_BroadcastFailure, envelope)
		return
	}

	notify[*definition.JobEnvelope](ctx, t.participant, definition.Topic_BroadcastSuccess, envelope)
}

func notify[T](ctx context.Context, participator participator, topic definition.Topic, data T) {
	go func() {
		b, err := json.Marshal(data)
		if err != nil {
			log.Printf("notify %s with error: %#v", topic, err)
		}

		stateBody := participant.StateBody{
			Topic: topic,
			Data:  b,
		}

		if err := participator.Notify(ctx, stateBody); err != nil {
			log.Printf("notify %s with error: %#v", topic, err)
		}
	}()
}
