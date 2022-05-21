package worker

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/parinpan/romusha/definition"
	"github.com/parinpan/romusha/internal/app/bridge"
	"github.com/parinpan/romusha/internal/app/participant"
	"github.com/parinpan/romusha/internal/app/scheduler"
)

func Start(ctx context.Context) error {
	rc := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	pcp := participant.NewParticipant(rc)
	jt := NewJobTracker(pcp, &definition.Member{})
	bridgeServer := bridge.NewServer(jt)

	errChan := make(chan error, 1)
	sigCapt := make(chan os.Signal, 1)
	signal.Notify(sigCapt, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGKILL)

	go func() {
		errChan <- startBridgerRpc(ctx, bridgeServer)
	}()

	go func() {
		scheduler.Schedule(
			ctx,
			time.Duration(20)*time.Millisecond,
			periodicJoinMasterPing(bridgeServer, pcp))
	}()

	// ping master to ask to participate
	_ = pingMaster(ctx, pcp)

	select {
	case err := <-errChan:
		log.Println("exited with an error: ", err.Error())
		return err
	case sig := <-sigCapt:
		log.Println("exited with a signal: ", sig.String())
	}

	return nil
}
