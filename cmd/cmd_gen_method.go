package cmd

import (
	"github.com/illidan33/tools/common"
	"github.com/illidan33/tools/gen/method"
	"github.com/spf13/cobra"
)

var genMethod method.CmdGenMethod
var cmdGenMethod = &cobra.Command{
	Use:   "method",
	Short: "generate gorm functions of gorm model",
	Run: func(cmd *cobra.Command, args []string) {
		common.Generate(&genMethod)
	},
}

func init() {
	cmdGenMethod.Flags().StringVarP(&genMethod.ModelName, "name", "", "", "(required) name of source model")
	cmdGenMethod.Flags().BoolVarP(&genMethod.IsDebug, "debug", "", false, "open debug flag,default: false")

	cmdGen.AddCommand(cmdGenMethod)
}
