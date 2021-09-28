package data

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

type Variable struct {
	inst value.Value
	typ  Type
}

func NewVariable(inst value.Value, typ Type) *Variable {
	return &Variable{
		inst: inst,
		typ:  typ,
	}
}

func (v *Variable) FetchAssig() value.Value {
	return v.inst
}

func (v *Variable) LLVal(block *ir.Block) value.Value {
	return block.NewLoad(v.typ.Type(), v.inst)
}

func (v *Variable) TType() Type {
	return v.typ
}

func (v *Variable) Type() types.Type {
	return v.TType().Type()
}

func (v *Variable) TypeString() string {
	return v.TType().TypeString()
}
