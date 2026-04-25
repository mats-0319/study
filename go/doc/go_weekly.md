# Go Weekly阅读笔记

记录自己在go weekly了解到的优秀文章、资源以及工具等内容。[subscribe](https://golangweekly.com/)

## 阅读

- https://beej.us/guide/bglcs/html/split/ ：计算机科学学习方法论，个人编写的学习计算机科学的宏观方法论
- https://atlas9.dev/blog/soft-delete.html ：pg软删除的缺点及替代方案（触发器或预写入日志）
- https://github.com/gibbok/typescript-book ：一本简洁的ts手册
- https://eblog.fly.dev/ginbad.html ：“gin是一个糟糕的软件库”
- https://www.sophielwang.com/blog/jpeg ：详细介绍jpeg图片格式压缩的原理
- https://www.youtube.com/watch?v=hJ4S-5MirvU ：一个讲使用go+pg实现outbox模式的视频，我刚和gemini合作实现了该模式
- https://medium.com/@m0t9_/lets-add-a-conditional-expression-to-go-language-3eec7783388e ：
  通过为go语言添加`?:`条件表达式，介绍go语言编译器（词法分析、类型检查等环节）的工作过程和编辑方法
- https://github.blog/engineering/infrastructure/how-github-uses-ebpf-to-improve-deployment-safety/ ：
  github的源代码也托管在github上，这篇文章介绍了github如何使用ebpf-go解决循环依赖的问题
  （举个例子，github无法访问时，github的开发者将无法下载和上传github的源代码）
- https://blog.iamvedant.in/containers-are-not-magic-namespaces-from-scratch ，从零开始创建一个容器，介绍容器技术原理
- https://anirudhology.com/blog/building-raftly-reproducing-production-failures ，
  实现一个可控的raft共识，通过删除一个环节来模拟和理解对应环节的作用

## 代码/工具

- github.com/go-drift/drift：使用go语言编写移动端跨平台程序(android/ios)，代码风格和flutter非常像。该项目由个人开发，且处于早期阶段
- github.com/aperturerobotics/goscript：将go代码在ast级转换为ts
- github.com/hajimehoshi/ebiten：go语言2D引擎
- github.com/templui/templui：一个让人可以使用go语言编写网页的库，内置支持tailwindCss
- github.com/markel1974/godoom：go语言3d引擎
