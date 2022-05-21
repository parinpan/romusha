package worker

import (
	"context"
	"net"

	"google.golang.org/grpc"

	"github.com/parinpan/romusha/definition"
	"github.com/parinpan/romusha/internal/app/bridge"
)

func startBridgerRpc(_ context.Context, bridger *bridge.Server) error {
	l, err := net.Listen("tcp", "localhost:9000")
	if err != nil {
		return err
	}

	srv := grpc.NewServer()
	definition.RegisterBridgeServer(srv, bridger)

	return srv.Serve(l)
}
