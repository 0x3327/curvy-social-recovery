package main

import (
	"context"
	"log"
	"os"

	"github.com/0x3327/curvy-social-recovery/internal/cskr"
)

func main() {
	cmd := cskr.NewCli()

	logger := log.New(os.Stderr, "Error: ", 0)

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		logger.Fatal(err)
	}
}
