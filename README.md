# tools

- install

需要提前配置GOBIN目录到系统环境变量；
```
git clone git@github.com:illidan33/tools.git

cd tools

go install

# show cmd of tools
tools -h
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

- gen dmethod

```
//go:generate tools gen dmethod [option]

根据entity批量生成kiple对应的常用的数据库操作func
```

- gen msync

```
//go:generate tools gen msync [option]

同步model的所有funcs到interface中
```

- gen swag
    - 需要在main.go文件所在目录运行命令;
    - 入参取名规则：{{func name}}Request
      - struc中每个字段需要定义标签'in'，标识属于'body' or 'query' or 'path' or 'header', query/path/header字段可定义标签: validate:"required";
      - 如果没有标签，默认整个struct为body（method不为GET）或者query（method为GET）；
    - 返回参数规则: {{func name}}Response;
    - BeforeActivation必须放置每个文件开头；
    - controller注册函数名称需要为：RegisterGlobalModel;