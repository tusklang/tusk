package data

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

type Variable struct {
	val             Value
	typ             *Type
	UnReferenceable bool
}

func NewVariable(val Value, typ *Type, unReferenceable bool) *Variable {
	return &Variable{
		val:             val,
		typ:             typ,
		UnReferenceable: unReferenceable,
	}
}

func (v *Variable) FetchVal() Value {
	return v.val
}

func (v *Variable) LLVal(block *ir.Block) value.Value {
	return block.NewLoad(v.typ.Type(), v.val.LLVal(block))
}

func (v *Variable) Type() types.Type {
	return v.typ.Type()
}

func (v *Variable) TypeString() string {
	return v.typ.TypeString()
}
