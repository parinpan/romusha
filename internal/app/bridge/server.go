package bridge

import (
	"context"

	"github.com/golang/protobuf/proto"

	"github.com/parinpan/romusha/definition"
	"github.com/parinpan/romusha/internal/app/job"
	"github.com/parinpan/romusha/internal/app/worker"
)

// Server is a grpc client for all workers only
type Server struct {
	tracker *worker.JobTracker
	status  definition.Status
}

func (s *Server) Assign(ctx context.Context, envelope *definition.Envelope) (resp *definition.Response, err error) {
	resp.Message = *proto.String("machine is available")
	resp.Status = definition.BridgeStatus_Success

	if s.status == definition.Status_Occupied {
		resp.Message = *proto.String("machine is occupied")
		resp.Status = definition.BridgeStatus_Occupied
		return
	}

	// set worker as occupied because it's about to use
	s.status = definition.Status_Occupied

	var dispatchedJob definition.Job
	var jobProcessor job.Processor

	if err = job.Decode[definition.Job](envelope.Job, &dispatchedJob); err != nil {
		s.status = definition.Status_Available
		resp.Message = *proto.String("machine can't decode job")
		resp.Status = definition.BridgeStatus_Error
		return
	}

	if err = job.Decode[job.Processor](dispatchedJob.Processor, &jobProcessor); err != nil {
		s.status = definition.Status_Available
		resp.Message = *proto.String("machine can't decode job processor")
		resp.Status = definition.BridgeStatus_Error
		return
	}

	go s.tracker.Track(ctx, envelope, &s.status, &dispatchedJob, jobProcessor)

	return
}
