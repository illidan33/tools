package cmd

import (
	"github.com/illidan33/tools/common"
	dao2 "github.com/illidan33/tools/kiple/common"
	"github.com/spf13/cobra"
)

var kipleInterfaceCheck dao2.CmdKipleInterfaceCheck
var cmdkipleInterfaceCheck = &cobra.Command{
	Use:   "daosync",
	Short: "sync funcs from impl to dao interface",
	Run: func(cmd *cobra.Command, args []string) {
		common.CmdDo(&kipleInterfaceCheck)
	},
}

func init() {
	cmdkipleInterfaceCheck.Flags().StringVarP(&kipleInterfaceCheck.InterfaceName, "interfaceName", "i", "", "(required) the interface name which you want to sync")
	cmdkipleInterfaceCheck.Flags().StringVarP(&kipleInterfaceCheck.ModelName, "moduleName", "m", "", "(required) the module name which you want to generate from")
	cmdkipleInterfaceCheck.Flags().BoolVarP(&kipleInterfaceCheck.IsDebug, "debug", "d", false, "open debug flag,default: false")

	cmdKiple.AddCommand(cmdkipleInterfaceCheck)
}
