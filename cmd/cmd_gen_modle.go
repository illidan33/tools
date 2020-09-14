package cmd

import (
	"github.com/spf13/cobra"
	"myprojects/tools/gen/model"
)

var flags model.CmdGenmodelFlags
var cmdGenmodel = &cobra.Command{
	Use:   "model",
	Short: "generate ddl sql to struct",
	Run: func(cmd *cobra.Command, args []string) {
		flags.CmdHandle()
	},
}

func init() {
	cmdGenmodel.Flags().StringVarP(&flags.CmdGenmodelName, "model", "m", "", "generate model with name (default table name)")
	cmdGenmodel.Flags().StringVarP(&flags.CmdGenmodelFilePath, "file", "f", "", "(required) generate model from file path")
	cmdGenmodel.Flags().BoolVarP(&flags.CmdGenmodelWithGormTag, "gorm", "", true, "generate struct with gorm tag or not")
	cmdGenmodel.Flags().BoolVarP(&flags.CmdGenmodelWithSimpleGormTag, "gorm-simple", "", false, "generate struct with simple gorm tag or not")
	cmdGenmodel.Flags().BoolVarP(&flags.CmdGenmodelWithJsonTag, "json", "", true, "generate struct with json tag or not")
	cmdGenmodel.Flags().BoolVarP(&flags.CmdGenmodelWithDefaultTag, "default", "", false, "generate struct with default tag or not")

	cmdGen.AddCommand(cmdGenmodel)
}
