package cmd

import (
	"github.com/spf13/cobra"
	"tools/common"
	"tools/gen/swagger"
)

var kipleSwagger swagger.CmdKipleSwagger
var cmdKipleSwagger = &cobra.Command{
	Use:   "swag",
	Short: "generate swagger doc from controller(for kiple)",
	Run: func(cmd *cobra.Command, args []string) {
		common.CmdDo(&kipleSwagger)
	},
}

func init() {
	cmdKipleSwagger.Flags().StringVarP(&kipleSwagger.Controller, "controller", "", "./controller", "(required) the controller dir which you want to generate swagger from, relative to the main file.")
	cmdKipleSwagger.Flags().StringVarP(&kipleSwagger.Pojo, "pojo", "", "./pojo", "(required) the pojo dir which you want to generate swagger from, relative to the main file.")
	cmdKipleSwagger.Flags().Uint8VarP(&kipleSwagger.IsInit, "init", "", 0, "only generate swagger tags into file or into dir 'docs'.0 写入docs文件夹 1 写入controller文件-跳过已存在；2 覆盖写入；3 全部清除；")
	cmdKipleSwagger.Flags().BoolVarP(&kipleSwagger.IsDebug, "debug", "", false, "open debug flag")

	cmdGen.AddCommand(cmdKipleSwagger)
}
