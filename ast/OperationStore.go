package ast

import (
	"github.com/tusklang/tusk/data"
)

type operationdef struct {
	left, right string
	operation   string
	handler     func(left, right data.Value, compiler *Compiler, function *data.Function, class *data.Class) data.Value
}

type OperationStore struct {
	operations []operationdef
}

func NewOperationStore() *OperationStore {
	return &OperationStore{}
}

func (os *OperationStore) NewOperation(operation string, ltype, rtype string, handler func(left, right data.Value, compiler *Compiler, function *data.Function, class *data.Class) data.Value) {
	os.operations = append(os.operations, operationdef{
		left:      ltype,
		right:     rtype,
		operation: operation,
		handler:   handler,
	})
}

func (os *OperationStore) RunOperation(lval, rval data.Value, operation string, compiler *Compiler, function *data.Function, class *data.Class) data.Value {

	for _, v := range os.operations {
		if operation == v.operation && matchOpdef(lval, v.left) && matchOpdef(rval, v.right) {
			//if the types match with the operation
			return v.handler(lval, rval, compiler, function, class)
		}
	}

	//there isn't an operation matching the given types
	return nil
}
