package data

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

type Pointer struct {
	typ    Type
	isType bool
}

func NewPointer(typ Type) *Pointer {
	return &Pointer{
		typ: typ,
	}
}

func (p *Pointer) SetToType() {
	p.isType = true //it is now a pointer type
}

func (p *Pointer) PType() Type {
	return p.typ
}

func (p *Pointer) TType() Type {
	return p
}

func (p *Pointer) Type() types.Type {
	return types.NewPointer(p.typ.Type())
}

func (p *Pointer) LLVal(block *ir.Block) value.Value {
	return nil
}

func (p *Pointer) TValue() Value {
	return p
}

func (p *Pointer) TypeData() *TypeData {
	td := *p.typ.TypeData()
	td.AddFlag("ptr")

	if p.isType {
		td.AddFlag("type")
	}

	return &td
}

func (p *Pointer) InstanceV() value.Value {
	return nil
}

func (p *Pointer) Equals(t Type) bool {
	return p.Type().Equal(t.Type())
}

func (p *Pointer) Default() constant.Constant {
	return constant.NewNull(p.Type().(*types.PointerType))
}

func (p *Pointer) Alignment() uint64 {
	return 8
}
