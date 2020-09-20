package cmd

import (
	"github.com/spf13/cobra"
	"myprojects/tools/gen"
	"myprojects/tools/gen/method"
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
	cmdGen.AddCommand(cmdGenMethod)
}
