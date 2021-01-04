module tools

go 1.14

replace (
	github.com/illidan33/tools => ../tools
	github.com/m2c/kiplestar => ../../m2c/kiplestar
)

require (
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/baiyubin/aliyun-sts-go-sdk v0.0.0-20180326062324-cfa1a18b161f // indirect
	github.com/dave/dst v0.26.0
	github.com/dave/kerr v0.0.0-20170318121727-bc25dd6abe8e
	github.com/dave/ktest v1.1.3 // indirect
	github.com/fatih/structtag v1.2.0
	github.com/ghodss/yaml v1.0.0
	github.com/go-redis/redis/v8 v8.4.0 // indirect
	github.com/illidan33/tools v0.0.0-00010101000000-000000000000 // indirect
	github.com/jinzhu/gorm v1.9.14
	github.com/kataras/iris/v12 v12.1.8
	github.com/m2c/kiplestar v0.0.1
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826
	github.com/nats-io/nats-server/v2 v2.1.9 // indirect
	github.com/spf13/cobra v1.0.0
	github.com/swaggo/swag v1.6.5
	github.com/valyala/fasthttp v1.16.0
	golang.org/x/tools v0.0.0-20200509030707-2212a7e161a5
)
