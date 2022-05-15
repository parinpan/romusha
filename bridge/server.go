package bridge

import (
	"context"

	"github.com/golang/protobuf/proto"

	"github.com/parinpan/romusha/internal/app/job"
	"github.com/parinpan/romusha/internal/app/participant"
)

// Server is a grpc client for all workers only
type Server struct {
	tracker *job.Track
	status  participant.Status
}

func (s *Server) Assign(ctx context.Context, j *Job) (resp *Response, err error) {
	message := proto.String("machine is available")
	state := State_SUCCESS

	if s.status == participant.Occupied {
		message = proto.String("machine is occupied")
		state = State_OCCUPIED
	}

	unpackedJob := job.Job{}
	s.status = participant.Occupied

	if err = job.Decode(j.Gob, &unpackedJob); err != nil {
		s.status = participant.Available
		message = proto.String("machine can't decode job")
		state = State_ERROR
	}

	go s.tracker.Track(ctx, &s.status, unpackedJob)

	return &Response{
		Message: *message,
		State:   state,
	}, err
}
