package cmd

import (
	"tools/common"
	"tools/gen/client"

	"github.com/spf13/cobra"
)

var genClient client.CmdGenClient
var cmdGenClient = &cobra.Command{
	Use:   "client",
	Short: "Generate client api from swagger doc",
	Run: func(cmd *cobra.Command, args []string) {
		common.CmdDo(&genClient)
	},
}

func init() {
	cmdGenClient.Flags().StringVarP(&genClient.DocUrl, "url", "u", "", "(required) Generate client from swagger url")
	cmdGenClient.Flags().StringVarP(&genClient.ServiceName, "name", "n", "", "(required) Generate client name")
	cmdGenClient.Flags().BoolVarP(&genClient.IsDebug, "debug", "d", false, "open debug flag,default: false")

	//cmdGen.AddCommand(cmdGenClient)
}
