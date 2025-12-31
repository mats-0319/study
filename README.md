# 学习笔记

## 目录

- demo：用于技术验证，一般不会直接调用
    - [x] server-sent event
    - [x] grpc
- go
    - [x] doc
    - [x] encrypt transmission：加密传输，可以让两个人在不安全的通信信道上安全的传递信息
    - [x] generate avatar：生成像素风格对称图片，可用作默认头像
    - [x] gocts：根据go语言定义的接口结构，生成相应的ts发送http请求的代码
    - [x] totp：概念学习与go语言demo，android章节有完整的应用
    - [x] utils
    - [x] listen_bilibili：使用B站作为音源的听歌工具，[代码地址](https://github.com/mats0319/listen_bilibili)
    - [x] unnamed_plan：微服务架构示例项目，与微服务相关的内容尽量全部手写、不调用第三方代码库，
      [代码地址](https://github.com/mats0319/unnamed_plan)
- vue3
    - [x] vue doc：官方文档阅读笔记
    - [x] vue ecology：vue生态，包括但不限于pinia、vite config、ts编译选项
    - [x] vue router
    - [ ] eslint
- android
    - [x] dart：简单过一遍dart语法，知道dart代码大概是什么样子的
    - [x] flutter：学习flutter官方文档提供的codelab
    - [x] totp：使用dart+flutter编写的符合RFC 6238标准的应用，[代码地址](https://github.com/mats0319/totp)
- 笔记
    - [x] 全栈开发：总结自己的工作经历和对软件开发的认识
    - [x] github使用记录
    - [x] makefile学习笔记
    - [x] 正则表达式学习笔记
    - [x] shell学习笔记
    - [x] travis-ci持续集成+codecov代码覆盖率使用记录

## 为release计算hash

如果一个项目的编译过程比较复杂，通常考虑将编译好的内容发布到release，而发布编译好的内容，又通常需要一个hash防伪。
记录如何计算一个文件/文件夹的hash：

- 计算文件hash：
    - linux：`sha1sum [file name]`
    - windows: `Get-FileHash -Path [file name] -Algorithm SHA1`
- 计算文件夹hash：
    - linux：`find [folder path] -type f -print0 | xargs -0 sha1sum | sha1sum`
    - windows：`Get-ChildItem -Path [folder path] -File -Recurse | ForEach-Object { Get-FileHash $_.FullName 
      -Algorithm SHA1 } | Sort-Object Path | ForEach-Object { $_.Hash } | Get-FileHash -Algorithm SHA1`
    - windows也可以装一个git，然后使用linux的命令

计算文件夹hash，它的过程是找到文件夹里的所有文件（排除文件夹、链接等），计算每个文件的hash，
然后再对这些hash计算一次hash，得到的就是文件夹的hash
