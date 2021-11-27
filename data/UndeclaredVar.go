package data

import (
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

//this is a type used when a variable in compilation was not declared
/*
	this has very limited usage
	- dot operator right hand operands

(thats it for now)
*/
type UndeclaredVar struct {
	Name string
}

func NewUndeclaredVar(name string) *UndeclaredVar {
	return &UndeclaredVar{
		Name: name,
	}
}

func (uv *UndeclaredVar) LLVal(function *Function) value.Value {
	return nil
}

func (uv *UndeclaredVar) TType() Type {
	return NewInvalidType()
}

func (uv *UndeclaredVar) Type() types.Type {
	return nil
}

func (uv *UndeclaredVar) TypeData() *TypeData {
	return NewTypeData("udvar")
}

func (uv *UndeclaredVar) InstanceV() value.Value {
	return nil
}
