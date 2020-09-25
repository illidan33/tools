package cmd

import (
	"github.com/spf13/cobra"
	"myprojects/tools/gen"
	"myprojects/tools/gen/model"
)

var genModel model.CmdGenModel
var cmdGenModel = &cobra.Command{
	Use:   "model",
	Short: "generate ddl sql to struct",
	Run: func(cmd *cobra.Command, args []string) {
		gen.Generate(&genModel)
	},
}

func init() {
	cmdGenModel.Flags().StringVarP(&genModel.DdlFilePath, "file", "f", "", "(required) generate model from file path")
	cmdGenModel.Flags().BoolVarP(&genModel.WithGormTag, "gorm", "", true, "generate struct with gorm tag or not")
	cmdGenModel.Flags().BoolVarP(&genModel.WithSimpleGormTag, "gorm-simple", "", false, "generate struct with simple gorm tag or not")
	cmdGenModel.Flags().BoolVarP(&genModel.WithJsonTag, "json", "", true, "generate struct with json tag or not")
	cmdGenModel.Flags().BoolVarP(&genModel.WithDefaultTag, "default", "", false, "generate struct with default tag or not")

	cmdGen.AddCommand(cmdGenModel)
}
