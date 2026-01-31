# go关键字（draft）

全部关键字列举，标识`x`表示该关键字没有什么可展开的知识点

| keywords |               |        |           | `x` means no expand |
|----------|---------------|--------|-----------|---------------------|
| break    | default       | func   | interface | select              |  
| case ×   | defer         | go     | map       | struct              |  
| chan     | else ×        | goto   | package × | switch              |  
| const    | fallthrough × | if     | range     | type                |  
| continue | for           | import | return    | var                 |  

## break

- 可以用在for循环、switch/select的case中
- 可以通过带一个参数（标签），退出并跳过标签标识的for/switch/select语句，一个**有效的标签**应满足以下条件：
    - 标签只能出现在for/switch/select语句的上一行（中间可以有注释和空行）
    - 标签紧挨着的语句的语句列表(statement list)包含break语句

```code 
ALL:
    for range [3]struct{}{} {
        for range [3]struct{}{} {
            break ALL
        }
    }

// break ALL退出后，不会再执行该循环，而是从这里继续执行
```

## chan

- 可以用来声明和定义`channel`，`channel`包含读写、只读、只写三种类型，
    - 只读和只写多见于函数输入参数位置，为了控制该函数权限
- 可以使用内置函数`close`关闭`channel`
- 读写不同状态`channel`的结果：

|            | read             | write                           | note                                                                  |
|------------|------------------|---------------------------------|-----------------------------------------------------------------------|
| 仅声明(值为nil) | 阻塞               | 阻塞                              | 整个进程没有活动`goroutine`时，将触发死锁，错误：`all goroutines are asleep - deadlock!` |
| 初始化-无缓冲    | 阻塞               | 阻塞                              |                                                                       |
| 初始化-有缓冲    | 正常返回             | 缓冲满时阻塞                          | 无数据时，见上一行                                                             |
| 已关闭-无缓冲    | 返回`channel`类型的零值 | `panic: send on closed channel` | `_, ok := <- ch`，`ok`为`false`                                         |
| 已关闭-有缓冲    | 正常返回             | `panic: send on closed channel` | 无数据时，见上一行                                                             |

## const

- 可使用iota，iota表示的自然数序列会在每个const语句块中初始化，在多行const中，iota的值仅与常量的位置有关 参考代码：

```code 
const (
    c = 'x'   // 'x', iota = 0 is covered
    a = iota  // 1
    b         // 2
    s = "xxx" // "xxx", iota = 3 is covered
    d         // 4
)

const (
    i = iota // 0, iota re-init
    j        // 1
)
```

## continue

- `continue`语句终止执行最近for循环的后续语句
- `continue`语句可以带标签，标签只能出现在for语句的上一行，且该语句的语句列表包含continue语句，
  带标签的continue语句表示跳过标签表示的for循环的后续语句，执行下一次循环。参考代码：

```code 
func f() {
ContinueLabel:
	for range [3]struct{}{} {
		fmt.Println("a")
		for i := 0; i < 3; i++ {
			fmt.Println(i)
			continue ContinueLabel
			fmt.Println("b")
		}
	}

	fmt.Println("finish")
	
	// output: a 0 a 0 a 0 finish
}

func f() {
	for range [3]struct{}{} {
		fmt.Println("a")
ContinueLabel:
		for i := 0; i < 3; i++ {
			fmt.Println(i)
			continue ContinueLabel // 本例与不带标签的continue相同
			fmt.Println("b")
		}
	}

	fmt.Println("finish")
	
	// output: a 0 1 2 a 0 1 2 a 0 1 2 finish
}
```

## default

可使用于switch/select，与case同级，详情参考对应关键字

## defer

- `defer`函数的参数会被立刻求值
- `defer`函数在`return`语句为返回值赋值之后、函数退出之前执行
    - `defer`函数可能改变当前函数的返回值
- 如果`defer`的函数为`nil`，则执行时panic
- `defer`函数自己的返回值会被丢弃
- `defer`语句会压栈（注册）一个函数：
    - 如果`defer`表达式为一串链式调用，则之前的调用函数会立刻执行，仅最后一个函数会被压栈
    - 后注册的`defer`函数先执行
- 若显式调用`os.Exit()`，则注册的`defer`函数不会执行

参考代码：

