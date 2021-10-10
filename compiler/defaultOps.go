package compiler

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
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

		callv := data.NewVariable(
			call,
			left.TType().(*data.Function).RetType(),
		)

		//set the load instruction to just fetch the assignment
		callv.SetLoadInst(func(v *data.Variable, block *ir.Block) value.Value {
			return v.FetchAssig()
		})

		return callv
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

	compiler.OperationStore.NewOperation("*", "-", "ptr&var", func(left, right data.Value, compiler *ast.Compiler, block *ir.Block) data.Value {
		vd := data.NewVariable(
			right.LLVal(block),
			right.TType().(*data.Pointer).PType(),
		)
		vd.SetLoadInst(func(v *data.Variable, block *ir.Block) value.Value {
			return block.NewLoad(v.Type(), v.FetchAssig())
		})
		return vd
	})

	compiler.OperationStore.NewOperation("*", "-", "type", func(left, right data.Value, compiler *ast.Compiler, block *ir.Block) data.Value {

		ptrt := data.NewPointer(right.TType())
		ptrt.SetToType() //make it a type, not a value

		return ptrt
	})

	compiler.OperationStore.NewOperation("==", "*", "*", func(left, right data.Value, compiler *ast.Compiler, block *ir.Block) data.Value {
		return data.NewInstVariable(
			block.NewICmp(enum.IPredEQ, left.LLVal(block), right.LLVal(block)),
			data.NewPrimitive(types.I1),
		)
	})

	compiler.OperationStore.NewOperation("!=", "*", "*", func(left, right data.Value, compiler *ast.Compiler, block *ir.Block) data.Value {
		return data.NewInstVariable(
			block.NewICmp(enum.IPredNE, left.LLVal(block), right.LLVal(block)),
			data.NewPrimitive(types.I1),
		)
	})

	compiler.OperationStore.NewOperation(">", "*", "*", func(left, right data.Value, compiler *ast.Compiler, block *ir.Block) data.Value {
		return data.NewInstVariable(
			block.NewZExt(block.NewICmp(enum.IPredUGT, left.LLVal(block), right.LLVal(block)), types.I32),
			data.NewPrimitive(types.I1),
		)
	})

	compiler.OperationStore.NewOperation(">=", "*", "*", func(left, right data.Value, compiler *ast.Compiler, block *ir.Block) data.Value {
		return data.NewInstVariable(
			block.NewICmp(enum.IPredUGE, left.LLVal(block), right.LLVal(block)),
			data.NewPrimitive(types.I1),
		)
	})

	compiler.OperationStore.NewOperation("<", "*", "*", func(left, right data.Value, compiler *ast.Compiler, block *ir.Block) data.Value {
		return data.NewInstVariable(
			block.NewICmp(enum.IPredULT, left.LLVal(block), right.LLVal(block)),
			data.NewPrimitive(types.I1),
		)
	})

	compiler.OperationStore.NewOperation("<=", "*", "*", func(left, right data.Value, compiler *ast.Compiler, block *ir.Block) data.Value {
		return data.NewInstVariable(
			block.NewICmp(enum.IPredULE, left.LLVal(block), right.LLVal(block)),
			data.NewPrimitive(types.I1),
		)
	})

	compiler.OperationStore.NewOperation("&", "-", "var", func(left, right data.Value, compiler *ast.Compiler, block *ir.Block) data.Value {
		vd := data.NewInstVariable(
			right.(*data.Variable).FetchAssig(),
			data.NewPointer(right.(*data.Variable).TType()),
		)
		return vd
	})

}
