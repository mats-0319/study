# dart 学习笔记

> [language](https://dart.dev/language)

## 安装与使用 install and use

dart基本上是和flutter绑定的，安装flutter时会自动安装dart

dart命令：

```cmd
dart create cli # 创建dart项目
dart run
```

## dart语言

代码格式：

- 一条语句之后必须要有分号(`;`)

### 变量 variable

```dart
void main() {
  var name = "mario"; // dart变量在使用前必须初始化（可空变量可以被默认初始化为null）
  String name = "mario"; // 显式声明变量类型
  String? name;
  late String name; // 将一个变量标记为延迟初始化，可以对冲dart的未初始化判断
  final String name = "mario"; // final/const 表示变量/编译时常量不可修改
  var _ = 1;
}
```

### 运算符 operator

```dart
void main() {
  assert(5 / 2 == 2.5); // 数学除法，结果是`double`类型
  assert(5 ~/ 2 = 2); // 其他编程语言中的整数除法

  (v as Struct).Method(); // 强制类型转换，转换失败会抛出异常
  if (v is Struct) {}; // 类型检查

  // 条件表达式
  condition ? expr1 : expr2;
  expr1 ?? expr2; // `expr1 == null ? expr1 : expr2`的简写，`expr1`不空时取`expr1`，否则取`expr2`

  // 级联符号：以下两种代码等价
  var paint = Paint() // 注意这里没有分号
    ..color = Colors.black; // 如果对象可空，则应在第一个级联项前使用'?..'，相当于`paint?.color`

  var paint = Paint();
  paint.color = Colors.black;
}
```

### 类型 type

#### 内置类型

- 数字类型：`int`/`double`
    - 可以定义成`num`类型，该类型的变量可以同时具有整型和浮点型的值
- 字符串类型：`String`
    - UTF-16
- 布尔类型：`bool`
- 记录类型：`(v1,v2)`，可以按照简易结构体来理解
    - 可以用序号或命名访问：`record.$1`/`record.a`
    - 为记录指定名称：
      `typedef ButtonItem = ({String label, Icon icon, void Function() ? onPressed});`
- 函数类型：`Function`
- 集合类型：
    - 数组：`List`，可以用下标访问
    - 集合：`Set`，无重复元素
    - 映射：`Map`
        - 尝试查询映射中不存在的`key`会得到`null`
- 字符类型：`Runes`
- 符号类型：`Symbol`
- null类型：`null`

```dart 
void main() {
  // 数字类型和字符串类型相互转换
  int n = int.parse("1");
  String s = 1.toString();

  // 字符串表达式
  // '{}'用于包裹表达式，如果表达式是标识符，则可以省略'{}'
  // dart会自动调用对象的'toString()'方法
  String s1 = "value: ${s}";

  // 集合初始化
  List<int> li = [1, 2, 3];
  Set<int> si = {1, 2, 3};
  Map<String, int> msi = {"key1": 1, "key2": 2};
}
```

#### 特殊类型

- Object：除了`Null`之外的所有dart类的超类
- Enum：枚举的超类
- Future/Stream：用于异步编程
- Iterable：用于循环和同步生成器函数
- Never：表示表达式永远无法成功求值，常用于总是抛出异常的函数(?)
- dynamic：禁用静态类型检查，通常应该使用`Object`/`Object?`代替
- void：表示该值永远不会被使用，通常用作函数的返回值类型

### 模式 pattern

todo

### 控制流 control

- 循环：主要记住for
- 分支：主要记住if、switch
- 错误处理：throw / try-catch-finally

### 函数 function

```dart
bool isEqual(int a, int b) {
  return a == b;
}

bool isEqual(int a, int b) => a == b; // 如果函数体只有一条表达式，则可以使用箭头('=>')简写

// 命名参数，使用'{}'包裹任意位置直到输入参数列表末尾，默认参数和可空参数不能在常规参数中出现
bool isEqual(int p, {int a = 1, int? b}) => a == b;
// isEqual(1,a:2,b:3); // 调用时需要指定命名参数的函数名

// 可选参数，这种写法更像是ts的可选参数，调用时可以给或不给该参数的值
bool isEqual(int a, [int b = 2]) => a == b;
```

### 类 class

### 构造函数 constructor

```dart
/* 6种构造函数 */

class Point {
  double x;
  double y;

  // 默认构造函数：未声明构造函数时，dart会生成一个不带参数和名称的生成式构造函数
  Point(this.x, this.y); // 生成式构造函数
  Point.origin() : x=0; // 命名构造函数
}

class ImmutablePoint {
  final double x, y;

  const ImmutablePoint(this.x, this.y); // 常量构造函数
}

class SingletonPoint {
  static final SingletonPoint _point = SingletonPoint(); // 使用默认构造函数

  // 工厂构造函数，每次调用返回相同实例（也可以使用不同的，取决于写法）
  factory SingletonPoint() {
    return _point;
  }
}
```

### 枚举 enum

```dart
enum Color { red, green, blue }
```
