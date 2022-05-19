package main

import (
	"context"

	"github.com/parinpan/romusha/internal/app/server"
)

func main() {
	server.Start(context.Background())
}
