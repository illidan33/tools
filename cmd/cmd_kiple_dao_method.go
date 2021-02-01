package cmd

import (
	"tools/common"
	"tools/gen/dao_create"
	"github.com/spf13/cobra"
)

var kipleDao dao_create.CmdKipleDao
var cmdkipleDao = &cobra.Command{
	Use:   "dmethod",
	Short: "generate methods of entity dao(for kiple)",
	Run: func(cmd *cobra.Command, args []string) {
		common.CmdDo(&kipleDao)
	},
}

func init() {
	cmdkipleDao.Flags().StringVarP(&kipleDao.Template.InterfaceName, "interface", "i", "", "(required) the interface name which you want to create")
	cmdkipleDao.Flags().StringVarP(&kipleDao.Template.ModelName, "moduleName", "m", "", "(required) the module name which you want to generate from")
	cmdkipleDao.Flags().StringVarP(&kipleDao.Template.DbName, "db", "", "", "(required) the db name which you want to use")
	cmdkipleDao.Flags().StringVarP(&kipleDao.Entity, "entity", "e", "", "(required) the entity place where generating code from")
	cmdkipleDao.Flags().BoolVarP(&kipleDao.IsDebug, "debug", "d", false, "open debug flag,default: false")

	cmdGen.AddCommand(cmdkipleDao)
}
