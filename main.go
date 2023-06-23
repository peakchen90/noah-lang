package main

import (
	"encoding/json"
	"fmt"
	"github.com/peakchen90/noah-lang/internal/compiler"
)

func main() {
	inst := compiler.NewCompiler("").Compile()

	jsonStr, _ := json.MarshalIndent(inst.Main.Ast, "", "  ")
	fmt.Println(string(jsonStr))
}
