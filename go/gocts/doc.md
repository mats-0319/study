# Gocts

根据go语言定义的接口结构，生成相应的ts发送http请求的代码

（ts使用axios包发送http请求）

程序可以生成axios配置文件(`config.ts`)、发送http请求的函数文件(`*.http.ts`)以及相关结构体和枚举的定义文件(`*.go.ts`)

## 安装

`go install github.com/mats9693/study/go/gocts`

## 使用

> 参考`gocts/sample/manual.md`

## 目录

- data: 生成器结构定义（生成器：保存go代码解析结果并据此生成ts代码的结构）
- generate_ts: 根据生成器生成ts代码
- initialize: 初始化功能，包含`-i`/`-g`等启动参数的功能实现
- parse: 解析go代码
- sample: 使用示例
- test: 测试
- utils: 工具库
- .run: goland运行配置，可以快速在sample路径构建并运行

## 约定

- 所有用到的资源应全部在一个包内
- 使用`const URI_[uri_name]=[uri_value]`结构定义接口uri
    - e.g. `const URI_ListUser = "/user/list"`
- `type xxx struct {}`将被识别为自定义结构类型，其中结构体名称形如`xxxReq`/`xxxRes`者，
  将被识别为对应接口的输入/输出参数，每个定义的接口，其uri、接口输入输出参数名称要求一致。
    - e.g.
    ```code
    const URI_ListUser = "/user/list"
    
    type ListUserReq struct {
      Operator     string       `json:"operator"`
      ListIdentify UserIdentify `json:"list_identify"`
      Page         Pagination   `json:"page"`
    }
    
    type ListUserRes struct {
      Res     ResBase  `json:"res"`
      Summary int64    `json:"summary"`
      Users   []string `json:"users"`
    }
    ```
    - 要求go结构体的每一个字段都有`json tag`，我们使用该tag作为结构体的名称
- 使用下方示例中的方式定义一个枚举：
    - e.g.
    ```code 
    type UserIdentify int8
    
    const (
      UserIdentify_Second UserIdentify = 20
      UserIdentify_Value0 UserIdentify = 10
      UserIdentify_Value2 UserIdentify = 40
    ) 
    ```
    - 枚举可以以类型别名的方式定义(`type UserIdentify = int8`)
    - 单个枚举项应符合结构：`[enum_item_name] [enum_name] = [enum_value]`
        - 枚举项的名称应以枚举名称+下划线开头，形如：`[enum_name]_[xxx]`
        - 每一个枚举项应显式给出枚举名称与枚举值

有这些约定是因为我们将会按照约定的格式来解析go代码，相应的，因为对格式作出了要求，我们也提供相应的代码生成工具（通过启动参数，详见使用示例）

## 开发计划

todo

## 说明

- `.http.ts`文件中，如果一个参数是自定义类型，是保留其自定义类型，还是将其展开成多个字段？  
  我们选择保留其自定义类型，因为一个请求的输入参数是结构体，
  常见情况是分页信息（包含页数和每页记录数）或创建请求的对象（例如创建一个用户，可能包含用户名、密码等）
  而这些情况其实都不适合展开字段：分页信息可以写一个通用的默认分页；
  而`传递一个想要创建的用户`应用场景中，把所有需要的参数写成结构体本身更有利于代码阅读。
- 为什么使用json格式在前后端之间传递参数？
  使用json格式无论在前端构造请求参数，还是在后端解析请求参数的时候都很方便
  （前端直接传一个object，后端则是直接调用json包的反序列化），
  也尝试过使用form格式传参，但是后端解析很麻烦。
  使用json格式的不便之处在于传流式参数（例如文件），这一点目前还是通过form格式做。
- 为什么所有的http请求方法都设计成了post方法？
  设计过程中，随着阅读了http请求结构、restful api风格，发现可能它们不适合直接照搬过来使用，举几个例子：
    - get作为幂等的接口，其结果可能被http缓存，但实际上，对于一些变化的表，我们不希望它缓存查询结果
    - http 1.0 只有get/post方法，为了兼容，一些流行的库put/delete方法甚至是使用post模拟的，这大可不必
