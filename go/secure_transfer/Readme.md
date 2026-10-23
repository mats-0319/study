# 安全的信息传输工具 Secure Message Transfer Tool

本工具可以让两个人在不安全的通信信道上，以文件为载体，安全的传递信息

## 使用

> 以命令行程序使用，参考：`/cmd/cli/script/manual.md`
> 以图形化界面程序使用，参考：`/cmd/gui/components/utils.go helpText`

> 我们会提供一些平台的可执行程序，但我们建议你阅读源代码后自行编译。

> 声明：`/cmd/cli/sample/`路径下的私钥等文件仅作为演示使用

程序会在当前目录创建和读取文件，请在一个单独的路径中启动程序。

> 在 Windows + GoLand 上遇到的非代码问题（例如被杀软沙盒拖去多跑一次、两次工作目录不一致），参考：
> `/doc/troubleshooting_windows.md`

编译：以命令行形式使用无需外部依赖 (go install)，以GUI形式使用需要安装依赖

- fyne存在交叉编译问题：仅根据文档中要求下载的内容，无法在linux上编译windows的可执行程序，所以gui的编译脚本也仅供参考

[下载go、C编译器和系统显卡驱动程序](https://docs.fyne.io/started/quick/)

```cmd
go get fyne.io/fyne/v2@latest
go install fyne.io/tools/cmd/fyne@latest
```

### GUI

fyne依赖cgo，如果编译环境和你的系统的glibc版本不兼容，则可执行程序无法使用。  
因此我们提供的gui的二进制文件在你的设备上无法运行是有可能的，附使用截图。

![linux](doc/linux_use.png)
![windows](doc/windows_use.png)

## 技术准备

[见tech_prepare.md](doc/tech_prepare.md)

## 技术细节

[见design.md](doc/design.md)

## 开发计划

- 优化流式加密性能：根据基准测试结果，当前流式加密代码分配内存比直接读取全文还多，尝试优化它

## 测试

常用命令：

```cmd
go test ./test/ -v                                   # 功能测试
go test ./test/ -race                                # 竞态检查
go test ./test/ -run '^$' -bench . -benchtime 1s     # 只跑基准
```

基准测试执行结果：

```cmd
$ go test .\test\ -run '^$' -bench . -benchtime 1s
                  不执行常规测试        每个基准测试至少运行1s
                            执行所有基准测试
goos: windows
goarch: amd64
pkg: github.com/mats0319/secure_transfer/test
cpu: AMD Ryzen 7 7700 8-Core Processor 
BenchmarkEncrypt/once/8MiB-16        268     4380662 ns/op     1914.92 MB/s    16800051 B/op      43 allocs/op
BenchmarkEncrypt/stream/8MiB-16      685     1896901 ns/op     4422.27 MB/s       15237 B/op      96 allocs/op
                                     迭代次数  每次操作平均用时                     每次操作平均分配内存  每次操作平均分配对象
                                                               吞吐量（来自-benchmem)           
```
