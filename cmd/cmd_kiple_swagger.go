package cmd

import (
	"github.com/illidan33/tools/common"
	"github.com/illidan33/tools/kiple/swagger"
	"github.com/spf13/cobra"
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
	cmdKipleSwagger.Flags().StringVarP(&kipleSwagger.ServideDir, "path", "p", "", "(required) the service dir path which you want to generate swagger from")
	cmdKipleSwagger.Flags().BoolVarP(&kipleSwagger.IsDebug, "debug", "d", false, "open debug flag,default: false")

	cmdKiple.AddCommand(cmdKipleSwagger)
}
