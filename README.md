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
//go:generate tools gen model [option]

根据ddl文件批量生成对应的golang struct
1. 如果字段有为NUll的情况，可以替换字段类型为sql.NullXXX,比如sql.NullTime;
```

- gen method

```
//go:generate tools gen method [option]

根据model生成常用的数据库通用func
```

- gen client

```
//go:generate tools gen client [option]

根据swagger文件反解析生成对应的api
```

- kiple daocreate

```
//go:generate tools kiple daocreate [option]

根据entity批量生成kiple对应的常用的数据库操作func
```

- kiple methodsync

```
//go:generate tools kiple methodsync [option]

同步model的所有funcs到interface中
```

- kiple swag
    - 需要在main.go文件所在目录运行命令;
    - 入参规则：
        - 仅body或query参数时，入参struct取名规则：{{func name}}Request;
        - 当存在path、header参数时,入参struct取名规则：{{func name}}RequestWl，并且struc中每个字段需要定义标签'in'，标识属于'body' or 'query' or 'path' or 'header';
        - query/path/header字段可定义标签: validate:"required"；
    - 返回参数规则: {{func name}}Response;
    - controller注册func名称必须为"RegisterGlobalModel";
    - BeforeActivation必须放置每个文件开头；