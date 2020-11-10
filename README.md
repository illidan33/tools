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

根据ddl文件批量生成对应的golang struc
```

- gen method
```
//go:generate tools gen method --name xxx

tools gen method -h

根据model生成常用的数据库通用func
```

- gen client
```
//go:generate tools gen client --url "http://xxx.json" -n "xxx"

tools gen client -h   

根据swagger文件反解析生成对应的api
```

- kiple daocreate
```
//go:generate tools kiple daocreate -i UserProfilesDao -e "../entity/xxx.go"

tools kiple daocreate -h

根据entity批量生成kiple对应的常用的数据库操作func
```

- kiple daosync
```
//go:generate tools kiple daosync -i xxx -m xxx

同步model的所有funcs到interface中
```