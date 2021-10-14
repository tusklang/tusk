package compiler

import (
	"fmt"

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

	compiler.OperationStore.NewOperation("=", "var", "*", func(left, right data.Value, compiler *ast.Compiler, block *ir.Block, class *data.Class) data.Value {

		varv := left.(*data.Variable).FetchAssig() //we know it's a variable, so we can assert it and fetch the assignment instruction
		toassign := right.LLVal(block)

		block.NewStore(toassign, varv)

		return left
	})

	compiler.OperationStore.NewOperation("+", "i32", "i32", func(left, right data.Value, compiler *ast.Compiler, block *ir.Block, class *data.Class) data.Value {
		return data.NewInstVariable(block.NewAdd(left.LLVal(block), right.LLVal(block)), data.NewPrimitive(types.I32))
	})

	compiler.OperationStore.NewOperation(".", "package", "udvar", func(left, right data.Value, compiler *ast.Compiler, block *ir.Block, class *data.Class) data.Value {

		pack := left.(*data.Package)
		sub := right.(*data.UndeclaredVar).Name

		//it can either be a class or a subpackage
		var (
			cclass  = pack.Classes[sub]
			subpack = pack.ChildPacks[sub]
		)

		if cclass == nil {
			return subpack
		}

		return class
	})

	compiler.OperationStore.NewOperation(".", "class", "udvar", func(left, right data.Value, compiler *ast.Compiler, block *ir.Block, class *data.Class) data.Value {

		cclass := left.(*data.Class)
		sub := right.(*data.UndeclaredVar).Name

		if cclass.Static[sub].Access == 0 && !cclass.Equals(class) {
			//error
			//trying to access a private field
			fmt.Println("trying to access private static field " + class.Name)
		}

		return cclass.Static[sub].Value
	})

	compiler.OperationStore.NewOperation(".", "instance", "udvar", func(left, right data.Value, compiler *ast.Compiler, block *ir.Block, class *data.Class) data.Value {

		inst := left.LLVal(block)
		sub := right.(*data.UndeclaredVar).Name

		classt := left.TType().(*data.Instance).Class

		if classt.Instance[sub].Access == 0 && !classt.Equals(class) {
			//error
			//trying to access a private field
			fmt.Println("trying to access private instance field " + class.Name)
		}

		return data.NewVariable(
			block.NewGetElementPtr(
				classt.SType,
				block.NewLoad(types.NewPointer(classt.SType), inst),
				constant.NewInt(types.I32, 0),
				constant.NewInt(types.I32, classt.Instance[sub].Index),
			),
			classt.Instance[sub].Type,
		)
	})

	compiler.OperationStore.NewOperation("()", "func", "fncallb", func(left, right data.Value, compiler *ast.Compiler, block *ir.Block, class *data.Class) data.Value {

		f := left.LLVal(block)
		fcb := right.(*data.FnCallBlock)

		var args []value.Value

		for _, v := range fcb.Args {
			args = append(args, v.LLVal(block))
		}

		call := block.NewCall(f, args...)

		return data.NewInstVariable(
			call,
			left.TType().(*data.Function).RetType(),
		)
	})

	compiler.OperationStore.NewOperation("()", "class", "fncallb", func(left, right data.Value, compiler *ast.Compiler, block *ir.Block, class *data.Class) data.Value {

		cclass := left.(*data.Class)
		fcb := right.(*data.FnCallBlock)

		var args []value.Value

		for _, v := range fcb.Args {
			args = append(args, v.LLVal(block))
		}

		return data.NewInstVariable(
			block.NewCall(cclass.Construct.LLFunc, args...),
			data.NewInstance(cclass),
		)
	})

	compiler.OperationStore.NewOperation("*", "-", "ptr&var", func(left, right data.Value, compiler *ast.Compiler, block *ir.Block, class *data.Class) data.Value {
		return data.NewVariable(
			right.LLVal(block),
			right.TType().(*data.Pointer).PType(),
		)
	})

	compiler.OperationStore.NewOperation("*", "-", "type", func(left, right data.Value, compiler *ast.Compiler, block *ir.Block, class *data.Class) data.Value {

		ptrt := data.NewPointer(right.TType())
		ptrt.SetToType() //make it a type, not a value

		return ptrt
	})

	compiler.OperationStore.NewOperation("==", "*", "*", func(left, right data.Value, compiler *ast.Compiler, block *ir.Block, class *data.Class) data.Value {
		return data.NewInstVariable(
			block.NewICmp(enum.IPredEQ, left.LLVal(block), right.LLVal(block)),
			data.NewPrimitive(types.I1),
		)
	})

	compiler.OperationStore.NewOperation("!=", "*", "*", func(left, right data.Value, compiler *ast.Compiler, block *ir.Block, class *data.Class) data.Value {
		return data.NewInstVariable(
			block.NewICmp(enum.IPredNE, left.LLVal(block), right.LLVal(block)),
			data.NewPrimitive(types.I1),
		)
	})

	compiler.OperationStore.NewOperation(">", "*", "*", func(left, right data.Value, compiler *ast.Compiler, block *ir.Block, class *data.Class) data.Value {
		return data.NewInstVariable(
			block.NewZExt(block.NewICmp(enum.IPredUGT, left.LLVal(block), right.LLVal(block)), types.I32),
			data.NewPrimitive(types.I1),
		)
	})

	compiler.OperationStore.NewOperation(">=", "*", "*", func(left, right data.Value, compiler *ast.Compiler, block *ir.Block, class *data.Class) data.Value {
		return data.NewInstVariable(
			block.NewICmp(enum.IPredUGE, left.LLVal(block), right.LLVal(block)),
			data.NewPrimitive(types.I1),
		)
	})

	compiler.OperationStore.NewOperation("<", "*", "*", func(left, right data.Value, compiler *ast.Compiler, block *ir.Block, class *data.Class) data.Value {
		return data.NewInstVariable(
			block.NewICmp(enum.IPredULT, left.LLVal(block), right.LLVal(block)),
			data.NewPrimitive(types.I1),
		)
	})

	compiler.OperationStore.NewOperation("<=", "*", "*", func(left, right data.Value, compiler *ast.Compiler, block *ir.Block, class *data.Class) data.Value {
		return data.NewInstVariable(
			block.NewICmp(enum.IPredULE, left.LLVal(block), right.LLVal(block)),
			data.NewPrimitive(types.I1),
		)
	})

	compiler.OperationStore.NewOperation("&", "-", "var", func(left, right data.Value, compiler *ast.Compiler, block *ir.Block, class *data.Class) data.Value {
		vd := data.NewInstVariable(
			right.(*data.Variable).FetchAssig(),
			data.NewPointer(right.(*data.Variable).TType()),
		)
		return vd
	})

}
