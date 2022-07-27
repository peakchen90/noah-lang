# 内置类型结构体（隐式）

## [Number]

```noah
type :[Number] {
    toStr()
    clone()
}
```

## [Bool]

```noah
type :[Bool] {
    toStr()
    clone()
}
```

## [Array]

```noah
type :[Array] {
    toStr()
    len() -> num
    clone() -> []T
}
```

## [VectorArray]

```noah
type :[VectorArray] extends [Array] {
    push(..item: [..]T)
    pop() -> T
    unshift(..item: [..]T)
    shift() -> T
    splice(index: num, len: num, ..item: [..]T) -> [..]T
    slice(start: num, end: num) -> [..]T
}
```

## [String]

```noah
type :[String] {
    len() -> num
    clone() -> []T
    split(str) -> [..]str
    toChars -> [..]T
    toUpperCase() -> str
    toLowerCase() -> str
    trim() -> str
    indexOf(ch: char) -> num 
    slice(start: num, end: num) -> str
}
```