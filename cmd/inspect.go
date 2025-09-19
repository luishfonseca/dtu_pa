package cmd

import (
	"fmt"

	"github.com/luishfonseca/dtu_pa/analyser"

	"github.com/spf13/cobra"

	"os"
	"path/filepath"
)

// inspectCmd represents the inspect command
var inspectCmd = &cobra.Command{
	Use:   "inspect [flags] file",
	Short: "Print the parsed compiled java class in a human readable format",
	Args:  cobra.ExactArgs(1),
	PreRunE: func(_ *cobra.Command, args []string) (err error) {
		if filepath.Ext(args[0]) != ".class" {
			return fmt.Errorf("file %s must have a .class extension", args[0])
		}

		if _, err = os.Stat(args[0]); os.IsNotExist(err) {
			return fmt.Errorf("file %s does not exist", args[0])
		}

		if args[0], err = filepath.Abs(args[0]); err != nil {
			return fmt.Errorf("could not get absolute path of %s: %w", args[0], err)
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if err := analyser.New(args[0]).Inspect(); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(inspectCmd)
}
