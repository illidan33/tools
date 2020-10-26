package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var cmdRoot = &cobra.Command{
	Use:   "tools",
	Short: "tools for generate code",
}

func Execute() {
	if err := cmdRoot.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
