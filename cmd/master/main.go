package main

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"

	"github.com/parinpan/romusha/internal/app/participant"
)

func main() {
	ctx := context.Background()

	rc := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	pcp := participant.NewParticipant(rc)

	if err := pcp.Watch(ctx, participant.AddParticipant(pcp)); err != nil {
		log.Fatal(err)
	}
}
