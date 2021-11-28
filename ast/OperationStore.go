package ast

import (
	"github.com/tusklang/tusk/data"
)

type operationdef struct {
	left, right string
	operation   string
	handler     func(left, right data.Value, lcg, rcg Group, compiler *Compiler, function *data.Function, class *data.Class) data.Value
}

type OperationStore struct {
	operations []operationdef
}

func NewOperationStore() *OperationStore {
	return &OperationStore{}
}

func (os *OperationStore) NewOperation(operation string, ltype, rtype string, handler func(left, right data.Value, lcg, rcg Group, compiler *Compiler, function *data.Function, class *data.Class) data.Value) {
	os.operations = append(os.operations, operationdef{
		left:      ltype,
		right:     rtype,
		operation: operation,
		handler:   handler,
	})
}

//precendence that untyped vals will be casted to automatically for operations
var (
	untypedintprec   = []string{"i32", "i64", "i16", "i8", "f64", "f32"}
	untypedfloatprec = []string{"f64", "f32"}
)

func (os *OperationStore) checkuntypedr(lval, rval data.Value, lcg, rcg Group, operation string, compiler *Compiler, function *data.Function, class *data.Class) data.Value {

	//store the original
	krval := rval

	var prec []string

	if rval != nil && rval.TType() != nil {
		if rval.TType().TypeData().Name() == "untypedint" {
			prec = untypedintprec
		} else if rval.TType().TypeData().Name() == "untypedfloat" {
			prec = untypedfloatprec
		}
	}

	for _, vi := range prec {
		if ok := os.runoperation(
			lval,
			compiler.CastStore.RunCast(true, Numtypes[vi], krval, rcg, compiler, function, class),
			lcg,
			rcg,
			operation,
			compiler,
			function,
			class,
		); ok != nil {
			return ok
		}
	}

	return os.runoperation(lval, rval, lcg, rcg, operation, compiler, function, class)
}

func (os *OperationStore) checkuntypedl(lval, rval data.Value, lcg, rcg Group, operation string, compiler *Compiler, function *data.Function, class *data.Class) data.Value {

	//store the original
	klval := lval

	var prec []string

	if lval != nil && lval.TType() != nil {
		if lval.TType().TypeData().Name() == "untypedint" {
			prec = untypedintprec
		} else if lval.TType().TypeData().Name() == "untypedfloat" {
			prec = untypedfloatprec
		}
	}

	for _, vi := range prec {
		if ok := os.checkuntypedr(
			compiler.CastStore.RunCast(true, Numtypes[vi], klval, lcg, compiler, function, class),
			rval,
			lcg,
			rcg,
			operation,
			compiler,
			function,
			class,
		); ok != nil {
			return ok
		}
	}

	return os.checkuntypedr(lval, rval, lcg, rcg, operation, compiler, function, class)
}

func (os *OperationStore) runoperation(lval, rval data.Value, lcg, rcg Group, operation string, compiler *Compiler, function *data.Function, class *data.Class) data.Value {
	for _, v := range os.operations {
		if operation == v.operation && matchOpdef(lval, v.left) && matchOpdef(rval, v.right) {
			//if the types match with the operation
			return v.handler(lval, rval, lcg, rcg, compiler, function, class)
		}
	}

	//there isn't an operation matching the given types
	return nil
}

func (os *OperationStore) RunOperation(lval, rval data.Value, lcg, rcg Group, operation string, compiler *Compiler, function *data.Function, class *data.Class) data.Value {
	return os.checkuntypedl(lval, rval, lcg, rcg, operation, compiler, function, class)
}
