package cmd

import (
	"myprojects/tools/gen"
	"myprojects/tools/gen/client"

	"github.com/spf13/cobra"
)

var genClient client.CmdGenClient
var cmdGenClient = &cobra.Command{
	Use:   "client",
	Short: "generate doc to client",
	Run: func(cmd *cobra.Command, args []string) {
		gen.Generate(&genClient)
	},
}

func init() {
	cmdGenClient.Flags().StringVarP(&genClient.CmdGenClientDocUrl, "url", "", "", "(required) generate client from url")
	cmdGenClient.Flags().StringVarP(&genClient.CmdGenClientServiceName, "client-name", "n", "", "(required) generate client name")

	cmdGen.AddCommand(cmdGenClient)
}
