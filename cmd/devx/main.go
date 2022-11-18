package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	configDir    string
	stackPath    string
	buildersPath string
	showTraitDef bool
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&configDir, "project", "p", ".", "project config dir")
	doCmd.PersistentFlags().StringVarP(&stackPath, "stack", "s", "stack", "stack field name in config file")
	doCmd.PersistentFlags().StringVarP(&buildersPath, "builders", "b", "builders", "builders field name in config file")
	discoverCmd.PersistentFlags().BoolVarP(&showTraitDef, "show-def", "s", false, "show trait definition")

	rootCmd.AddCommand(
		doCmd,
		projectCmd,
	)

	projectCmd.AddCommand(
		initCmd,
		updateCmd,
		validateCmd,
		discoverCmd,
		genCmd,
	)
}

var rootCmd = &cobra.Command{
	Use:   "devx",
	Short: "guku DevX cloud native self-service magic",
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
