package compiler

import (
	"github.com/llir/llvm/ir"
	"github.com/tusklang/tusk/ast"
	"github.com/tusklang/tusk/data"
	"github.com/tusklang/tusk/operations"
)

func initDefaultOps(compiler *ast.Compiler) {

	compiler.OperationStore = operations.NewOperationStore()

	compiler.OperationStore.NewOperation("+", "i32", "i32", func(left, right data.Value, block *ir.Block) data.Value {
		return data.NewInstruction((block.NewAdd(left.LLVal(block), right.LLVal(block))))
	})

}
