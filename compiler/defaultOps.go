package compiler

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	"github.com/tusklang/tusk/ast"
	"github.com/tusklang/tusk/data"
)

func initDefaultOps(compiler *ast.Compiler) {

	compiler.OperationStore = ast.NewOperationStore()

	compiler.OperationStore.NewOperation("=", "ptr", "*", func(left, right data.Value, compiler *ast.Compiler, block *ir.Block) data.Value {

		varv := left.(*data.Variable).FetchAssig() //we know it's a variable, so we can assert it and fetch the assignment instruction
		toassign := right.LLVal(block)

		block.NewStore(toassign, varv)

		return left
	})

	compiler.OperationStore.NewOperation("+", "i32", "i32", func(left, right data.Value, compiler *ast.Compiler, block *ir.Block) data.Value {
		return data.NewVariable(block.NewAdd(left.LLVal(block), right.LLVal(block)), data.NewPrimitive(types.I32))
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

	compiler.OperationStore.NewOperation(".", "instance", "udvar", func(left, right data.Value, compiler *ast.Compiler, block *ir.Block) data.Value {

		inst := left.(*data.Variable).FetchAssig() //it's always a variable
		sub := right.(*data.UndeclaredVar).Name

		classt := left.TType().(*data.Instance).Class

		return data.NewVariable(
			block.NewGetElementPtr(
				classt.SType,
				inst,
				constant.NewInt(types.I32, 0),
				constant.NewInt(types.I32, classt.Instance[sub].Index),
			),
			data.NewPointer(classt.Instance[sub].Type),
		)
	})

	compiler.OperationStore.NewOperation("()", "func", "fncallb", func(left, right data.Value, compiler *ast.Compiler, block *ir.Block) data.Value {

		f := left.LLVal(block)
		fcb := right.(*data.FnCallBlock)

		var args []value.Value

		for _, v := range fcb.Args {
			args = append(args, v.LLVal(block))
		}

		call := block.NewCall(f, args...)

		return data.NewVariable(
			call,
			left.TType().(*data.Function).RetType(),
		)
	})

	compiler.OperationStore.NewOperation("()", "class", "fncallb", func(left, right data.Value, compiler *ast.Compiler, block *ir.Block) data.Value {

		class := left.(*data.Class)
		fcb := right.(*data.FnCallBlock)

		var args []value.Value

		for _, v := range fcb.Args {
			args = append(args, v.LLVal(block))
		}

		return data.NewVariable(
			block.NewCall(class.Construct.LLFunc, args...),
			data.NewInstance(class),
		)
	})

}
