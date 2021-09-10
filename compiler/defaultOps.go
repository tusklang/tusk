package compiler

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/value"
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

	compiler.OperationStore.NewOperation(".", "class", "udvar", func(left, right data.Value, compiler *ast.Compiler, block *ir.Block) data.Value {

		class := left.(*data.Class)
		sub := right.(*data.UndeclaredVar).Name

		return class.Static[sub]
	})

	compiler.OperationStore.NewOperation("()", "func", "fncallb", func(left, right data.Value, compiler *ast.Compiler, block *ir.Block) data.Value {

		f := left.LLVal(block)
		fcb := right.(*data.FnCallBlock)

		var args []value.Value

		for _, v := range fcb.Args {
			args = append(args, v.LLVal(block))
		}

		return data.NewInstruction(
			block.NewCall(f, args...),
		)
	})

}
