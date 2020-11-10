package cmd

import (
	"github.com/illidan33/tools/common"
	"github.com/illidan33/tools/gen/client"

	"github.com/spf13/cobra"
)

var genClient client.CmdGenClient
var cmdGenClient = &cobra.Command{
	Use:   "client",
	Short: "Generate swagger doc to client api",
	Run: func(cmd *cobra.Command, args []string) {
		common.CmdDo(&genClient)
	},
}

func init() {
	cmdGenClient.Flags().StringVarP(&genClient.DocUrl, "url", "u", "", "(required) Generate client from swagger url")
	cmdGenClient.Flags().StringVarP(&genClient.ServiceName, "name", "n", "", "(required) Generate client name")
	cmdGenClient.Flags().BoolVarP(&genClient.IsDebug, "debug", "d", false, "open debug flag,default: false")

	cmdGen.AddCommand(cmdGenClient)
}
