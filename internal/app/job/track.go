package job

import (
	"context"
	"log"

	"github.com/parinpan/romusha/internal/app/participant"
)

type participantClient interface {
	Notify(ctx context.Context, state participant.StateBody) error
}

type Track struct {
	participant participantClient
}

func (t *Track) Track(ctx context.Context, status *participant.Status, job Job) {
	defer func() {
		*status = participant.Available
	}()

	if err := job.Processor(ctx, job.FilePaths); err == nil {
		return
	} else {
		log.Printf("processed job with error: %#v", err)
	}

	b, err := Encode(job)
	log.Printf("job encoded with error: %#v", err)

	sb := participant.StateBody{
		Topic: participant.Fault,
		Data:  b,
	}

	if err := t.participant.Notify(ctx, sb); err != nil {
		log.Printf("notify with error: %#v", err)
	}
}
