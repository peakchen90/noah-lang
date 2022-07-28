package main

import (
	"encoding/json"
	"fmt"
	"github.com/peakchen90/noah-lang/internal/lexer"
	"github.com/peakchen90/noah-lang/internal/parser"
	"os"
)

func main() {
	code, _ := os.ReadFile("example.noah")
	node := *parser.NewParser(string(code))

	jsonStr, _ := json.MarshalIndent(node, "", "  ")
	fmt.Println(string(jsonStr))

	fmt.Println(string(lexer.TTColon))

	//switch v := node.Node.(type) {
	//case *ast.Program:
	//	fmt.Println(v.Body)
	//default:
	//	fmt.Println("none")
	//}
}
