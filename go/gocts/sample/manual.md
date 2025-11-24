# sample

```cmd
> go install github.com/mats9693/study/go/gocts
> gocts -h

2025/11/24 10:39:37 Options: 
(In this help, './go/' means our go files dir)
  -c string
        config file (default "./go/config.json")
  -g    generate go files from './go/init.json'
        
  -genFrom string
        generate go files from given file
  -h    this help
  -i    initializeFlag basic files
        overwrite './go/config.json' and './go/init.json'
  -v    show version
```

因为本工具有一些默认的约束条件（例如命名），所以我们为此提供了相应的支持（启动参数）：

```cmd
gocts -i
```

命令可以生成一个**配置文件`(config.json)`**和一个**go接口文件的初始化文件`(init.json)`**

根据你的需求自行编辑以上文件
> 比较重要的配置项有：  
> `base_url`：前端向哪儿发送请求（ip+端口）  
> `basic_go_type`：没有写在这儿的类型将被视作自定义类型，我们提前将常用的go语言内置类型写在这里了

初始化文件则是用于生成符合约定的go接口文件，默认的初始化文件及其生成的go接口文件如下：

```cmd
gocts -g （也可以通过`-genFrom`参数指定初始化文件）
```

```code
// init.json
{
  "package_name": "api",
  "files": [
    {
      "file_name": "demo",
      "api_list": [
        {
          "name": "ListUser",
          "uri": "/user/list"
        },
        {
          "name": "CreateUser",
          "uri": "/user/create"
        }
      ],
      "enum_list": [
        {
          "name": "UserIdentify",
          "number": 3
        }
      ]
    }
  ]
}

// demo.go
package api

type UserIdentify int8

const (
    UserIdentify_Value0 UserIdentify = 0
    UserIdentify_Value1 UserIdentify = 1
    UserIdentify_Value2 UserIdentify = 2
)

const URI_ListUser = "/user/list"

type ListUserReq struct {}

type ListUserRes struct {}

const URI_CreateUser = "/user/create"

type CreateUserReq struct {}

type CreateUserRes struct {}
```

> 我们为生成文件提供备份功能，假设你在修改`.json/.go`文件的过程中忘记了约定的格式，
> 可以通过`-i/-g`命令重新生成示例代码，而你的修改将被保存到同目录下`backup`文件夹。

根据你的需求，编辑go接口文件，例如我们对`demo`文件进行如下调整：（新增两个struct、修改两个struct）

```code
type ResBase struct {
IsSuccess bool   `json:"is_success"`
Err       string `json:"err"`
}

type Pagination struct {
PageNum  int `json:"page_num"`
PageSize int `json:"page_size"`
}

type ListUserReq struct {
Operator string     `json:"operator"`
Page     Pagination `json:"page"`
}

type ListUserRes struct {
Res     ResBase  `json:"res"`
Summary int64    `json:"summary"`
Users   []string `json:"users"`
}
```

最后生成ts代码

```cmd
gocts
```
