package cskr

import (
	"io"

	"github.com/0x3327/curvy-social-recovery/internal/cskr/commands/recoverCmd"
	"github.com/0x3327/curvy-social-recovery/internal/cskr/commands/splitCmd"
	"github.com/urfave/cli/v3"
)

func NewCli() *cli.Command {
	cmd := &cli.Command{
		Name:  "cskr",
		Usage: "Curvy Social Key Recovery",
		Commands: []*cli.Command{
			splitCmd.NewCommand(),
			recoverCmd.NewCommand(),
		},
		ErrWriter: io.Discard,
	}

	return cmd
}
