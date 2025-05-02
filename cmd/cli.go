package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ParteeLabs/gomiger/config"
	"github.com/ParteeLabs/gomiger/generator"
	"github.com/urfave/cli/v3"
)

var rcPath string

var initCmd = &cli.Command{
	Name:    "init",
	Aliases: []string{"i"},
	Usage:   "generate migration and migrator",
	Action: func(_ context.Context, _ *cli.Command) error {
		rc, err := config.GetGomigerRC(rcPath)
		if err != nil {
			return fmt.Errorf("Cannot load the gomiger.rc file: %w", err)
		}
		if generator.IsSrcCodeInitialized(rc) {
			return fmt.Errorf("The source code is ALREADY INITIALIZED")
		}
		if err := generator.InitSrcCode(rc); err != nil {
			return fmt.Errorf("Cannot init gomiger: %w", err)
		}
		return nil
	},
}

var newCmd = &cli.Command{
	Name:    "new",
	Aliases: []string{"n"},
	Usage:   "generate a new migration",
	Action: func(_ context.Context, cmd *cli.Command) error {
		rc, err := config.GetGomigerRC(rcPath)
		if err != nil {
			return fmt.Errorf("Cannot load the gomiger.rc file: %w", err)
		}
		if !generator.IsSrcCodeInitialized(rc) {
			return fmt.Errorf("The source code is NOT INITIALIZED")
		}
		if err := generator.GenMigrationFile(rc, cmd.Args().Get(0)); err != nil {
			return fmt.Errorf("Cannot generate migration file: %w", err)
		}
		return nil
	},
}

func main() {
	cmd := &cli.Command{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "rc-path",
				Category:    "global",
				Value:       "./gomiger.rc.yaml",
				Usage:       "Path to the gomiger.rc file",
				Destination: &rcPath,
			},
		},
		Commands: []*cli.Command{
			initCmd,
			newCmd,
		},
	}
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
