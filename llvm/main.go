package main

import (
	"fmt"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
)

func main() {
	// Create a new LLVM IR module.
	// i32 := types.I32
	m := ir.NewModule()

	x := m.NewFunc("main", types.Void)
	block := x.NewBlock("")
	printf := ir.NewFunc("printf", types.Void)
	printf.Sig.Variadic = true
	m.Funcs = append(m.Funcs, printf)
	def := m.NewGlobalDef("$const_str", constant.NewCharArrayFromString("Hello World!!"))
	block.NewCall(printf, def)
	block.NewRet(nil)
	// Print the LLVM IR assembly of the module.
	fmt.Println(m)
}
