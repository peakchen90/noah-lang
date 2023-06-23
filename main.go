package main

import (
	"encoding/json"
	"fmt"
	"github.com/peakchen90/noah-lang/internal/compiler"
	"github.com/peakchen90/noah-lang/internal/parser"
	"os"
)

func main() {
	code, _ := os.ReadFile("example.noah")
	node := parser.NewParser(string(code))

	c := compiler.NewCompiler(string(code), "")
	c.Compile()
	//c.mainModule.preCompile()

	firstNode := node.Body[0].Node

	//expr, ok := firstNode.(*ast.ExprStmt)
	//if ok {
	//	//a := c.inferKind(expr.Expression)
	//	//fmt.Println(a)
	//}

	jsonStr, _ := json.MarshalIndent(node, "", "  ")
	fmt.Println(string(jsonStr))

	fmt.Println(c, firstNode)

	//switch v := node.Node.(type) {
	//case *ast.Program:
	//	fmt.Println(v.Body)
	//default:
	//	fmt.Println("none")
	//}
}
