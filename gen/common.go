package gen

type CmdGen interface {
	CmdHandle()
}

func Generate(cmd CmdGen) {
	cmd.CmdHandle()
}