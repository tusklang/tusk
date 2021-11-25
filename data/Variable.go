package data

import (
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

type Variable struct {
	inst value.Value
	typ  Type

	loadinst func(*Variable, *Function) value.Value
}

func NewVariable(inst value.Value, typ Type) *Variable {
	return &Variable{
		inst: inst,
		typ:  typ,
		loadinst: func(v *Variable, function *Function) value.Value {
			return function.ActiveBlock.NewLoad(v.Type(), v.FetchAssig())
		},
	}
}

func NewInstVariable(inst value.Value, typ Type) *Variable {
	vd := NewVariable(inst, typ)
	vd.SetLoadInst(func(v *Variable, f *Function) value.Value {
		return v.FetchAssig()
	})
	return vd
}

func (v *Variable) FetchAssig() value.Value {
	return v.inst
}

func (v *Variable) SetLoadInst(f func(*Variable, *Function) value.Value) {
	v.loadinst = f
}

func (v *Variable) LLVal(function *Function) value.Value {
	return v.loadinst(v, function)
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
