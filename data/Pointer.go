package data

import (
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
)

type Pointer struct {
	typ Type
}

func NewPointer(typ Type) *Pointer {
	return &Pointer{
		typ: typ,
	}
}

func (p *Pointer) TType() Type {
	return p.typ
}

func (p *Pointer) Type() types.Type {
	return types.NewPointer(p.typ.Type())
}

func (p *Pointer) TypeData() *TypeData {
	td := *p.typ.TypeData()
	td.AddFlag("ptr")
	return &td
}

func (p *Pointer) Default() constant.Constant {
	return constant.NewNull(p.Type().(*types.PointerType))
}
