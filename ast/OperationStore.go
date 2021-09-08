package ast

import (
	"github.com/llir/llvm/ir"
	"github.com/tusklang/tusk/data"
)

type operationdef struct {
	left, right string
	operation   string
	handler     func(left, right data.Value, compiler *Compiler, block *ir.Block) data.Value
}

type OperationStore struct {
	operations []operationdef
}

func NewOperationStore() *OperationStore {
	return &OperationStore{}
}

func (os *OperationStore) NewOperation(operation string, ltype, rtype string, handler func(left, right data.Value, compiler *Compiler, block *ir.Block) data.Value) {
	os.operations = append(os.operations, operationdef{
		left:      ltype,
		right:     rtype,
		operation: operation,
		handler:   handler,
	})
}

func (os *OperationStore) RunOperation(lval, rval data.Value, operation string, compiler *Compiler, block *ir.Block) data.Value {

	var (
		ltyp = lval.TypeString()
		rtyp = rval.TypeString()
	)

	for _, v := range os.operations {
		if operation == v.operation && matchOpdef(ltyp, v.left) && matchOpdef(rtyp, v.right) {
			//if the types match with the operation
			return v.handler(lval, rval, compiler, block)
		}
	}

	//there isn't an operation matching the given types
	return nil
}
