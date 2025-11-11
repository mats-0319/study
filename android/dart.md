# dart 官方文档阅读笔记

## 变量

使用`var name="mario"`或`String name="mario"`定义变量，有类型推断也可以指定

### 空安全

dart不允许访问**一个可能为空的表达式的属性或方法**（除了toString()方法和hashcode属性）

定义变量的时候就建议初始化（可空变量可以自动初始化为null），使用变量之前必须赋过值

（所有可空类型的默认值都是null）

late关键字：延迟初始化，定义在类型前面，可以在第一次使用变量的时候初始化

final/const关键字：定义常量，final定义对象的字段可以修改，const定义的不行

## 运算符

整数除法会产生小数结果，想要得到整数结果应使用`~/`

强制类型转换：`(userIns as User).username = "mario"`

类型判断：`if (userIns is User) {}`

## 类型

### 内置类型

- 数字类型（int，double）
  - String -> int: var one = int.parse("1");
  - int -> String: var str = 1.toString();
  - double -> String: var str = 3.14159.toStringAsFixed(2); // 3.14
  - 数字类型中间可以添加下划线`_`来方便阅读，例如`const n = 1_000_000`
- 字符串（String）
- 布尔类型（bool）
- Records（(value1,value2)）：一串值放在一起
  - 形如：`(String, int) recordIns = ("abc", 1)`
  - 可以为每个值命名：`({int a, int b}) recordIns = (a: 1, b: 2)`
  - 如果上一条中没有`{}`，则每个值的命名将被视为文档说明/注释，不具有实际意义（初始化时无法使用，常用于函数返回值）
  - 两个Record参数类型和顺序全部相同时视为相同类型、可以相互赋值；任意具名Record不与其他Record属于相同类型
  - 访问Record：`var recordIns = ("first", a: 1, b: true, "last")`使用`recordIns.a`/`recordIns.$1`访问
- 函数（Function）
- 数组（List）：形如`List<int> list = [1,2,3]`
- 集合（Set）：形如`Set<String> set = {"a", "b", "c"}`，无序且不重复的数组，使用`.add()`方法新增
- 映射（Map）：
  - 形如`Map<String, String> m = Map<String, String>()`，使用`m["key"]="value"`新增
  - 尝试获取不存在的key，会获得一个null
- 字符（Runes，常使用characters接口代替）
- 符号（Symbol）
- 空（null）

Object：dart所有的类型都是对象，所以它也是所有类型的超类（除了null）

Enum：枚举类型的超类，形如`enum Color { red, green, blue }`

Future/Stream：用于异步编程

Iterable：用于for-in循环和同步生成器函数

typedef: 类型别名，形如`typedef IntList = List<int>； IntList ins = [1,2,3]`

## 控制流

### 异常处理

dart可以throw任何内容作为异常，例如`throw "a string"`

捕获：
``` dart
try {
  breedMoreLlamas();
} on OutOfLlamasException {
  // A specific exception
  buyMoreLlamas();
} on Exception catch (e) {
  // Anything else that is an exception
  print('Unknown exception: $e');
} catch (e, s) {
  // No specified type, handles all
  print('Something really unknown: $e');
  print('Stack trace:\n $s');
} finally {
  // Always clean up, even if an exception is thrown.
  cleanLlamaStalls();
}
```

## 函数

dart中的函数也是一个对象

## 使用库

`import "dart:js_intgerop"` / `import "package:test/test.dart" as t`

下划线（`_`）开头的内容仅库内部可见

## 类和对象

### implements & extends

一个类可以实现多个接口，但是只能继承一个类

使用`@override`重写父类的方法
