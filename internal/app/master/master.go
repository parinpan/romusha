package master

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/parinpan/romusha/internal/app/assignor"
	"github.com/parinpan/romusha/internal/app/bridge"
	"github.com/parinpan/romusha/internal/app/job"
	"github.com/parinpan/romusha/internal/app/participant"
	"github.com/parinpan/romusha/internal/app/scheduler"
)

func Start(ctx context.Context) error {
	jobStatusTTL := time.Duration(24) * time.Hour
	jobPollDelay := time.Duration(100) * time.Millisecond

	bm := &bridge.Manager{}
	rc := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	pcp := participant.NewParticipant(rc)

	jobQ := job.NewQueue(rc)
	jobStatus := job.NewStatus(rc, jobStatusTTL)
	asg := assignor.NewAssignor(bm, jobQ, jobStatus, pcp)

	errChan := make(chan error, 1)
	sigCapt := make(chan os.Signal, 1)
	signal.Notify(sigCapt, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGKILL)

	go func() {
		errChan <- pcp.Watch(
			ctx,
			participant.AddParticipant(pcp, bm),
			participant.RemoveParticipant(pcp),
			job.RequeueJob(jobQ, jobStatus),
			jobStatusWatcher(jobStatus))
	}()

	go func() {
		errChan <- jobQ.Poll(ctx, jobPollDelay, job.AssignJob(asg))
	}()

	go func() {
		scheduler.Schedule(
			ctx,
			time.Duration(500)*time.Millisecond,
			periodicCallToJoin(pcp))
	}()

	// ask all available workers to join
	_ = callToJoin(ctx, pcp)

	select {
	case err := <-errChan:
		log.Println("exited with an error: ", err.Error())
		return err
	case sig := <-sigCapt:
		log.Println("exited with a signal: ", sig.String())
	}

	return nil
}
