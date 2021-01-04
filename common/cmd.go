package common

import (
	"fmt"
	"reflect"
)

type CmdGen interface {
	CmdHandle()
}

type CmdHandle interface {
	Init() error
	Parse() error
	String() string
}

func CmdDo(cmd CmdHandle) {
	var err error
	if err = cmd.Init(); err != nil {
		panic(fmt.Errorf("\033[1;31mInit error: \n%s\033[0m\n", err.Error()))
	}
	if err = cmd.Parse(); err != nil {
		panic(fmt.Errorf("\033[1;31mParse error: \n%s\033[0m\n", err.Error()))
	}
	fmt.Printf("\033[1;32m%s %s-generate success\033[0m\n", cmd.String(), reflect.TypeOf(cmd).String())
}
