package compiler

import (
	"github.com/tusklang/tusk/ast"
	"github.com/tusklang/tusk/data"
)

//function to add numerical operations
//TODO: find a way to automate all these
//i copy pasted the code in most of these
func addNumOps(compiler *ast.Compiler) {
	//add all arithmetic operators for numeric types

	var intpuint = make(map[string]data.Type)

	for k, v := range ast.InttypeV {
		intpuint[k] = v
	}

	for k, v := range ast.UinttypeV {
		intpuint[k] = v
	}

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

	}

	//division for int and uint are a bit different
	for k, _v := range ast.InttypeV {
		v := _v

		compiler.OperationStore.NewOperation("/", k, k, func(left, right data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class) data.Value {
			return data.NewInstVariable(function.ActiveBlock.NewSDiv(left.LLVal(function), right.LLVal(function)), v)
		})
	}
	for k, _v := range ast.UinttypeV {
		v := _v

		compiler.OperationStore.NewOperation("/", k, k, func(left, right data.Value, compiler *ast.Compiler, function *data.Function, class *data.Class) data.Value {
			return data.NewInstVariable(function.ActiveBlock.NewUDiv(left.LLVal(function), right.LLVal(function)), v)
		})
	}

	for k, _v := range ast.FloattypeV {
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

	}
}
