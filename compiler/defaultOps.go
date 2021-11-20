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

		var varv value.Value

		switch varvt := left.(type) {
		case *data.Variable:
			varv = varvt.FetchAssig()
		case *data.InstanceVariable:
			varv = varvt.FetchAssig()
		}

		toassign := right.LLVal(block)

		block.NewStore(toassign, varv)

		return left
	})

	//add all arithmetic operators for numeric types
	for k, _v := range numtypes {

		//the v declared in the loop changes per iteration
		//because we use v in the operations, we need a v value that is persistent for each iteration
		var v = _v

		compiler.OperationStore.NewOperation("+", k, k, func(left, right data.Value, compiler *ast.Compiler, block *ir.Block, class *data.Class) data.Value {
			return data.NewInstVariable(block.NewAdd(left.LLVal(block), right.LLVal(block)), v)
		})

		compiler.OperationStore.NewOperation("-", k, k, func(left, right data.Value, compiler *ast.Compiler, block *ir.Block, class *data.Class) data.Value {
			return data.NewInstVariable(block.NewSub(left.LLVal(block), right.LLVal(block)), v)
		})

		compiler.OperationStore.NewOperation("*", k, k, func(left, right data.Value, compiler *ast.Compiler, block *ir.Block, class *data.Class) data.Value {
			return data.NewInstVariable(block.NewMul(left.LLVal(block), right.LLVal(block)), v)
		})

		compiler.OperationStore.NewOperation("/", k, k, func(left, right data.Value, compiler *ast.Compiler, block *ir.Block, class *data.Class) data.Value {
			return data.NewInstVariable(block.NewSDiv(left.LLVal(block), right.LLVal(block)), v)
		})
	}

	compiler.OperationStore.NewOperation("->", "type", "*", func(left, right data.Value, compiler *ast.Compiler, block *ir.Block, class *data.Class) data.Value {
		return compiler.CastStore.RunCast(false, left.TypeData().Name(), right, compiler, block, class)
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

		return cclass
	})

	compiler.OperationStore.NewOperation(".", "class", "udvar", func(left, right data.Value, compiler *ast.Compiler, block *ir.Block, class *data.Class) data.Value {

		cclass := left.(*data.Class)
		sub := right.(*data.UndeclaredVar).Name

		if cclass.Static[sub].Access == 2 && !cclass.Equals(class) {
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

		var (
			ivar     *data.ClassField
			ok       bool
			fieldtyp string
		)

		if ivar, ok = classt.Instance[sub]; !ok {
			if ivar, ok = classt.Methods[sub]; !ok {
				//error
				//field `sub` does not exist in class
			} else {
				fieldtyp = "method"
			}
		} else {
			fieldtyp = "var"
		}

		if ivar.Access == 2 && !classt.Equals(class) {
			//error
			//trying to access a private field
			fmt.Println("trying to access private instance field " + class.Name)
		}

		switch fieldtyp {
		case "method":
			//method

			cloned := data.CloneFunc(ivar.Value.(*data.Function))
			cloned.Instance = inst

			return cloned
		case "var":
			//instance variable
			gep := block.NewGetElementPtr(
				classt.SType,
				inst,
				constant.NewInt(types.I32, 0),
				constant.NewInt(types.I32, ivar.Index),
			)
			gep.InBounds = true

			return data.NewInstanceVariable(
				data.NewVariable(
					gep,
					classt.Instance[sub].Type,
				),
				inst,
			)
		}

		return nil
	})

	compiler.OperationStore.NewOperation("()", "func", "fncallb", func(left, right data.Value, compiler *ast.Compiler, block *ir.Block, class *data.Class) data.Value {

		f := left.LLVal(block)
		fcb := right.(*data.FnCallBlock)

		tf := left.TType().(*data.Function)

		var args []value.Value
		var tad int //this value is a boolean (int) used to store if the function is a method or not

		if left.TypeData().HasFlag("method") {
			args = append(args, left.InstanceV())
			tad = 1
		}

		if len(fcb.Args) != len(tf.ParamTypes)+tad {
			//error
			//args given doesn't match args in sig
		}

		for k, v := range tf.ParamTypes {
			if !v.Equals(fcb.Args[k].TType()) {
				if cast := compiler.CastStore.RunCast(true, v.TypeData().Name(), fcb.Args[k], compiler, block, class); cast != nil {
					fcb.Args[k] = cast
				} else {
					//compiler error
					//variable value type doesn't match inputted type
				}
			}
			args = append(args, fcb.Args[k].LLVal(block))
		}

		var call value.Value = block.NewCall(f, args...)

		if left.TypeData().HasFlag("linked") {
			call.(*ir.InstCall).Sig().Params = nil
			call.(*ir.InstCall).Sig().Variadic = true

			//linked functions always have a pointer, integer, or void return type
			rettype := tf.RetType().Type()

			if types.IsPointer(rettype) {
				//use a bitcast for a pointer return
				call = block.NewBitCast(call, rettype)
			} else if types.IsInt(rettype) {
				//use an ptrtoint cast for an integer return
				call = block.NewPtrToInt(call, rettype)
			}

		}

		return data.NewInstVariable(
			call,
			tf.RetType(),
		)
	})

	compiler.OperationStore.NewOperation("()", "class", "fncallb", func(left, right data.Value, compiler *ast.Compiler, block *ir.Block, class *data.Class) data.Value {

		cclass := left.(*data.Class)
		fcb := right.(*data.FnCallBlock)

		return compiler.OperationStore.RunOperation(cclass.Construct, fcb, "()", compiler, block, class)
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
