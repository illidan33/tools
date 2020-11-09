# tools

- install
```
git clone git@github.com:illidan33/tools.git

cd tools

go install
```
the cmd file 'tools' will be installed to your $GOBIN directory.

### cmd

- gen model
```
//go:generate tools gen model -f xxx

tools gen model -h      

doc: 
generate ddl sql to struct

Usage:
  tools gen model [flags]

Flags:
      --debug         open debug flag (default false)
      --default       generate struct with default tag or not (default false)
  -f, --file string   (required) generate model from file path, make sure not has single quote in your field comment of ddl string.
      --gmsimple      generate struct with simple gorm tag or not (default true) (default true)
      --gorm          generate struct with gorm tag or not (default true) (default true)
  -h, --help          help for model
      --json          generate struct with json tag or not (default true) (default true)
```

- gen method
```
//go:generate tools gen method --name xxx

tools gen method -h

doc:
generate gorm functions of gorm model

Usage:
  tools gen method [flags]

Flags:
      --debug         open debug flag,default: false
  -h, --help          help for method
      --name string   (required) name of source model

```

- gen client
```
//go:generate tools gen client --url "http://xxx.json" -n "xxx"

tools gen client -h   

doc:                                                              
Generate swagger doc to client

Usage:
  tools gen client [flags]

Flags:
  -n, --client-name string   (required) Generate client name
      --debug                open debug flag,default: false
  -h, --help                 help for client
      --url string           (required) Generate client from swagger url
```

- kiple dao
```
//go:generate tools kiple dao -i UserProfilesDao -e "../entity/xxx.go"

tools kiple dao -h

doc:
generate methods of entity dao

Usage:
  tools kiple dao [flags]

Flags:
  -d, --debug                  open debug flag,default: false
  -e, --entity string          (required) the entity place where generating code from
  -h, --help                   help for dao
  -i, --interfaceName string   (required) the interface name which you want to create
```

- kiple interface
```
//go:generate tools kiple method -i xxx -m xxx

tools kiple interface -h

doc:
generate methods of interface

Usage:
  tools kiple interface [flags]

Flags:
  -d, --debug                  open debug flag,default: false
  -h, --help                   help for interface
  -i, --interfaceName string   (required) the interface name which you want to create
  -m, --moduleName string      (required) the module name which you want to generate from
```