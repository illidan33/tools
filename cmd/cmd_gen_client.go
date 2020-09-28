package cmd

import (
	"myprojects/tools/gen"
	"myprojects/tools/gen/client"

	"github.com/spf13/cobra"
)

var genClient client.CmdGenClient
var cmdGenClient = &cobra.Command{
	Use:   "client",
	Short: "Generate swagger doc to client",
	Run: func(cmd *cobra.Command, args []string) {
		gen.Generate(&genClient)
	},
}

func init() {
	cmdGenClient.Flags().StringVarP(&genClient.DocUrl, "url", "", "", "(required) Generate client from swagger url")
	cmdGenClient.Flags().StringVarP(&genClient.ServiceName, "client-name", "n", "", "(required) Generate client name")
	cmdGenClient.Flags().BoolVarP(&genClient.IsDebug, "debug", "", false, "open debug flag,default: false")

	cmdGen.AddCommand(cmdGenClient)
}
