package cmd

import (
	"tools/common"
	"tools/kiple/method_sync"
	"github.com/spf13/cobra"
)

var kipleInterfaceCheck method_sync.CmdKipleInterfaceCheck
var cmdkipleInterfaceCheck = &cobra.Command{
	Use:   "methodsync",
	Short: "sync funcs from impl to interface",
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
