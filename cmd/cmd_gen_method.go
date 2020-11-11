package cmd

import (
	"github.com/illidan33/tools/common"
	"github.com/illidan33/tools/gen/method"
	"github.com/spf13/cobra"
)

var genMethod method.CmdGenMethod
var cmdGenMethod = &cobra.Command{
	Use:   "method",
	Short: "generate db handle functions for one model",
	Run: func(cmd *cobra.Command, args []string) {
		common.CmdDo(&genMethod)
	},
}

func init() {
	cmdGenMethod.Flags().StringVarP(&genMethod.ModelName, "name", "n", "", "(required) name of source model")
	cmdGenMethod.Flags().StringVarP(&genMethod.ModelFile, "filepath", "f", "", "(required) filepath of source model")
	cmdGenMethod.Flags().BoolVarP(&genMethod.IsDebug, "debug", "d", false, "open debug flag (default: false)")

	cmdGen.AddCommand(cmdGenMethod)
}
