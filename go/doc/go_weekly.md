# Go Weekly阅读笔记

记录自己在go weekly了解到的优秀文章、资源以及工具等内容。[subscribe](https://golangweekly.com/)

## 阅读

- https://beej.us/guide/bglcs/html/split/ ：计算机科学学习方法论，个人编写的学习计算机科学的宏观方法论
- https://atlas9.dev/blog/soft-delete.html ：pg软删除的缺点及替代方案（触发器或预写入日志）
- https://github.com/gibbok/typescript-book ：一本简洁的ts手册
- https://eblog.fly.dev/ginbad.html ：“gin是一个糟糕的软件库”
- https://www.sophielwang.com/blog/jpeg ：详细介绍jpeg图片格式压缩的原理
- https://medium.com/@m0t9_/lets-add-a-conditional-expression-to-go-language-3eec7783388e ：
  通过为go语言添加`?:`条件表达式，介绍go语言编译器（词法分析、类型检查等环节）的工作过程和编辑方法
- https://github.blog/engineering/infrastructure/how-github-uses-ebpf-to-improve-deployment-safety/ ：
  github的源代码也托管在github上，这篇文章介绍了github如何使用ebpf-go解决循环依赖的问题
  （举个例子，github无法访问时，github的开发者将无法下载和上传github的源代码）
- https://blog.iamvedant.in/containers-are-not-magic-namespaces-from-scratch ，从零开始创建一个容器，介绍容器技术原理
- https://anirudhology.com/blog/building-raftly-reproducing-production-failures ，
  实现一个可控的raft共识，通过删除一个环节来模拟和理解对应环节的作用
- https://corrode.dev/learn/migration-guides/go-to-rust/ ，一篇go和rust的对比文章
- https://www.alexedwards.net/blog/go-experiments-explained ，go实验性功能，它们可能默认开启或关闭、正在推进加入正式版本或已搁置
- https://segflow.github.io/post/fast-file-search-go/ ，流处理的优化过程，将基础代码（加载到内存然后按byte读取）的效率一路提高到66倍
- https://zackoverflow.dev/writing/why-does-tsgo-use-so-much-memory ，介绍ts的新编译器tsgo为什么占用这么多内存

## 代码/工具

- github.com/go-drift/drift：使用go语言编写移动端跨平台程序(android/ios)，代码风格和flutter非常像。该项目由个人开发，且处于早期阶段
- github.com/aperturerobotics/goscript：将go代码在ast级转换为ts
- github.com/hajimehoshi/ebiten：go语言2D引擎
- github.com/templui/templui：一个让人可以使用go语言编写网页的库，内置支持tailwindCss
- github.com/markel1974/godoom：go语言3d引擎
- github.com/six-ddc/plow:http基准测试工具
- github.com/gookit/validate:通用验证器和过滤器

## 知识点

- go通常使用首字母大小写控制标识符可见性，但还有一条控制规则：internal文件夹
    - 这条规则并非go语言规范的一部分，而是go工具链强制执行的
    - 标识符：包括常量、变量、类型、函数、方法
    - 规则：假设一个模块有一个internal文件夹，那么：
        - 该模块以外无法访问internal文件夹
        - 该模块内，internal文件夹的父文件夹以外，无法访问internal文件夹
      ```text
      module_1
      - go.mod
      - main.go // &#10006;
      
      module_2
      - go.mod
      - main.go // &#10006;
      - level_1
          - main.go // &#10006;
          - level_2
              - main.go // &#10004;
              - external
                  - external.go // &#10004;
              - internal
                  - internal.go // define: PublicVar
      
      &#10004;：允许使用`internal.PublicVar`
      &#10006;：不允许使用`internal.PublicVar`
      ```
    - 错误信息：在编译时，提示import对应行出错，错误描述为**不允许使用internal包**
      （这里用词是`package`，但实际上与包名无关，而是指文件夹名称）
