package bridge

import (
	"context"

	"github.com/golang/protobuf/proto"

	"github.com/parinpan/romusha/definition"
	"github.com/parinpan/romusha/internal/app/job"
)

type jobTracker interface {
	Track(ctx context.Context, envelope *definition.JobEnvelope, exec definition.Executor, status *definition.Status)
}

// Server is a grpc client for all workers only
type Server struct {
	definition.UnimplementedBridgeServer
	tracker jobTracker
	status  definition.Status
}

func NewServer(jobTracker jobTracker) *Server {
	return &Server{
		tracker: jobTracker,
		status:  definition.Status_Available,
	}
}

func (s *Server) Assign(ctx context.Context, envelope *definition.JobEnvelope) (resp *definition.Response, err error) {
	resp.Message = *proto.String("machine is available")
	resp.Status = definition.BridgeStatus_Success

	if s.status == definition.Status_Occupied {
		resp.Message = *proto.String("machine is occupied")
		resp.Status = definition.BridgeStatus_Occupied
		return
	}

	// set worker as occupied because it's about to use
	s.status = definition.Status_Occupied

	var exec definition.Executor

	if err = job.Decode(envelope.Executor, &exec); err != nil {
		s.status = definition.Status_Available
		resp.Message = *proto.String("machine can't decode job processor")
		resp.Status = definition.BridgeStatus_Error
		return
	}

	go s.tracker.Track(ctx, envelope, exec, &s.status)

	return
}

func (s *Server) Status() definition.Status {
	return s.status
}
