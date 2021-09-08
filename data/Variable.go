package data

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

type Variable struct {
	val value.Value
	typ types.Type
}

func NewVariable(val value.Value, typ types.Type) *Variable {
	return &Variable{
		val: val,
		typ: typ,
	}
}

func (v *Variable) LLVal(block *ir.Block) value.Value {
	return block.NewLoad(v.typ, v.val)
}

func (v *Variable) Type() types.Type {
	return v.typ
}

func (v *Variable) TypeString() string {
	return v.Type().LLString()
}
