# 语言设计

## 类型

**内置类型**:

- `num`: 数字类型，8 字节浮点型（默认值: `0`）
- `char`: 字符类型，1 字节（默认值: `0`）
- `str`: 字符串类型（默认值: `""`）
- `bool`: 布尔类型（默认值: `false`）
- `[n]T` : 数组类型，如：`[3]str`、`[]num` 等（默认值: `null`）
- `[..]T` : 可变长数组类型，如：`[..]str`、`[..]num` 等（默认值: `null`）
- `any`: 动态类型 `interface any {}`（默认值: `null`）

**自定义类型**:
- `interface T {}`: 接口
- `type T num` : 类型别名
- `type T {a: str, b: num}` : 结构体类型（默认值: `null`）
- `type T {A, B}` : 枚举类型（默认值: `null`）

```hera
// 定义数字类型的别名
type TypeNum num

// 定义结构体
type Person {
    name: str
    age: num
}

// 定义枚举类型
type Color {
    Red,
    Green
}

// 定义一个 `Student` 结构体，继承 `Person` 的属性和实现的方法
type Student extends Person {
    grade: num
}
```

## 变量

**基础类型**：

```hera
// 声明一个字符串类型的变量
let hello: str = "hello world"

// 声明常量
const PI: num = 3.14159

// 声明变量，可省略类型，系统会自动推断为布尔类型
let flag = true
```

**结构体**：

```hera
let s1: Person; // null

let s2 = Person{ name: "Hera" } // { name: "Hera", age: 0 }

let s3: Person = { name: "Hera" } // { name: "Hera", age: 0 }

// 修改 `age` 属性值
s3.age = 22
```

**枚举类型**：

```hera
let e1: Color; // null

let e2 = Color.Red // Color.Red
```

**数组类型**：

```hera
let arr1: [3]num; // [0, 0, 0]

// 有初始值时数组长度可省略
let arr2: []str = ["a", "b", "c"] // ["a", "b", "c"]

let arr3: [3]num = [1, 2] // [1, 2, 0]

// 结构体数组
let arr4: [3]Person = [
    { name: "Alice", age: 18 },
    { name: "Bob"}
]
//  [
//      { name: "Alice", age: 18 },
//      { name: "Bob", age: 0 },
//      null
//  ]

// 修改第 3 个元素的值
arr4[2] = { name: "Hera", age: 28 }

// 修改第 3 个元素的值的 `age` 属性
arr4[2].age = 18
```

## 可变长数组

```hera
let arr: [..]num = [1]

arr.push(2) // arr == [1, 2]
arr.unshift(3) // arr == [3, 1, 2]
```

## 函数

```hera
// 声明一个名为 `foo` 的函数，入参 `name` 的类型是字符串，返回值是布尔类型
fn foo(name: str) -> bool {
    return true
}

// 为 `Person` 结构体实现一个名为 `foo` 的方法，方法内部可以使用 `self` 关键字指向结构体的实例 
fn Person foo(..name: [..]str) -> str {
    return self.name
}

// 函数调用
fn main() {
    foo("hello world")
}
```

**剩余参数**：

```hera
fn add(..nums: [..]num) -> num {
    let sum = 0
    for n: nums {
        sum = sum + n
    }
    return sum
}

fn main() {
    add(1, 2) // 3
    add(10, 20, 30) // 60
}
```

## 内存引用

在函数参数传递及赋值语句中，除**字符串、数组、可变长数组、结构体**是传递内存地址引用外，其他类型都是传递值的拷贝

## 逻辑控制

**函数返回**：

```hera
return // 空返回
return "abc" // 返回字符串
```

**条件控制**：

```hera
if expr1 {
    // do something
} else if expr2 {
    // do something
} else {
    // do something
}
```

**循环**：

```hera
let arr: []num = [1, 2, 3]

// 遍历数组的元素及索引
for item, index: arr {
    println(item, index)
}

// 包含初始值声明、条件、更新语句的循环
for let i = 0; i < arr.len(); i = i + 1 {
    println(arr[i], i)
}

// 仅有条件的循环
let i = 0;
for i < arr.len() {
    i = i + 1
}

// 无限循环，可通过 `break` 或 `continue` 跳出循环
for {
    break
    continue
}
```

## 多态

```hera
interface Man {
    name: str
    say() -> str
}

// 结构体 `Person` 实现 `Man` 接口
type Person: Man {}

// 必须实现 `say()` 方法
fn Person say() -> str {
    return "Person: " + self.name
}

type Student: Man {}

// 必须实现 `say()` 方法
fn Student say() -> str {
    return "Student: " + self.name
}
```

**`interface` 继承**

```hera
interface Woman extend Man {
    eat(n: num)
}

type Alice: Woman {
    weight: num
}

// 必须实现 `say()` 方法
fn Alice say() -> str {
    return "Woman: " + self.name
}

// 必须实现 `eat()` 方法
fn Alice eat(n: num) {
    self.weight = self.weight + n
}
```

**使用多态**：

```hera
// 作为函数参数
fn hello(m: Man) {
    m.say()
}

let w = Woman{}
hello(w)

// 作为变量类型
let man: Man
man = Person{}
man = Woman{}
```

**动态类型**

```hera
fn hello(value: any) {
    if type(value) == str {
    
    } else if type(value) == num {
    
    } 
}
```

## 模块

### 定义模块

同一个文件夹的所有变量、类型都同属一个模块。模块内的变量、类型可以互相访问，但是对外部模块来说这些变量及类型默认都是私有的，可通过 `pub` 向外部暴露

`foo.hera`

```hera
var abc = 123 // private

pub const PI = 3.14159 // public

// public struct
pub type Person {
    name: str
}
```

`sub/foo.hera`

```hera
import ".." as entry

let n = entry.abc // 123 (from `foo.hera`)

pub say() -> str {
    return "Hello World"
}
```

`main.hera`

```hera
import "sub" as sub

let n = abc // 123 (from `foo.hera`)

fn main() {
    let s = sub.say()
}
```

**全局模块(模块名以 `mod:` 开始)**:

```hera
// 导入标准库模块
import "mod:std/numbers" as numbers
numbers.toNum("1.2") // 1.2

// 引用第三方模块
import "mod:github.com/foo/bar" as third
third.xxx
```

## 其他

[内置类型接口（隐式）](./implicit-interface.md)