```code 
func f() {
	defer print(1)
	defer print(2)

	return
	defer print(3) // 没有注册，所以不会执行
}

// output: 21
```

```code 
type chainCall struct {
	n int
}

func (c *chainCall) set(value int) *chainCall {
	c.n = value
	print("set value: ", value)

	return c
}

func f() {
    cc := &chainCall{}

	a := 1
	b := 2
	defer cc.set(1).set(2).set(a+b) // 这里压栈的函数是：cc.set(3)

	a += 10
	cc.set(a)
	
	// output: 1 2 11 3
}
```

## for

三种形式：

- 单一条件(single condition)，类似其他语言的while循环，形如`for a < b { }`
- for子句(for clause)，类似其他语言的for循环，形如`for i := 0; i < 10; i++ { }`
- range子句(range clause)，与range关键字联合使用，形如`for i, v := range arr { }`

在for语句中声明的变量（在for子句和range子句形式中，使用短变量声明），作用域是for语句体，每次迭代会被重复使用

## func

- 函数可以返回多个值
- 函数输出参数可以命名，然后当做常规变量使用，就像输入参数一样
    - 具名的输出参数，会在函数开始时，被初始化为对应类型的零值
    - 若函数有多个输出参数，要么全部命名、要么全部不命名，不允许只为部分输出参数命名
    - 拥有具名输出参数的函数，`return`可以不带参数

## go

go语句在相同的地址空间内，使用独立的并发线程执行函数调用，形如`go func() { }()`

## goto

从接触编程开始，几乎所有的意见都是：不要用goto  
可一方面，go语言保留了goto关键字，作为仅有的25个关键字之一，在go源码中也有所体现；  
另一方面，go github有issue提议：取消使用goto，重写对应源码。

个人倾向于不使用goto，所以本节仅简单描述goto语句规则，为了能看懂使用goto的代码。

goto语句需要一个标签，然后将控制转移到对应标签的位置，要求标签与goto语句在同一函数内。

不建议使用goto跳过变量声明

不能使用goto跳转到另一个局部代码块中（例如不能在一个for循环外面，跳转到里面）

## if

条件判断

## import

`import mlog "github.com/mats0319/.../log"`

- 导入外部代码包，导入的内容仅在当前文件有效
- 可以通过包名访问包内的导出标识符，例如常量、变量、类型、函数、结构体导出字段和方法
- 导入的包名可以任意修改，有两个特殊的例子：`_`、`.`
    - `_`：仅执行该包的`init`函数，不导入任何标识符
    - `.`：可以让你像使用当前包内标识符一样，使用该包的导出标识符

## interface

接口

## map

map的`key`可以是任意类型，只要定义了判等运算符(`==`/`!=`)

## range

可用于`for...range...`句式，遍历`array`、`slice`、`string`、`map`、`channel`、`整数`类型

| type               | e.g.                    | 1st value | 2nd value |
|--------------------|-------------------------|-----------|-----------|
| array or slice     | [3]int *[3]int [ ]int   | index     | value     |
| string             | string                  | index     | rune      |
| map                | map[int]string          | key       | value     |
| channel            | (chan int) (<-chan int) | value     |           |
| integer            | 5                       | index     |           |
| function, 0 values | func(func() bool)       |           |           |
| function, 1 value  | func(func(V) bool)      | V         |           |
| function, 2 values | func(func(K, V) bool)   | K         | V         |

遍历map：

- 顺序随机，且多次遍历顺序可能不同
- 如果在一次遍历中删除一个尚未执行到的key，则本轮遍历不会执行该条目
- 如果在一次遍历中新增一个k-v，则本轮遍历中可能会执行该条目，也可能不会

遍历channel：不断读channel，直到channel关闭。如果读不到或channel为nil则阻塞

遍历函数：实现类似python生成器的功能

## return

实现约束：如果在局部作用域内，定义了与返回值同名的变量，则该作用域内，不能使用隐式返回。参考代码：

```code 
func f() (err error) {
	{
		err := errors.New("new error")
		if err != nil {
			return // 编译报错：err is shadowed during return
		}
	}

	return
}
```

## select

- select与switch结构上相似，但不支持fallthrough关键字
- select语句要求每一个`case`都是通信操作(communication operations)

select语句执行过程：

