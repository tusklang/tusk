package operations

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/value"
	"github.com/tusklang/tusk/ast"
	"github.com/tusklang/tusk/initialize"
)

//package to initialize all the operations

func InitOperations(compiler *ast.Compiler, prog *initialize.Program) {
	compiler.Operations = make(map[string]func(value.Value, value.Value, *ir.Block) value.Value)
	initIntOps(compiler)
}
