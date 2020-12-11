package cmd

import (
	"tools/common"
	"tools/gen/model"
	"github.com/spf13/cobra"
)

var genModel model.CmdGenModel
var cmdGenModel = &cobra.Command{
	Use:   "model",
	Short: "convert ddl to golang struct with index.",
	Run: func(cmd *cobra.Command, args []string) {
		common.CmdDo(&genModel)
	},
}

func init() {
	cmdGenModel.Flags().StringVarP(&genModel.DdlFilePath, "file", "f", "", "(required) the path of ddl file")
	cmdGenModel.Flags().BoolVarP(&genModel.WithGormTag, "gorm", "", true, "tag has gorm or not(default true)")
	cmdGenModel.Flags().BoolVarP(&genModel.WithSimpleGormTag, "gmsimple", "", true, "tag has simple gorm or total gorm (default true，need gorm set to true)")
	cmdGenModel.Flags().BoolVarP(&genModel.WithJsonTag, "json", "", true, "tag has json or not (default true)")
	cmdGenModel.Flags().BoolVarP(&genModel.WithDefaultTag, "default", "", false, "tag has default or not (default false)")
	cmdGenModel.Flags().BoolVarP(&genModel.IsDebug, "debug", "d", false, "open debug flag (default false)")

	cmdGen.AddCommand(cmdGenModel)
}
