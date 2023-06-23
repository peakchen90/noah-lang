package parser

import (
	"os"
	"testing"
)

var parserFixtures = [...]string{
	`
let a = 1
const b: string = "a"
let c: []number = [1, 2]
let d: [2]string
let e: []number
let f: []number = [1]
	`,
	`
type A {
	a: number
	b: string
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
	a: number
	b(n: number) -> string
}
interface B extends A {
	c()
}
type C: A,B {
	d: string
}
type D: B {}
	`,
	`
fn a() {}

fn b(a: string, ...b: []number) -> bool {
	return true
}
	`,
	`
type A {
	a: number
}

fn A a() {
	self.a = 3
}
	`,
	`
type A { a: number }

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
let arr: [3]number
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
		if file.IsDir() {
			continue
		}

		code, err := os.ReadFile("testdata/" + file.Name())
		if err != nil {
			panic(err)
		}

		NewParser(string(code)).Parse()

	}

	//for _, fixture := range parserFixtures {
	//	NewParser(fixture)
	//}
}
