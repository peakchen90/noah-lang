package codegen

import (
	"fmt"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
)

// TODO

func LLVMTest() {
	module := ir.NewModule()

	module.NewGlobalDef("aaa", constant.NewInt(types.I32, 1000000))

	fmt.Println(module.String())
}
