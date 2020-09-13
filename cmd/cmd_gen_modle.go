package cmd

import (
	"github.com/spf13/cobra"
	"myprojects/tools/gen/modle"
)

var flags modle.CmdGenModleFlags
var cmdGenModle = &cobra.Command{
	Use:   "modle",
	Short: "generate ddl sql to struct",
	Run: func(cmd *cobra.Command, args []string) {
		flags.CmdHandle()
	},
}

func init() {
	cmdGenModle.Flags().StringVarP(&flags.CmdGenModleName, "modle", "m", "", "generate modle with name (default table name)")
	cmdGenModle.Flags().StringVarP(&flags.CmdGenModleFilePath, "file", "f", "", "(required) generate modle from file path")
	cmdGenModle.Flags().BoolVarP(&flags.CmdGenModleWithGormTag, "gorm", "", true, "generate struct with gorm tag or not")
	cmdGenModle.Flags().BoolVarP(&flags.CmdGenModleWithSimpleGormTag, "gorm-simple", "", false, "generate struct with simple gorm tag or not")
	cmdGenModle.Flags().BoolVarP(&flags.CmdGenModleWithJsonTag, "json", "", true, "generate struct with json tag or not")
	cmdGenModle.Flags().BoolVarP(&flags.CmdGenModleWithDefaultTag, "default", "", false, "generate struct with default tag or not")

	cmdGen.AddCommand(cmdGenModle)
}
