# 学习笔记

## 数学知识复习

指数：

- X<sup>A</sup>X<sup>B</sup> = X<sup>A+B</sup>
- X<sup>A</sup> / X<sup>B</sup> = X<sup>A-B</sup>
- (X<sup>A</sup>)<sup>B</sup> = X<sup>AB</sup>

对数：

- 假设X<sup>A</sup> = B，则有log<sub>X</sub>B = A
- log<sub>A</sub>B = log<sub>C</sub>B / log<sub>C</sub>A (C > 0)
- log AB = log A + log B （默认为2底对数）
- log A/B = log A - log B
- log (A<sup>B</sup>) = B log A

级数：累加求和

## 抽象数据类型 abstract data type ADT

抽象数据类型包括数据结构本身和它的操作，就像整型、字符串类型等普通数据类型一样。

[理解]go语言规范中对类型的定义是：一个取值范围、一组运算符、一个方法集

- 取值范围很好理解，你不能把整数值赋值给字符串类型的变量
- 运算符指该类型可以使用哪些运算符，例如能否使用判等运算符 (`==`)、移位运算符 (`>>`)
    - 这一点在go中体现的不明显，在其他允许运算符重载的语言中更加清晰，例如rust和kotlin
- 方法集指定义在该类型上的方法的集合
