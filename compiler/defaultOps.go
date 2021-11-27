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

	addNumOps(compiler)

	compiler.OperationStore.NewOperation("=", "var", "*", func(left, right data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class) data.Value {

		if !left.TType().Equals(right.TType()) {
			right = compiler.CastStore.RunCast(true, left.TType(), right, compiler, function, class)
			if right == nil {
				//error
				//dst and src types don't match
			}
		}

		var varv value.Value

		switch varvt := left.(type) {
		case *data.Variable:
			varv = varvt.FetchAssig()
		case *data.InstanceVariable:
			varv = varvt.FetchAssig()
		}

		toassign := right.LLVal(function)

		function.ActiveBlock.NewStore(toassign, varv)

		return left
	})

	compiler.OperationStore.NewOperation("->", "type", "*", func(left, right data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class) data.Value {
		return compiler.CastStore.RunCast(false, left.(data.Type), right, compiler, function, class)
	})

	compiler.OperationStore.NewOperation(".", "package", "udvar", func(left, right data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class) data.Value {

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

	compiler.OperationStore.NewOperation(".", "class", "udvar", func(left, right data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class) data.Value {

		cclass := left.(*data.Class)
		sub := right.(*data.UndeclaredVar).Name

		if cclass.Static[sub].Access == 2 && !cclass.Equals(class) {
			//error
			//trying to access a private field
			fmt.Println("trying to access private static field " + class.Name)
		}

		return cclass.Static[sub].Value
	})

	compiler.OperationStore.NewOperation(".", "instance", "udvar", func(left, right data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class) data.Value {

		inst := left.LLVal(function)
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
			gep := function.ActiveBlock.NewGetElementPtr(
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

	compiler.OperationStore.NewOperation("()", "func", "fncallb", func(left, right data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class) data.Value {

		f := left.LLVal(function)
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
				if cast := compiler.CastStore.RunCast(true, v, fcb.Args[k], compiler, function, class); cast != nil {
					fcb.Args[k] = cast
				} else {
					//compiler error
					//variable value type doesn't match inputted type
				}
			}
			args = append(args, fcb.Args[k].LLVal(function))
		}

		var call value.Value = function.ActiveBlock.NewCall(f, args...)

		if left.TypeData().HasFlag("linked") {
			call.(*ir.InstCall).Sig().Params = nil
			call.(*ir.InstCall).Sig().Variadic = true

			//linked functions always have a pointer, integer, or void return type
			rettype := tf.RetType().Type()

			if types.IsPointer(rettype) {
				//use a bitcast for a pointer return
				call = function.ActiveBlock.NewBitCast(call, rettype)
			} else if types.IsInt(rettype) {
				//use an ptrtoint cast for an integer return
				call = function.ActiveBlock.NewPtrToInt(call, rettype)
			}

		}

		return data.NewInstVariable(
			call,
			tf.RetType(),
		)
	})

	compiler.OperationStore.NewOperation("()", "class", "fncallb", func(left, right data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class) data.Value {

		cclass := left.(*data.Class)
		fcb := right.(*data.FnCallBlock)

		return compiler.OperationStore.RunOperation(cclass.Construct, fcb, "()", compiler, function, class)
	})

	//array indexing
	compiler.OperationStore.NewOperation("[]", "slice&array", "i32", func(left, right data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class) data.Value {
		gept := left.TType().(*data.SliceArray).ValType().Type()
		gep := function.ActiveBlock.NewGetElementPtr(gept, left.LLVal(function), right.LLVal(function))
		gep.InBounds = true
		return data.NewVariable(
			gep,
			left.TType().(*data.SliceArray).ValType(),
		)
	})

	compiler.OperationStore.NewOperation("[]", "fixed&array", "i32", func(left, right data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class) data.Value {
		farr := left.TType().(*data.FixedArray)
		gept := farr.Type()

		//might be a bit of a hack solution, but since a fixed array's llval gives the loaded array rather than the array's pointer
		//we need to create a pointer for GEP here...
		//so we need another alloca...

		alc := function.ActiveBlock.NewAlloca(gept)
		function.ActiveBlock.NewStore(left.LLVal(function), alc)

		//llvm optimization should take care of it... right?

		gep := function.ActiveBlock.NewGetElementPtr(gept, alc, constant.NewInt(types.I32, 0), right.LLVal(function))
		gep.InBounds = true
		return data.NewVariable(
			gep,
			farr.ValType(),
		)
	})

	compiler.OperationStore.NewOperation("[]", "varied&array", "i32", func(left, right data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class) data.Value {
		varr := left.TType().(*data.VariedLengthArray)
		gept := varr.ValType().Type()
		gep := function.ActiveBlock.NewGetElementPtr(gept, left.LLVal(function), right.LLVal(function))
		gep.InBounds = true
		return data.NewVariable(
			gep,
			varr.ValType(),
		)
	})
	////////////////

	compiler.OperationStore.NewOperation("#", "-", "ptr&var", func(left, right data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class) data.Value {
		return data.NewVariable(
			right.LLVal(function),
			right.TType().(*data.Pointer).PType(),
		)
	})

	compiler.OperationStore.NewOperation("#", "-", "type", func(left, right data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class) data.Value {

		ptrt := data.NewPointer(right.TType())
		ptrt.SetToType() //make it a type, not a value

		return ptrt
	})

	compiler.OperationStore.NewOperation("==", "*", "*", func(left, right data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class) data.Value {
		return data.NewInstVariable(
			function.ActiveBlock.NewICmp(enum.IPredEQ, left.LLVal(function), right.LLVal(function)),
			data.NewPrimitive(types.I1),
		)
	})

	compiler.OperationStore.NewOperation("!=", "*", "*", func(left, right data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class) data.Value {
		return data.NewInstVariable(
			function.ActiveBlock.NewICmp(enum.IPredNE, left.LLVal(function), right.LLVal(function)),
			data.NewPrimitive(types.I1),
		)
	})

	compiler.OperationStore.NewOperation(">", "*", "*", func(left, right data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class) data.Value {
		return data.NewInstVariable(
			function.ActiveBlock.NewZExt(function.ActiveBlock.NewICmp(enum.IPredUGT, left.LLVal(function), right.LLVal(function)), types.I32),
			data.NewPrimitive(types.I1),
		)
	})

	compiler.OperationStore.NewOperation(">=", "*", "*", func(left, right data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class) data.Value {
		return data.NewInstVariable(
			function.ActiveBlock.NewICmp(enum.IPredUGE, left.LLVal(function), right.LLVal(function)),
			data.NewPrimitive(types.I1),
		)
	})

	compiler.OperationStore.NewOperation("<", "*", "*", func(left, right data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class) data.Value {
		return data.NewInstVariable(
			function.ActiveBlock.NewICmp(enum.IPredULT, left.LLVal(function), right.LLVal(function)),
			data.NewPrimitive(types.I1),
		)
	})

	compiler.OperationStore.NewOperation("<=", "*", "*", func(left, right data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class) data.Value {
		return data.NewInstVariable(
			function.ActiveBlock.NewICmp(enum.IPredULE, left.LLVal(function), right.LLVal(function)),
			data.NewPrimitive(types.I1),
		)
	})

	compiler.OperationStore.NewOperation("@", "-", "var", func(left, right data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class) data.Value {
		vd := data.NewInstVariable(
			right.(*data.Variable).FetchAssig(),
			data.NewPointer(right.(*data.Variable).TType()),
		)
		return vd
	})

}
