package parser

import (
	"os"
	"testing"
)

var parserFixtures = [...]string{
	`
let a = 1
const b: str = "a"
let c []num = [1, 2]
let d [2]str
let e [..]num
let f [..]num = [1]
	`,
	`
type A {
	a: num
	b: str
}
type B extends A {
	c: bool
}
type C { A, B }

let a: A = null
let c: C = C.B 
	`,
	`
interface A {
	a: num
	b(n: num) -> str
}
interface B extends A {
	c()
}
type C: A,B {
	d: str
}
type D: B {}
	`,
	`
fn a() {}

fn b(a: str, ..b: [..]num) -> bool {
	return true
}
	`,
	`
type A {
	a: num
}

fn A a() {
	self.a = 3
}
	`,
	`
type A { a: num }

let a: A = { a: 1 }
let b: A = { 2 }
let c = A{ a: 3 }
let d = A{ 4 }
	`,
	`
if 1 + 1 == 2 {
	let a = 1
} else if true {
	let a = 2
} else {
	let a = 3
} 
	`,
	`
let arr: [3]num
for item, index: arr {
    println(item, index)
}
for let i = 0; i < arr.len(); i = i + 1 {
    println(arr[i], i)
}
let i = 0;
for i < arr.len() {
    i = i + 1
}
for {
    break
    continue
}
	`,
}

func TestParser(t *testing.T) {
	files, err := os.ReadDir("testdata")
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		file.IsDir()
	}

	for _, fixture := range parserFixtures {
		NewParser(fixture)
	}
}
