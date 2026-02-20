# Go tools

> [命令列表](https://go.dev/doc/cmd)提供有链接，直接指向各个命令最新版本的文档

`go help`可以查看go命令列表  
`go help <command>`可以查看一个命令的详细信息

## Commands

- go tool compile –S [main.go] >> [main.S] 生成main.go的go汇编，并保存在main.S文件中
- go test -coverprofile=[coverage.txt] 执行当前目录下的测试文件，保存结果为coverage.txt
- go tool cover -html=[coverage.txt] 生成网页形式的代码覆盖率报告
- go build -gcflags -m [file name] 查看编译过程中，优化了哪些代码

## vet

检查go代码，报告可疑结构(suspicious constructs)，例如参数与格式字符串不一致的printf调用

vet使用的启发式方法不保证所有的报告都是真正的问题，但它可以找到编译器未捕获的错误

vet不一定能找到所有的错误，所以不要依赖它判断程序正确性

需要在go module目录执行：`go vet [path]`

我的使用方式：`go vet -json ./...`

- 使用json格式的错误报告，主要是因为默认格式下，如果没有检查到可疑结构，vet不会有任何输出，初学时容易误解成命令没有执行

## fmt

格式化代码

`go fmt [flags] [path ...]`
or
`gofmt [flags] [path ...]`

flags:

```txt 
-d
	Do not print reformatted sources to standard output.
	If a file's formatting is different than gofmt's, print diffs
	to standard output.
-e
	Print all (including spurious) errors.
-l
	Do not print reformatted sources to standard output.
	If a file's formatting is different from gofmt's, print its name
	to standard output.
-r rule
	Apply the rewrite rule to the source before reformatting.
-s
	Try to simplify code (after applying the rewrite rule, if any).
-w
	Do not print reformatted sources to standard output.
	If a file's formatting is different from gofmt's, overwrite it
	with gofmt's version. If an error occurred during overwriting,
	the original file is restored from an automatic backup.
	
-cpuprofile filename
	Write cpu profile to the specified file.
```

我的使用方式：`gofmt -w -s -l .`，修改源文件、简化代码、列举修改了哪些文件

## 问题

- `gofmt`和`go fmt`结果可能不一样，因为两边调的不是一个可执行文件
