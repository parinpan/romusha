package worker

import (
	"context"
	"encoding/json"
	"log"

	"github.com/parinpan/romusha/definition"
	"github.com/parinpan/romusha/internal/app/job"
	"github.com/parinpan/romusha/internal/app/participant"
)

type participator interface {
	Notify(ctx context.Context, state participant.StateBody) error
}

type JobTracker struct {
	member      []byte
	participant participator
}

func NewJobTracker(participant participator, member participant.Member) (*JobTracker, error) {
	b, err := json.Marshal(member)
	if err != nil {
		return nil, err
	}

	return &JobTracker{
		member:      b,
		participant: participant,
	}, nil
}

func (t *JobTracker) Track(
	ctx context.Context,
	envelope *definition.Envelope,
	status *participant.Status, job *definition.Job, processor job.Processor) {

	t.notifyBusy(ctx)

	defer func() {
		*status = participant.Available
		t.notifyJoin(ctx)
	}()

	if err := processor(ctx, job.Sources); err != nil {
		t.notifyJob(ctx, definition.Topic_Fault, envelope)
		return
	}

	t.notifyJob(ctx, definition.Topic_Fault, envelope)
}

func (t *JobTracker) notifyBusy(ctx context.Context) {
	go func() {
		sb := participant.StateBody{
			Topic: participant.Busy,
			Data:  []byte(nil),
		}

		if err := t.participant.Notify(ctx, sb); err != nil {
			log.Printf("notify-busy with error: %#v", err)
		}
	}()
}

func (t *JobTracker) notifyJoin(ctx context.Context) {
	go func() {
		sb := participant.StateBody{
			Topic: participant.Join,
			Data:  t.member,
		}

		if err := t.participant.Notify(ctx, sb); err != nil {
			log.Printf("notify-join with error: %#v", err)
		}
	}()
}

func (t *JobTracker) notifyJob(ctx context.Context, topic definition.Topic, envelope *definition.Envelope) {
	go func() {
		b, err := job.Encode[definition.Envelope](envelope)
		if err != nil {
			return
		}

		sb := participant.StateBody{
			Topic: participant.Topic(topic.String()),
			Data:  b,
		}

		if err := t.participant.Notify(ctx, sb); err != nil {
			log.Printf("notify-fault with error: %#v", err)
		}
	}()
}