- 计算所有case子句中，**符合条件的表达式**，所有副作用(side effects)正常执行；接收语句左侧的赋值语句（包括短变量声明）不会执行。
    - 符合条件的表达式：指发送或接收操作的额外表达式，即channel操作符右侧的表达式，  
      举例：`ch <- <-ch1`/`<- getValue()`中的`<-ch1`/`getValue()`，会在下一步*选择分支*之前计算  
      进一步，如果`<-ch1`阻塞，则整个select阻塞
- 选择分支：
    - 若存在可继续的通信（case），则通过统一的伪随机(uniform pseudo-random)选择一个
    - 若不存在可继续的通信（case），则选择default语句执行；若无default子句，则select语句阻塞，直到有一个case可以继续或死锁（全部goroutine阻塞）
- 执行准备：
    - 执行选中case子句的表达式（case到冒号之间的部分，多是短变量声明或赋值）
    - 若选中default子句，直接进入下一步
- 执行选中case的语句块

以下代码包含**执行读写channel的case表达式，且副作用正常执行**、**随机选择一个case执行**以及
**全部计算case表达式、只执行选中case的赋值**的演示，参考代码：

```code 
func main() {
	var (
		times = 1000
		c1    = make(chan int, times)
		c2    = make(chan int, times)
		c3    = make(chan int, times)
	)

	sideEffects := make([]int, 3)
	selected := make([]int, 3)
	{
		count := 1
		c3 <- count
		for i := 0; i < times; i++ {
			select {
			case c1 <- addOne(sideEffects, 0):
				selected[0]++
			case c2 <- addOne(sideEffects, 1):
				selected[1]++
			case sideEffects[2] = <-c3:
				selected[2]++
				count++
				c3 <- count
			}
		}

		<-c3
	}

	remain := make([]int, 3)
	{
		chs := []chan int{c1, c2, c3}
		for i := 0; i < len(remain); i++ {
			remain[i] = len(chs[i])
		}
	}

	fmt.Printf("side effects  : %v\n", sideEffects)
	fmt.Printf("selected times: %v\n", selected)
	fmt.Printf("remain amount : %v\n", remain)
	
        // side effects  : [1000 1000 336] // mainly two first values
        // selected times: [328 336 336]
        // remain amount : [328 336 0] // mainly two first values, compare with last line
}

func addOne(slice []int, index int) int {
	slice[index]++
	return slice[index]
}
```

## struct

定义结构体

## switch

- 类似C语言的switch：`switch [expression] {}`
- 把一连串`if-else`写成switch：`switch {case [expression]: // do sth}`
    - 要求`expression`值为`bool`类型
    - 第一个表达式值为`true`的`case`会被执行
- 类型转换(type switch)，要求t是any类型或类型参数（泛型）
    ```code 
    // from official doc
    var t interface{}
    t = functionOfSomeType()
    switch t := t.(type) {
        default:
            fmt.Printf("unexpected type %T\n", t)     // %T prints whatever type t has
        case bool:
            fmt.Printf("boolean %t\n", t)             // t has type bool
        case int:
            fmt.Printf("integer %d\n", t)             // t has type int
        case *bool:
            fmt.Printf("pointer to boolean %t\n", *t) // t has type *bool
        case *int:
            fmt.Printf("pointer to integer %d\n", *t) // t has type *int
    }
    ```

## type

定义类型  
新类型不继承底层类型的方法集，类型别名可以，但直接继承的方法集受到package的限制（可导出的方法可以正常调用非导出的方法），参考代码：

```code 
type S struct {
}

func (s *S) ExportedFunc() {
   s.nonExportedFunc()
}

func (s *S) nonExportedFunc() {
}

type sNew S
type sAlias = S // 类型别名声明(alias declarations)
type sAlias2 = another_package.S // 假设another_package包，定义有同样的结构体S与方法集

func main() {
    var sn = &sNew{}
    var sl = &sAlias{}
    var sl2 = &sAlias2{}
    
    sn.ExportedFunc()       // wrong
    sn.nonExportedFunc()    // wrong
    sl.ExportedFunc()       // right
    sl.nonExportedFunc()    // right
    sl2.ExportedFunc()      // right，即使该可导出方法中调用了不可导出的方法，也能正确执行
    sl2.nonExportedFunc()   // wrong
}
```

## var

定义变量
