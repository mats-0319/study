# Go Code Style

记录我的go代码风格

## 判断error，是否将函数调用写在if行内

``` go 
// style 1
err := f()
if err != nil {}

// style 2
if err := f(); err != nil {}
```

- 风格2可以限定`err`的作用域，让代码更清晰
- 风格1在多返回值等场景下更简洁，举个例子：
  ``` go
  amount, users, err := db.listUser()
  if err != nil {}
  ```
  在上例中，风格2需要提前定义全部左值，并且`err`的作用域也没有限制住

综上所述，多返回值场景下，使用风格1；简单单一返回值场景下，使用风格2
- 简单：指返回值为error、ok等不包含业务数据的类型，且返回值最多判断一下、后续不会再使用的场景
