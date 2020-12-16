package cmd

import (
	"github.com/spf13/cobra"
	"tools/common"
	"tools/kiple/swagger"
)

var kipleSwagger swagger.CmdKipleSwagger
var cmdKipleSwagger = &cobra.Command{
	Use:   "swag",
	Short: "generate swagger doc",
	Run: func(cmd *cobra.Command, args []string) {
		common.CmdDo(&kipleSwagger)
	},
}

func init() {
	cmdKipleSwagger.Flags().StringVarP(&kipleSwagger.Controller, "controller", "", "./controller", "(required) the controller dir name which you want to generate swagger from")
	cmdKipleSwagger.Flags().StringVarP(&kipleSwagger.Pojo, "pojo", "", "./pojo", "(required) the pojo dir name which you want to generate swagger from")
	cmdKipleSwagger.Flags().Uint8VarP(&kipleSwagger.IsInit, "init", "", 0, "only generate swagger tag into file or into dir 'docs'.")
	cmdKipleSwagger.Flags().BoolVarP(&kipleSwagger.IsDebug, "debug", "", false, "open debug flag")

	cmdKiple.AddCommand(cmdKipleSwagger)
}
