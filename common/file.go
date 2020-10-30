package common

type CmdFilePath struct {
	CmdFileName string
	PackageName string
	SysArch     string // 系统架构
	Sys         string // 操作系统
	CmdLine     string // 当前执行命令所在文件中的行号
	CmdDir      string // 当前执行命令所在文件的路径
}
