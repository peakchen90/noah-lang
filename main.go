package main

import (
	"encoding/json"
	"fmt"
	"github.com/peakchen90/hera-lang/internal/parser"
)

func main() {
	node := *parser.NewParser(`
import "axe.aa/fwa" as ab
`)

	jsonStr, _ := json.MarshalIndent(node, "", "  ")
	fmt.Println(string(jsonStr))

	//switch v := node.Node.(type) {
	//case *ast.Program:
	//	fmt.Println(v.Body)
	//default:
	//	fmt.Println("none")
	//}
}
