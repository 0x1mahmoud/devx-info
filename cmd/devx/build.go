package main

import (
	"fmt"

	"cuelang.org/go/cue/errors"
	"github.com/devopzilla/devx/pkg/client"
	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:     "build [environment]",
	Short:   "Build DevX magic for the specified environment",
	Args:    cobra.ExactArgs(1),
	Aliases: []string{"do"},
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := client.Run(args[0], configDir, stackPath, buildersPath, reserve, dryRun, server, noStrict, stdout); err != nil {
			return fmt.Errorf(errors.Details(err, nil))
		}
		return nil
	},
}
