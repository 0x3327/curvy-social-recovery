package recoverCmd

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

const (
	outputFlag    = "output"
	thresholdFlag = "threshold"
	patternFlag   = "pattern"
)

func NewCommand() *cli.Command {
	return &cli.Command{
		Name:      "recover",
		Aliases:   []string{"r"},
		Usage:     "Recover the private (k,v) pair from the given shares",
		ArgsUsage: "<share directory>",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    outputFlag,
				Aliases: []string{"o"},
				Usage:   "Output path for (k,v) JSON file",
			},
			&cli.IntFlag{
				Name:      thresholdFlag,
				Aliases:   []string{"t"},
				Usage:     "Number of shares needed for recovery",
				Required:  true,
				Validator: thresholdValidator,
			},
			&cli.StringFlag{
				Name:    patternFlag,
				Aliases: []string{"p"},
				Usage:   "Pattern for matching share files",
				Value:   "share*.json",
			},
		},
		Action: recoverKey,
	}
}

func recoverKey(_ context.Context, cmd *cli.Command) error {
	t := cmd.Int(thresholdFlag)
	keyPath := cmd.String(outputFlag)

	if cmd.Args().Len() != 1 {
		return errors.New("invalid number of arguments")
	}

	sharePath := cmd.Args().First()
	err := shareDirValidator(sharePath)
	if err != nil {
		return fmt.Errorf("unable to validete share dir: %w", err)
	}

	// Find and print out all žson files in dir
	sharePattern := cmd.String(patternFlag)
	shareFilePaths, err := filepath.Glob(filepath.Join(sharePath, sharePattern))
	if err != nil {
		return fmt.Errorf("unable to find shares: %w", err)
	}

	if len(shareFilePaths) < 2 {
		fmt.Println("Not enough shares found.\nAborting.")
		return nil
	}

	fmt.Println("Share files found: ")
	for _, shareFilePath := range shareFilePaths {
		fmt.Println("\t" + shareFilePath)
	}
	fmt.Println()

	// Read the share files and create a []Share slice
	shares := make([]keyrecovery.Share, 0, len(shareFilePaths))
	var shareFile commands.ShareFile
	for i := range shareFilePaths {
		err = shareFile.ReadFromFile(shareFilePaths[i])
		if err != nil {
			return fmt.Errorf("unable to read share file: %w", err)
		}
		shares = append(shares, keyrecovery.Share(shareFile))
	}

	// Recover the key
	skStr, vkStr, err := keyrecovery.Recover(int(t), shares)

	// Handle the errors (some should not be considered "error state")
	if err != nil {
		var duplicateError *keyrecovery.DuplicatePointInSharesError
		var verificationError *keyrecovery.RecoveredKeysDoNotMatchError

		if errors.As(err, &duplicateError) {
			idx := duplicateError.Idx
			pointStr := duplicateError.Point

			fmt.Printf("Point 0x%s in file %s is not unique, possible tampering.\nAborting.\n", pointStr, shareFilePaths[idx])

			return nil
		} else if errors.As(err, &verificationError) {
			fmt.Println("Tampering of shares detected.\nAborting.")
			return nil
		} else {
			return fmt.Errorf("unable to recover the key: %w", err)
		}
	}

	key := commands.KeyFile{
		SpendingKey: skStr,
		ViewingKey:  vkStr,
	}

	err = key.WriteFile(keyPath)
	if err != nil {
		return err
	}

	fmt.Println("Key successfully recovered!")

	return nil
}

func thresholdValidator(t int64) error {
	if int(t) < 2 {
		return errors.New("threshold must be greater than 1")
	}
	return nil
}

func shareDirValidator(sharePath string) error {
	info, err := os.Stat(sharePath)
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
