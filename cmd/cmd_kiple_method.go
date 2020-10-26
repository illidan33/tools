package cmd

import (
	"github.com/illidan33/tools/gen"
	"github.com/illidan33/tools/kiple/method"
	"github.com/spf13/cobra"
)

var kipleMethod method.CmdKipleMethod
var cmdKipleMethod = &cobra.Command{
	Use:   "method",
	Short: "generate gorm functions of gorm model",
	Run: func(cmd *cobra.Command, args []string) {
		gen.Generate(&kipleMethod)
	},
}

func init() {
	cmdKipleMethod.Flags().StringVarP(&kipleMethod.ModelName, "name", "", "", "(required) name of source model")
	cmdKipleMethod.Flags().StringVarP(&kipleMethod.Entity, "entity", "", "", "the entity place where generating code from (default '../entity/{ModelName}.go')")
	cmdKipleMethod.Flags().BoolVarP(&kipleMethod.IsDebug, "debug", "", false, "open debug flag,default: false")

	cmdKiple.AddCommand(cmdKipleMethod)
}
