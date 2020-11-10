package common

import (
	"fmt"
)

type CmdGen interface {
	CmdHandle()
}

type CmdHandle interface {
	Init() error
	Parse() error
	String() string
}

func Generate(cmd CmdGen) {
	cmd.CmdHandle()
}

func CmdDo(cmd CmdHandle) {
	var err error
	if err = cmd.Init(); err != nil {
		panic(fmt.Errorf("\033[1;31mInit error: \n%s\033[0m\n", err.Error()))
	}
	if err = cmd.Parse(); err != nil {
		panic(fmt.Errorf("\033[1;31mInit error: \n%s\033[0m\n", err.Error()))
	}
	fmt.Printf("\033[1;32m%s generate success\033[0m\n", cmd.String())
}
