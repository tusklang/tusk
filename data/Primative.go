package data

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

type Primitive struct {
	typ types.Type
	nam string
}

func NewPrimitive(typ types.Type) *Primitive {
	return &Primitive{
		typ: typ,
		nam: typ.LLString(),
	}
}

func NewNamedPrimitive(typ types.Type, nam string) *Primitive {
	p := NewPrimitive(typ)
	p.nam = nam
	return p
}

func (p *Primitive) Default() constant.Constant {
	switch p.typ {
	case types.I32:
		return constant.NewInt(types.I32, 0)
	case types.Float:
		return constant.NewFloat(types.Float, 0)
	default:
		return &constant.Null{}
	}
}

func (p *Primitive) LLVal(block *ir.Block) value.Value {
	return nil
}

func (p *Primitive) TType() Type {
	return p
}

func (p *Primitive) Type() types.Type {
	return p.typ
}

func (p *Primitive) TypeData() *TypeData {
	td := NewTypeData(p.nam)
	td.AddFlag("type")
	return td
}

func (p *Primitive) InstanceV() value.Value {
	return nil
}

func (p *Primitive) Equals(t Type) bool {
	return p.TypeData().Name() == t.TypeData().Name()
}
