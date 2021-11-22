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

	switch g := p.typ.(type) {
	case *types.IntType:
		return constant.NewInt(g, 0)
	case *types.FloatType:
		return constant.NewFloat(g, 0)
	}

	return &constant.Null{}
}

func (p *Primitive) LLVal(block *ir.Block) value.Value {
	return nil
}

func (p *Primitive) TValue() Value {
	return p
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

func (p *Primitive) Alignment() uint64 {
	switch p.typ {
	case types.I64:
		return 8
	case types.I32:
		return 4
	case types.I16:
		return 2
	case types.I8:
		return 1
	case types.I1:
		return 1
	case types.Double:
		return 8
	case types.Float:
		return 4
	default:
		return 8
	}
}
