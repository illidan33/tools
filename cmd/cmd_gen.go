package cmd

import (
	"github.com/spf13/cobra"
)

var cmdGen = &cobra.Command{
	Use:   "gen",
	Short: "generate tools",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	cmdRoot.AddCommand(cmdGen)
}
