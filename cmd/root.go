package cmd

import (
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dtu_pa",
	Short: "Program Analyser for DTU's 25/26 edition of Program Analysis",
}
