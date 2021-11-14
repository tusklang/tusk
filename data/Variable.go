package data

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

type Variable struct {
	inst value.Value
	typ  Type

	loadinst func(*Variable, *ir.Block) value.Value
}

func NewVariable(inst value.Value, typ Type) *Variable {
	return &Variable{
		inst: inst,
		typ:  typ,
		loadinst: func(v *Variable, block *ir.Block) value.Value {
			return block.NewLoad(v.Type(), v.FetchAssig())
		},
	}
}

func NewInstVariable(inst value.Value, typ Type) *Variable {
	vd := NewVariable(inst, typ)
	vd.SetLoadInst(func(v *Variable, b *ir.Block) value.Value {
		return v.FetchAssig()
	})
	return vd
}

func (v *Variable) FetchAssig() value.Value {
	return v.inst
}

func (v *Variable) SetLoadInst(f func(*Variable, *ir.Block) value.Value) {
	v.loadinst = f
}

func (v *Variable) LLVal(block *ir.Block) value.Value {
	return v.loadinst(v, block)
}

func (v *Variable) TType() Type {
	return v.typ.TType()
}

func (v *Variable) Type() types.Type {
	return v.TType().Type()
}

func (v *Variable) TypeData() *TypeData {

	td := *v.typ.TypeData()
	td.AddFlag("var")

	return &td
}

func (v *Variable) InstanceV() value.Value {
	return v.typ.InstanceV()
}
