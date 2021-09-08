package data

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

type Type struct {
	typ types.Type
}

func NewType(typ types.Type) *Type {
	return &Type{
		typ: typ,
	}
}

func (t *Type) Default() constant.Constant {
	switch t.typ {
	case types.I32:
		return constant.NewInt(types.I32, 0)
	default:
		return &constant.Null{}
	}
}

func (t *Type) LLVal(block *ir.Block) value.Value {
	return nil
}

func (t *Type) Type() types.Type {
	return t.typ
}

func (t *Type) TypeString() string {
	return t.Type().LLString()
}
