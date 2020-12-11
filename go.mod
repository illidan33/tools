module tools

go 1.14

replace (
	github.com/illidan33/tools => ../tools
	github.com/m2c/kiplestar => ../../m2c/kiplestar
)

require (
	github.com/dave/dst v0.26.0
	github.com/dave/kerr v0.0.0-20170318121727-bc25dd6abe8e
	github.com/dave/ktest v1.1.3 // indirect
	github.com/fatih/structtag v1.2.0
	github.com/ghodss/yaml v1.0.0
	github.com/go-redis/redis/v8 v8.4.0 // indirect
	github.com/illidan33/tools v0.0.0 // indirect
	github.com/jinzhu/gorm v1.9.14
	github.com/m2c/kiplestar v0.0.1
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826
	github.com/spf13/cobra v1.0.0
	github.com/valyala/fasthttp v1.16.0
)
