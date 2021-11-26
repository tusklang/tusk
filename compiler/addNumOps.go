package compiler

import (
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/tusklang/tusk/ast"
	"github.com/tusklang/tusk/data"
)

//function to add numerical operations
//TODO: find a way to automate all these
//i copy pasted the code in most of these
func addNumOps(compiler *ast.Compiler) {
	//add all arithmetic operators for numeric types

	var intpuint = make(map[string]data.Type)

	for k, v := range inttypeV {
		intpuint[k] = v
	}

	for k, v := range uinttypeV {
		intpuint[k] = v
	}

	compiler.OperationStore.NewOperation("+", "untypedint", "untypedint", func(left, right data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class) data.Value {
		return data.NewInstVariable(function.ActiveBlock.NewAdd(left.LLVal(function), right.LLVal(function)), left.TType())
	})

	compiler.OperationStore.NewOperation("-", "untypedint", "untypedint", func(left, right data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class) data.Value {
		return data.NewInstVariable(function.ActiveBlock.NewSub(left.LLVal(function), right.LLVal(function)), left.TType())
	})

	compiler.OperationStore.NewOperation("*", "untypedint", "untypedint", func(left, right data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class) data.Value {
		return data.NewInstVariable(function.ActiveBlock.NewMul(left.LLVal(function), right.LLVal(function)), left.TType())
	})

	compiler.OperationStore.NewOperation("/", "untypedint", "untypedint", func(left, right data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class) data.Value {
		return data.NewInstVariable(function.ActiveBlock.NewSDiv(left.LLVal(function), right.LLVal(function)), left.TType())
	})

	compiler.OperationStore.NewOperation("+", "untypedfloat", "untypedint", func(left, right data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class) data.Value {
		ll := right.LLVal(function).(*constant.Int).X
		return data.NewInstVariable(function.ActiveBlock.NewFAdd(left.LLVal(function), constant.NewFloat(types.Double, float64(ll.Int64()))), left.TType())
	})

	compiler.OperationStore.NewOperation("-", "untypedfloat", "untypedint", func(left, right data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class) data.Value {
		ll := right.LLVal(function).(*constant.Int).X
		return data.NewInstVariable(function.ActiveBlock.NewFSub(left.LLVal(function), constant.NewFloat(types.Double, float64(ll.Int64()))), left.TType())
	})

	compiler.OperationStore.NewOperation("*", "untypedfloat", "untypedint", func(left, right data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class) data.Value {
		ll := right.LLVal(function).(*constant.Int).X
		return data.NewInstVariable(function.ActiveBlock.NewFMul(left.LLVal(function), constant.NewFloat(types.Double, float64(ll.Int64()))), left.TType())
	})

	compiler.OperationStore.NewOperation("/", "untypedfloat", "untypedint", func(left, right data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class) data.Value {
		ll := right.LLVal(function).(*constant.Int).X
		return data.NewInstVariable(function.ActiveBlock.NewFDiv(left.LLVal(function), constant.NewFloat(types.Double, float64(ll.Int64()))), left.TType())
	})

	compiler.OperationStore.NewOperation("+", "untypedint", "untypedfloat", func(left, right data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class) data.Value {
		ll := left.LLVal(function).(*constant.Int).X
		return data.NewInstVariable(function.ActiveBlock.NewFAdd(constant.NewFloat(types.Double, float64(ll.Int64())), right.LLVal(function)), left.TType())
	})

	compiler.OperationStore.NewOperation("-", "untypedint", "untypedfloat", func(left, right data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class) data.Value {
		ll := left.LLVal(function).(*constant.Int).X
		return data.NewInstVariable(function.ActiveBlock.NewFSub(constant.NewFloat(types.Double, float64(ll.Int64())), right.LLVal(function)), left.TType())
	})

	compiler.OperationStore.NewOperation("*", "untypedint", "untypedfloat", func(left, right data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class) data.Value {
		ll := left.LLVal(function).(*constant.Int).X
		return data.NewInstVariable(function.ActiveBlock.NewFMul(constant.NewFloat(types.Double, float64(ll.Int64())), right.LLVal(function)), left.TType())
	})

	compiler.OperationStore.NewOperation("/", "untypedint", "untypedfloat", func(left, right data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class) data.Value {
		ll := left.LLVal(function).(*constant.Int).X
		return data.NewInstVariable(function.ActiveBlock.NewFDiv(constant.NewFloat(types.Double, float64(ll.Int64())), right.LLVal(function)), left.TType())
	})

	compiler.OperationStore.NewOperation("+", "untypedfloat", "untypedfloat", func(left, right data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class) data.Value {
		return data.NewInstVariable(function.ActiveBlock.NewFAdd(left.LLVal(function), right.LLVal(function)), left.TType())
	})

	compiler.OperationStore.NewOperation("-", "untypedfloat", "untypedfloat", func(left, right data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class) data.Value {
		return data.NewInstVariable(function.ActiveBlock.NewFSub(left.LLVal(function), right.LLVal(function)), left.TType())
	})

	compiler.OperationStore.NewOperation("*", "untypedfloat", "untypedfloat", func(left, right data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class) data.Value {
		return data.NewInstVariable(function.ActiveBlock.NewFMul(left.LLVal(function), right.LLVal(function)), left.TType())
	})

	compiler.OperationStore.NewOperation("/", "untypedfloat", "untypedfloat", func(left, right data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class) data.Value {
		return data.NewInstVariable(function.ActiveBlock.NewFDiv(left.LLVal(function), right.LLVal(function)), left.TType())
	})

	for k, _v := range intpuint {

		//the v declared in the loop changes per iteration
		//because we use v in the operations, we need a v value that is persistent for each iteration
		var v = _v

		compiler.OperationStore.NewOperation("+", k, k, func(left, right data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class) data.Value {
			return data.NewInstVariable(function.ActiveBlock.NewAdd(left.LLVal(function), right.LLVal(function)), v)
		})

		compiler.OperationStore.NewOperation("-", k, k, func(left, right data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class) data.Value {
			return data.NewInstVariable(function.ActiveBlock.NewSub(left.LLVal(function), right.LLVal(function)), v)
		})

		compiler.OperationStore.NewOperation("*", k, k, func(left, right data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class) data.Value {
			return data.NewInstVariable(function.ActiveBlock.NewMul(left.LLVal(function), right.LLVal(function)), v)
		})

		compiler.OperationStore.NewOperation("+", k, "untypedint", func(left, right data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class) data.Value {
			return data.NewInstVariable(function.ActiveBlock.NewAdd(left.LLVal(function), right.LLVal(function)), v)
		})

		compiler.OperationStore.NewOperation("-", k, "untypedint", func(left, right data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class) data.Value {
			return data.NewInstVariable(function.ActiveBlock.NewSub(left.LLVal(function), right.LLVal(function)), v)
		})

		compiler.OperationStore.NewOperation("*", k, "untypedint", func(left, right data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class) data.Value {
			return data.NewInstVariable(function.ActiveBlock.NewMul(left.LLVal(function), right.LLVal(function)), v)
		})

		compiler.OperationStore.NewOperation("+", "untypedint", k, func(left, right data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class) data.Value {
			return data.NewInstVariable(function.ActiveBlock.NewAdd(left.LLVal(function), right.LLVal(function)), v)
		})

		compiler.OperationStore.NewOperation("-", "untypedint", k, func(left, right data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class) data.Value {
			return data.NewInstVariable(function.ActiveBlock.NewSub(left.LLVal(function), right.LLVal(function)), v)
		})

		compiler.OperationStore.NewOperation("*", "untypedint", k, func(left, right data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class) data.Value {
			return data.NewInstVariable(function.ActiveBlock.NewMul(left.LLVal(function), right.LLVal(function)), v)
		})
	}

	//add division operations
	for k, _v := range inttypeV {

		v := _v

		compiler.OperationStore.NewOperation("/", k, k, func(left, right data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class) data.Value {
			return data.NewInstVariable(function.ActiveBlock.NewSDiv(left.LLVal(function), right.LLVal(function)), v)
		})

		compiler.OperationStore.NewOperation("/", k, "untypedint", func(left, right data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class) data.Value {
			return data.NewInstVariable(function.ActiveBlock.NewSDiv(left.LLVal(function), right.LLVal(function)), v)
		})

		compiler.OperationStore.NewOperation("/", "untypedint", k, func(left, right data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class) data.Value {
			return data.NewInstVariable(function.ActiveBlock.NewSDiv(left.LLVal(function), right.LLVal(function)), v)
		})
	}

	for k, _v := range uinttypeV {

		v := _v

		compiler.OperationStore.NewOperation("/", k, k, func(left, right data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class) data.Value {
			return data.NewInstVariable(function.ActiveBlock.NewUDiv(left.LLVal(function), right.LLVal(function)), v)
		})

		compiler.OperationStore.NewOperation("/", k, "untypedint", func(left, right data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class) data.Value {
			return data.NewInstVariable(function.ActiveBlock.NewUDiv(left.LLVal(function), right.LLVal(function)), v)
		})

		compiler.OperationStore.NewOperation("/", "untypedint", k, func(left, right data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class) data.Value {
			return data.NewInstVariable(function.ActiveBlock.NewUDiv(left.LLVal(function), right.LLVal(function)), v)
		})
	}

	for k, _v := range floattypeV {
		var v = _v

		compiler.OperationStore.NewOperation("+", k, k, func(left, right data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class) data.Value {
			return data.NewInstVariable(function.ActiveBlock.NewFAdd(left.LLVal(function), right.LLVal(function)), v)
		})

		compiler.OperationStore.NewOperation("-", k, k, func(left, right data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class) data.Value {
			return data.NewInstVariable(function.ActiveBlock.NewFSub(left.LLVal(function), right.LLVal(function)), v)
		})

		compiler.OperationStore.NewOperation("*", k, k, func(left, right data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class) data.Value {
			return data.NewInstVariable(function.ActiveBlock.NewFMul(left.LLVal(function), right.LLVal(function)), v)
		})

		compiler.OperationStore.NewOperation("/", k, k, func(left, right data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class) data.Value {
			return data.NewInstVariable(function.ActiveBlock.NewFDiv(left.LLVal(function), right.LLVal(function)), v)
		})

		compiler.OperationStore.NewOperation("+", k, "untypedint", func(left, right data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class) data.Value {
			ll := right.LLVal(function).(*constant.Int).X
			return data.NewInstVariable(function.ActiveBlock.NewFAdd(left.LLVal(function), constant.NewFloat(types.Double, float64(ll.Int64()))), v)
		})

		compiler.OperationStore.NewOperation("-", k, "untypedint", func(left, right data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class) data.Value {
			ll := right.LLVal(function).(*constant.Int).X
			return data.NewInstVariable(function.ActiveBlock.NewFSub(left.LLVal(function), constant.NewFloat(types.Double, float64(ll.Int64()))), v)
		})

		compiler.OperationStore.NewOperation("*", k, "untypedint", func(left, right data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class) data.Value {
			ll := right.LLVal(function).(*constant.Int).X
			return data.NewInstVariable(function.ActiveBlock.NewFMul(left.LLVal(function), constant.NewFloat(types.Double, float64(ll.Int64()))), v)
		})

		compiler.OperationStore.NewOperation("/", k, "untypedint", func(left, right data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class) data.Value {
			ll := right.LLVal(function).(*constant.Int).X
			return data.NewInstVariable(function.ActiveBlock.NewFDiv(left.LLVal(function), constant.NewFloat(types.Double, float64(ll.Int64()))), v)
		})

		compiler.OperationStore.NewOperation("+", k, "untypedfloat", func(left, right data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class) data.Value {
			return data.NewInstVariable(function.ActiveBlock.NewFAdd(left.LLVal(function), right.LLVal(function)), v)
		})

		compiler.OperationStore.NewOperation("-", k, "untypedfloat", func(left, right data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class) data.Value {
			return data.NewInstVariable(function.ActiveBlock.NewFSub(left.LLVal(function), right.LLVal(function)), v)
		})

		compiler.OperationStore.NewOperation("*", k, "untypedfloat", func(left, right data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class) data.Value {
			return data.NewInstVariable(function.ActiveBlock.NewFMul(left.LLVal(function), right.LLVal(function)), v)
		})

		compiler.OperationStore.NewOperation("/", k, "untypedfloat", func(left, right data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class) data.Value {
			return data.NewInstVariable(function.ActiveBlock.NewFDiv(left.LLVal(function), right.LLVal(function)), v)
		})

		compiler.OperationStore.NewOperation("+", "untypedint", k, func(left, right data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class) data.Value {
			ll := left.LLVal(function).(*constant.Int).X
			return data.NewInstVariable(function.ActiveBlock.NewFAdd(constant.NewFloat(types.Double, float64(ll.Int64())), right.LLVal(function)), v)
		})

		compiler.OperationStore.NewOperation("-", "untypedint", k, func(left, right data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class) data.Value {
			ll := left.LLVal(function).(*constant.Int).X
			return data.NewInstVariable(function.ActiveBlock.NewFSub(constant.NewFloat(types.Double, float64(ll.Int64())), right.LLVal(function)), v)
		})

		compiler.OperationStore.NewOperation("*", "untypedint", k, func(left, right data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class) data.Value {
			ll := left.LLVal(function).(*constant.Int).X
			return data.NewInstVariable(function.ActiveBlock.NewFMul(constant.NewFloat(types.Double, float64(ll.Int64())), right.LLVal(function)), v)
		})

		compiler.OperationStore.NewOperation("/", "untypedint", k, func(left, right data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class) data.Value {
			ll := left.LLVal(function).(*constant.Int).X
			return data.NewInstVariable(function.ActiveBlock.NewFDiv(constant.NewFloat(types.Double, float64(ll.Int64())), right.LLVal(function)), v)
		})

		compiler.OperationStore.NewOperation("+", "untypedfloat", k, func(left, right data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class) data.Value {
			return data.NewInstVariable(function.ActiveBlock.NewFAdd(left.LLVal(function), right.LLVal(function)), v)
		})

		compiler.OperationStore.NewOperation("-", "untypedfloat", k, func(left, right data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class) data.Value {
			return data.NewInstVariable(function.ActiveBlock.NewFSub(left.LLVal(function), right.LLVal(function)), v)
		})

		compiler.OperationStore.NewOperation("*", "untypedfloat", k, func(left, right data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class) data.Value {
			return data.NewInstVariable(function.ActiveBlock.NewFMul(left.LLVal(function), right.LLVal(function)), v)
		})

		compiler.OperationStore.NewOperation("/", "untypedfloat", k, func(left, right data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class) data.Value {
			return data.NewInstVariable(function.ActiveBlock.NewFDiv(left.LLVal(function), right.LLVal(function)), v)
		})

	}
}
