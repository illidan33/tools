package cmd

import (
	"github.com/illidan33/tools/common"
	"github.com/illidan33/tools/kiple/dao"
	"github.com/spf13/cobra"
)

var kipleDao dao.CmdKipleDao
var cmdkipleDao = &cobra.Command{
	Use:   "daocreate",
	Short: "generate methods of entity dao",
	Run: func(cmd *cobra.Command, args []string) {
		common.CmdDo(&kipleDao)
	},
}

func init() {
	cmdkipleDao.Flags().StringVarP(&kipleDao.InterfaceName, "interface", "i", "", "(required) the interface name which you want to create")
	cmdkipleDao.Flags().StringVarP(&kipleDao.Entity, "entity", "e", "", "(required) the entity place where generating code from")
	cmdkipleDao.Flags().BoolVarP(&kipleDao.IsDebug, "debug", "d", false, "open debug flag,default: false")

	cmdKiple.AddCommand(cmdkipleDao)
}
