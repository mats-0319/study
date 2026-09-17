# Android 开发

## Android开发技术选型

android开发主流方案：

- flutter
- kotlin + jetpack compose
- react native

暂时不考虑react，react jsi我们都没有了解过，加上项目要求快速迭代。

flutter的优势是简单，体现在dart语言比kotlin简单、相同项目代码量更少等方面；
劣势是涉及原生功能时需要通过methodChannel调用kotlin、很多sdk没有dart语言版本（例如google play、AWS）。

## flutter

- [安装flutter、android开发组件](https://docs.flutter.dev/install/quick)
- [flutter学习](https://docs.flutter.dev/learn/pathway/tutorial)
- [flutter实战](https://book.flutterchina.club/chapter1/install_flutter.html)
- [codelabs](https://docs.flutter.dev/reference/learning-resources)

codelabs带有代码的学习资料举例：

- [WebView](https://codelabs.developers.google.com/codelabs/flutter-webview?hl=zh-cn#0)

dart sdk的位置：`~/flutter/bin/cache/dart-sdk`

项目启动依赖很多环境，例如：`flutter`/`java`/`gradle`/`AGP`/`AGP plugin`/`kotlin`，它们当中只要有一个不兼容，
程序就无法运行。所以不建议单独升级某一个环境，而是交给flutter统一管理，想要升级也是升级flutter、创建新项目然后合并。

### 常用命令

设置国内镜像：`android/gradle/wrapper/gradle-wrapper.properties`文件，
`distributionUrl`的值修改为：`https\://mirrors.cloud.tencent.com/gradle/gradle-9.3.1-all.zip`

- `flutter create --platforms=android`
- `flutter clean`
- `flutter pub get` 下载依赖
- `flutter run -v`
- `flutter build apk --split-per-abi` 打包
- `flutter install --use-application-binary=build/app/outputs/flutter-apk/app-arm64-v8a-release.apk`
    - 这里要具体指到对应架构的包，不然flutter会安装成示例代码；或者打包时不区分架构也行

ADB:

- `adb shell getprop | grep cpu` 查看手机cpu架构（需要手机开启usb调试）
- `adb shell ls -l /sdcard/Download/` 列举模拟器指定目录的文件
- `adb push ./message.txt /sdcard/Download/` 向模拟器发送文件
- `adb pull /sdcard/Download/ .` 从模拟器下载文件
