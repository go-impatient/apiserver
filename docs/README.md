## Go包管理工具, govendor的使用
### 安装
```
go get -u -v github.com/kardianos/govendor
```
### 使用
```
# 进入到项目目录
cd $GOPATH/src/github.com/moocss/apiserver

# 初始化vendor目录
govendor init
# add gin
govendor fetch github.com/gin-gonic/gin@v1.2

# 将GOPATH中本工程使用到的依赖包自动移动到vendor目录中
# 说明：如果本地GOPATH没有依赖包，先go get相应的依赖包
govendor add +external
或使用缩写： govendor add +e

# Go 1.6以上版本默认开启 GO15VENDOREXPERIMENT 环境变量，可忽略该步骤。
# 通过设置环境变量 GO15VENDOREXPERIMENT=1 使用vendor文件夹构建文件。
# 可以选择 export GO15VENDOREXPERIMENT=1 或 GO15VENDOREXPERIMENT=1 go build 执行编译
export GO15VENDOREXPERIMENT=1
```

## 运行
```
# 进入 apiserver 目录编译源代码
$ cd $GOPATH/src/github.com/moocss/apiserver
$ gofmt -w .
$ go tool vet .
$ go build -v .
$ ./apiserver
```

## Go语言完整的应用项目结构最佳实践

首先，项目下的目录和文件命名一律小写，有必要的用下划线 _，目录名一律单数形式，目录下的包名尽量与 目录名一致。
cmd 目录存放用于编译可运行程序的 main 源码，它又分成了子级目录，主要是考虑一个项目可能有多种可运行程序。
src 目录放主要源码，集中在这个目录主要是为了方便查找和替换。src 目录下除了 app.go，router.go 这种顶层入口，又细分如下：
  - util，工具函数，不会依赖本项目的任何其它逻辑，只会被其它源码依赖；
  - service，对外部服务的封装，如对 mongodb、redis、zipkin 等 client 的封装，也不会依赖本项目 util 之外的任何其它逻辑，只会被其它源码依赖；
  - schema，数据模型，与数据库无关，也不会依赖本项目 util 之外的任何其它逻辑，只会被其它源码依赖；
  - model，通常依赖 util，service 和 schema，实现对数据库操作的主要逻辑，各个 model 内部无相互依赖；
  - bll，Business logic layer，通常依赖 util，schema 和 model，通过组合调用 model 实现更复杂的业务逻辑；
  - api，API 接口，通常依赖 util，schema 和 bll，挂载于 Router 上，直接受理客户端请求、提取和验证数据，调用 bll 层处理数据，然后响应给客户端；
  - ctl，Controller，类似 api 层，通常依赖 util，schema 和 bll，挂载于 Router 上，为客户端响应 View 页面；
  - 其它如 auth、logger 等则是一些带状态的被其它组件依赖的全局性组件。

与 cmd、src 平级的目录可能还会有：web 前端源码目录；config 配置文件目录；vendor go 依赖包目录；dist 编译后的可执行文件目录；docs 文档目录；k8s k8s 配置文件目录等。

## 目录结构

```
apiserver
├── docs                        # 帮助文档和资料
├── main.go                     # 程序入口与可以运行程序
├── src                              #
│   ├── app.go                       #
│   ├── swagger.go                   #
│   ├── config                       # 配置文件统一存放目录;用来处理配置和配置文件的Go package, 
│   │   ├── config.go                # 用来处理配置
│   │   ├── config.yaml              # 配置文件
│   │   ├── server.crt               # TLS配置文件
│   │   └── server.key
│   ├── db.sql                       # 在部署新环境时，可以登录MySQL客户端，执行source db.sql创建数据库和表
│   ├── api                          # 类似MVC架构中的C，用来读取输入，并将处理流程转发给实际的处理函数，最后返回结果
│   │   ├── handler.go
│   │   ├── sd                       # 健康检查handler
│   │   │   └── check.go
│   │   └── user                     # 核心：用户业务逻辑handler
│   │       ├── create.go            # 新增用户
│   │       ├── delete.go            # 删除用户
│   │       ├── get.go               # 获取指定的用户信息
│   │       ├── list.go              # 查询用户列表
│   │       ├── login.go             # 用户登录
│   │       ├── update.go            # 更新用户
│   │       └── user.go              # 存放用户handler公用的函数、结构体等
│   ├── model                        # 数据库相关的操作统一放在这里，包括数据库初始化和对表的增删改查
│   │   ├── init.go                  # 初始化和连接数据库
│   │   ├── model.go                 # 存放一些公用的go struct
│   │   └── user.go                  # 用户相关的数据库CURD操作
│   ├── pkg                          # 引用的包
│   │   ├── auth                     # 认证包
│   │   │   └── auth.go
│   │   ├── constvar                 # 常量统一存放位置
│   │   │   └── constvar.go
│   │   ├── errno                    # 错误码存放位置
│   │   │   ├── code.go
│   │   │   └── errno.go
│   │   ├── token
│   │   │   └── token.go
│   │   └── version                  # 版本包
│   │       ├── base.go
│   │       ├── doc.go
│   │       └── version.go
│   ├── router                       # 路由相关处理
│   │   ├── middleware               # API服务器用的是Gin Web框架，Gin中间件存放位置
│   │   │   ├── auth.go
│   │   │   ├── header.go
│   │   │   ├── logging.go
│   │   │   └── requestid.go
│   │   └── router.go
│   ├── service                      # 实际业务处理函数存放位置
│   │   └── service.go
│   └── util                         # 工具类函数存放目录
│       ├── util.go
│       └── util_test.go
│   
├── Makefile                     # Makefile文件，一般大型软件系统都是采用make来作为编译工具
├── README.md                    # API目录README
└── vendor                       # vendor目录用来管理依赖包
    ├── github.com
    ├── golang.org
    ├── gopkg.in
    └── vendor.json
```
