package compiler

import (
	"github.com/llir/llvm/ir"
	"github.com/tusklang/tusk/ast"
	"github.com/tusklang/tusk/data"
)

func initDefaultOps(compiler *ast.Compiler) {

	compiler.OperationStore = ast.NewOperationStore()

	compiler.OperationStore.NewOperation("+", "i32", "i32", func(left, right data.Value, compiler *ast.Compiler, block *ir.Block) data.Value {
		return data.NewInstruction((block.NewAdd(left.LLVal(block), right.LLVal(block))))
	})

	compiler.OperationStore.NewOperation(".", "package", "udvar", func(left, right data.Value, compiler *ast.Compiler, block *ir.Block) data.Value {

		pack := left.(*data.Package)
		sub := right.(*data.UndeclaredVar).Name

		//it can either be a class or a subpackage
		var (
			class   = pack.Classes[sub]
			subpack = pack.ChildPacks[sub]
		)

		if class == nil {
			return subpack
		}

		return class
	})

}
