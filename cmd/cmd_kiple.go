package cmd

import (
	"github.com/spf13/cobra"
)

var cmdKiple = &cobra.Command{
	Use:   "kiple",
	Short: "generate code of kiple",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	cmdRoot.AddCommand(cmdKiple)
}
