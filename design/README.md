# 语言设计

## 类型

**内置类型**:

- `number`: 数字类型，64 位浮点型（默认值: `0`）
- `byte`: 无符号 8 位（默认值: `0`）
- `char`: UTF-8 字符类型，32位（默认值: `0`）
- `string`: 字符串类型（默认值: `""`）
- `bool`: 布尔类型（默认值: `false`）
- `[n]T` : 数组类型，如：`[3]string`、`[]number` 等（默认值: `null`）
- `[]T` : 可变长数组类型，如：`[]string`、`[]number` 等（默认值: `null`）
- `any`: 动态类型（默认值: `null`）

**自定义类型**:

- `type T T2` : 声明自定义类型
- `interface T {}`: 接口
- `struct T {a: string, b: number}` : 结构体类型（默认值: `null`）
- `enum T {A, B}` : 枚举类型（默认值: `null`）

```noah
// 定义数字类型的别名
type TypeNum number

// 定义结构体
struct Person {
    name: string
    age: number
}

// 定义枚举类型
enum Color {
    Red,
    Green
}

// 定义接口
interface Man {
    fn say(a: number) -> string
}

// 定义一个 `Student` 结构体，继承 `Person` 的属性和实现的方法
struct Student <- Person {
    grade: number
}
```

## 变量

**基础类型**：

```noah
// 声明一个字符串类型的变量
let hello: string = "hello world"

// 声明常量
const PI: number = 3.14159

// 声明变量，可省略类型，系统会自动推断为布尔类型
let flag = true
```

**结构体**：

```noah
let s1: Person; // null

let s2 = Person{ name: "noah" } // { name: "noah", age: 0 }

let s3: Person = { name: "noah" } // { name: "noah", age: 0 }

fn main() {
    // 修改 `age` 属性值
    s3.age = 22
}
```

**枚举类型**：

```noah
let e1: Color; // null

let e2 = Color.Red // Color.Red
```

**数组类型**：

```noah
let arr1: [3]number; // [0, 0, 0]

// 有初始值时数组长度可省略
let arr2: []string = ["a", "b", "c"] // ["a", "b", "c"]

let arr3: [3]number = [1, 2] // [1, 2, 0]

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

fn main() {
    // 修改第 3 个元素的值
    arr4[2] = { name: "noah", age: 28 }
    
    // 修改第 3 个元素的值的 `age` 属性
    arr4[2].age = 18
}
```

## 可变长数组

```noah
fn main() {
    let arr: []number = [1]
    
    arr.push(2) // arr == [1, 2]
    arr.unshift(3) // arr == [3, 1, 2]
}
```

## 函数

除 **定义变量**、**定义类型**、**定义函数** 外，其他语句必须放在函数里执行，`main` 函数会程序的入口。

```noah
// 声明一个名为 `foo` 的函数，入参 `name` 的类型是字符串，返回值是布尔类型
fn foo(name: string) -> bool {
    return true
}

// 为 `Person` 结构体实现一个名为 `foo` 的方法，方法内部可以使用 `self` 关键字指向结构体的实例 
fn Person foo(...name: []string) -> string {
    return self.name
}

// 函数调用
fn main() {
    foo("hello world")
}
```

**剩余参数**：

```noah
fn add(...nums: []number) -> number {
    let sum = 0
    for v: nums {
        sum += v
    }
    return sum
}

fn main() {
    add(1, 2) // 3
    add(10, 20, 30) // 60
}
```

## 内存引用

在函数参数传递及赋值语句中，除**数组、可变长数组、结构体、函数引用**是传递内存地址引用外，其他类型都是传递值的拷贝

## 逻辑控制

**函数返回**：

```noah
fn main() {
    return // 空返回
    return "abc" // 返回字符串
}
```

**条件控制**：

```noah
fn main() {
    if expr1 {
        // do something
    } else if expr2 {
        // do something
    } else {
        // do something
    }
}
```

**循环**：

```noah
fn main() {
    let arr: []number = [1, 2, 3]
    
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
    label: for {
        break
        continue
    }
}
```

## 多态

```noah
pub interface Person {
    say() -> string
}

pub struct Man {
    name: string
}

pub struct Woman {
    nick: string
}

impl (Person) Man {
    pub fn say() {
        return "Man: " + self.name
    }
}

impl (Person) Woman {
    pub fn say() {
        return "Woman: " + self.nick
    }
}
```

**`strcut` 继承**

```noah
struct OldMan <- Man {
    age: number
}

struct Student <- Man, Woman {
    scores: []number
}
```

**使用结构体**：

```noah
// 作为函数参数
fn hello(p: Person) {
    p.say() // return string
}

fn main() {
    let s = Student{}
    s.scores
    s.name
    s.nick
    
    // 作为变量类型
    let p: Persion
    p = Man{}
    p = Woman{}
    hello(p)
}
```

## 动态类型

```noah
fn hello(value: any) {
    if value is string {
    
    } else if value is number {
    
    } 
}
```

**类型转换**

```noah
let a: number = 'a' as number
```

## 模块

### 定义模块

同一个文件的变量、类型都同属一个模块。模块内的变量、类型可以互相访问，但是对外部模块来说这些变量及类型默认都是私有的，可通过 `pub`
向外部暴露

`lib/foo.noah`

```noah
var abc = 123 // private

pub const PI = 3.14159 // public

// public struct
pub struct Person {
    name: string
}
```

`lib/bar.noah`

```noah
use "foo"
pub use "foo" as foo2 // re-export

let n1 = foo.PI // 3.14159 (from `lib/foo.noah`)
type P foo2.Person

pub fn say() -> string {
    return "Hello World"
}
```

`main.noah`

```noah
use "lib/bar"

let n = bar.foo2.PI // 3.14159 (from `lib/foo.noah`)

fn main() {
    let s = bar.say() // (from `lib/bar.noah`)    
}
```

**三方模块(模块名以 `moduleName:` 开始)**:

```noah
use "std:numbers" // 导入标准库模块
use "third:lib/foo" // 导入三方库模块
```

## 其他

[内置类型接口（隐式）](./implicit-interface.md)