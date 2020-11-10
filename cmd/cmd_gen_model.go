package cmd

import (
	"github.com/illidan33/tools/common"
	"github.com/illidan33/tools/gen/model"
	"github.com/spf13/cobra"
)

var genModel model.CmdGenModel
var cmdGenModel = &cobra.Command{
	Use:   "model",
	Short: "转换sql的ddl语句为golang的struct，并携带索引.",
	Run: func(cmd *cobra.Command, args []string) {
		common.CmdDo(&genModel)
	},
}

func init() {
	cmdGenModel.Flags().StringVarP(&genModel.DdlFilePath, "file", "f", "", "(required) ddl存放路径")
	cmdGenModel.Flags().BoolVarP(&genModel.WithGormTag, "gorm", "", true, "tag是否需要gorm标签 (default true)")
	cmdGenModel.Flags().BoolVarP(&genModel.WithSimpleGormTag, "gmsimple", "", true, "tag需要简化版gorm标签 or 完整版gorm标签 (default true，需要gorm配置为true才生效)")
	cmdGenModel.Flags().BoolVarP(&genModel.WithJsonTag, "json", "", true, "tag是否需要json标签 (default true)")
	cmdGenModel.Flags().BoolVarP(&genModel.WithDefaultTag, "default", "", false, "tag是否需要default标签 (default false)")
	cmdGenModel.Flags().BoolVarP(&genModel.IsDebug, "debug", "d", false, "open debug flag (default false)")

	cmdGen.AddCommand(cmdGenModel)
}
