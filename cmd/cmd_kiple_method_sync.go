package cmd

import (
	"tools/common"
	"tools/gen/method_sync"
	"github.com/spf13/cobra"
)

var kipleInterfaceCheck method_sync.CmdKipleInterfaceCheck
var cmdkipleInterfaceCheck = &cobra.Command{
	Use:   "msync",
	Short: "sync funcs from impl to interface(for kiple)",
	Run: func(cmd *cobra.Command, args []string) {
		common.CmdDo(&kipleInterfaceCheck)
	},
}

func init() {
	cmdkipleInterfaceCheck.Flags().StringVarP(&kipleInterfaceCheck.InterfaceName, "interface", "i", "", "(required) the interface name which you want to sync")
	cmdkipleInterfaceCheck.Flags().StringVarP(&kipleInterfaceCheck.ModelName, "module", "m", "", "(required) the module name which you want to generate from")
	cmdkipleInterfaceCheck.Flags().BoolVarP(&kipleInterfaceCheck.IsDebug, "debug", "d", false, "open debug flag,default: false")

	cmdGen.AddCommand(cmdkipleInterfaceCheck)
}
