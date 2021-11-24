package data

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

type InstanceVariable struct {
	variable *Variable
	object   value.Value
}

func NewInstanceVariable(v *Variable, o value.Value) *InstanceVariable {
	return &InstanceVariable{
		variable: v,
		object:   o,
	}
}

func (v *InstanceVariable) LLVal(block *ir.Block) value.Value {
	return v.variable.LLVal(block)
}

func (v *InstanceVariable) TType() Type {
	return v.variable.TType()
}

func (v *InstanceVariable) Type() types.Type {
	return v.TType().Type()
}

func (v *InstanceVariable) TypeData() *TypeData {

	td := *v.variable.TypeData()
	td.AddFlag("instancevar")

	return &td
}

func (v *InstanceVariable) FetchAssig() value.Value {
	return v.variable.FetchAssig()
}

func (v *InstanceVariable) GetObj() value.Value {
	return v.object
}

func (v *InstanceVariable) InstanceV() value.Value {
	return v.object
}
