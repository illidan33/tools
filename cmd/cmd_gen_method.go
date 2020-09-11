package cmd

import (
	"github.com/spf13/cobra"
	"myprojects/tools/gen/method"
)

var cmdGenMethodFlags method.CmdGenMethodFlags
var cmdGenMethod = &cobra.Command{
	Use:   "method",
	Short: "generate functions of gorm modle",
	Run: func(cmd *cobra.Command, args []string) {
		cmdGenMethodFlags.CmdHandle()
	},
}

func init() {
	cmdGenMethod.Flags().StringVarP(&cmdGenMethodFlags.CmdGenModleFilePath, "file-path", "f", "", "(required) file to generate method")

	cmdGen.AddCommand(cmdGenMethod)
}
