# 内置类型结构体（隐式）

## [Number]

```noah
interface [Number] {
    toStr()
    clone()
}
```

## [Bool]

```noah
interface [Bool] {
    toStr()
    clone()
}
```

## [Array]

```noah
interface [Array] {
    toStr()
    len() -> num
    clone() -> []T
}
```

## [VectorArray]

```noah
interface [VectorArray] extends [Array] {
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
interface [String] {
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