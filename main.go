package main

import "github.com/peakchen90/hera-lang/compiler/ast"

func main() {
	ast.NewParser(`
"\76g 0a\tb\n\a\"c\\"
`)
}
