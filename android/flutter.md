# flutter

## 打包和安装

`flutter build appbundle`
`flutter build apk --split-per-abi`
`flutter install --use-application-binary=build/app/outputs/flutter-apk/app-arm64-v8a-release.apk`

查看手机cpu架构：（需要手机开启usb调试）
`adb shell getprop | grep cpu`

注：不要像flutter文档中提到的，直接使用`flutter install`，这会导致下载成flutter默认的演示项目，
实际安装时应指明下载的软件包（或者在打包的时候不区分不同的cpu架构、把几种包打在一起）

## codelab

flutter文档中提供了一个codelab，完成之后对flutter编程有了一些基本的了解：

- 在使用过程中大概了解了哪些Widget有长宽、padding/margin等字段可以设置
- 用flutter写样式整体上比起使用html+css更容易

## widget的状态

有/无状态的widget：（stateful/stateless）指widget是否需要跟随指定条件而变化（例如用户输入）  
无状态的widget不会变化；有状态的widget需要关联其他对象来构建：

```flutter
class MyCounter extends StatefulWidget {
  const MyCounter({super.key});

  @override
  State<MyCounter> createState() => _MyCounterState();
}

class _MyCounterState extends State<MyCounter> {
  int count = 0;

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        Text('Count: $count'),
        TextButton(
          onPressed: () {
            setState(() {
              count++;
            });
          },
          child: Text('Increment'),
        )
      ],
    );
  }
}
```

## widget状态管理

父子widget间共享状态：创建子widget时传入变量（使用widget构造函数）

ChangeNotifier:

```flutter
class CounterNotifier extends ChangeNotifier {
  int _count = 0;
  int get count => _count;

  void increment() {
    _count++;
    notifyListeners();
  }
}

// 调用
Column(
  children: [
    ListenableBuilder(
      listenable: counterNotifier,
      builder: (context, child) {
        return Text('counter: ${counterNotifier.count}');
      },
    ),
    TextButton(
      child: Text('Increment'),
      onPressed: () {
        counterNotifier.increment();
      },
    ),
  ],
)
```

MVVM：数据与UI分开

```flutter
ListenableBuilder(
  listenable: viewModel,
  builder: (context, child) {
```
