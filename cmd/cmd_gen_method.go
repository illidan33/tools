package cmd

import (
	"github.com/spf13/cobra"
	"github.com/illidan33/tools/gen"
	"github.com/illidan33/tools/gen/method"
)

var genMethod method.CmdGenMethod
var cmdGenMethod = &cobra.Command{
	Use:   "method",
	Short: "generate gorm functions of gorm model",
	Run: func(cmd *cobra.Command, args []string) {
		gen.Generate(&genMethod)
	},
}

func init() {
	cmdGenMethod.Flags().StringVarP(&genMethod.ModelName, "name", "", "", "(required) name of source model")
	cmdGenMethod.Flags().BoolVarP(&genMethod.IsDebug, "debug", "", false, "open debug flag,default: false")

	cmdGen.AddCommand(cmdGenMethod)
}
