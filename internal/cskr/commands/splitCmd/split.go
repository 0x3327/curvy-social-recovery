package splitCmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/0x3327/curvy-social-recovery/internal/cskr/commands"
	keyrecovery "github.com/0x3327/curvy-social-recovery/key_recovery"
	"github.com/urfave/cli/v3"
)

// Flag names
const (
	outputFlag    = "output"
	nOfSharesFlag = "nOfShares"
	thresholdFlag = "threshold"
)

func NewCommand() *cli.Command {
	return &cli.Command{
		Name:      "split",
		Aliases:   []string{"s"},
		Usage:     "Generate the shares for the private (k,v) pair",
		ArgsUsage: "<JSON key file>",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:      nOfSharesFlag,
				Aliases:   []string{"n"},
				Required:  true,
				Usage:     "Number of shares",
				Validator: nOfSharesValidator,
			},
			&cli.IntFlag{
				Name:      thresholdFlag,
				Aliases:   []string{"t"},
				Required:  true,
				Usage:     "Number of shares required to reconstruct the (k,v) pair",
				Validator: thresholdValidator,
			},
			&cli.StringFlag{
				Name:      outputFlag,
				Aliases:   []string{"o"},
				Value:     ".",
				Usage:     "Output directory",
				Validator: outputValidator,
			},
		},
		Action: split,
	}
}

func split(_ context.Context, cmd *cli.Command) error {
	if cmd.Args().Len() != 1 {
		return errors.New("invalid number of arguments")
	}

	n := cmd.Int(nOfSharesFlag)
	t := cmd.Int(thresholdFlag)
	outPath := cmd.String(outputFlag)

	if 2*t <= n {
		fmt.Println("Threshold should be larger than nOfShares/2.\nAborting.")
		return nil
	}

	inPath := cmd.Args().First()

	// Load (k,v) from file
	var key commands.KeyFile
	err := key.LoadFromFile(inPath)
	if err != nil {
		return fmt.Errorf("failed to load key from file: %w", err)
	}

	// Generate shares
	shares, err := keyrecovery.Split(int(t), int(n), key.SpendingKey, key.ViewingKey)
	if err != nil {
		return fmt.Errorf("could not split shares: %w", err)
	}

	// Write shares to files
	var shareFile commands.ShareFile
	for i, share := range shares {
		shareFile.FromShare(share)
		err = shareFile.WriteFile(filepath.Join(outPath, "share"+fmt.Sprintf("%02d", i+1)+".json"))
		if err != nil {
			return fmt.Errorf("unable to write share file: %w", err)
		}
	}

	fmt.Println("Shares successfully generated!")

	return nil
}

func nOfSharesValidator(n int64) error {
	if int(n) < 2 {
		return errors.New("number of shares must be greater than 1")
	}
	return nil
}

func thresholdValidator(t int64) error {
	if int(t) < 2 {
		return errors.New("threshold must be greater than 1")
	}
	return nil
}

func outputValidator(outPath string) error {
	info, err := os.Stat(outPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return errors.New("directory does not exist")
		} else {
			return err
		}
	}

	if !info.IsDir() {
		return errors.New("is not a directory")
	}

	return nil
}
