package cmd

import (
	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:               "mesg-core",
	Short:             "Core",
	DisableAutoGenTag: true,
}
