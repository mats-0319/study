# Android 开发

## Android开发技术选型

常见技术栈：

- kotlin + jetpack compose
- kotlin + XML
- flutter
- React Native
- Unity

主要就是在`kotlin+jetpack compose`和`flutter`中间选择，列举我的情况：

- 编写简单的Android app，或许会涉及少许android naive（例如webview）
- 基本不考虑ios平台、非移动端平台(web、windows、linux等)

这样来看，`flutter`的优势就不大了，这两种技术路线都是google维护的，而`flutter`的劣势是调用android naive的能力：

- flutter调用到android native的时候还是需要写kotlin，例如我要用到的webview
    - 我没亲手写kotlin，是因为引用的库做了这部分工作
    - 我个人对这一点还是挺在意的，flutter这回是有对应的包了，下回呢？但凡遇到一个flutter没有的包，我还是得去学kotlin
        - 例如google play、AWS等第三方SDK，通常只有kotlin版本而没有dart版本

所以现在让我写android app，我可能会选择`kotlin+jetpack compose`，因为它全栈都是kotlin语言。