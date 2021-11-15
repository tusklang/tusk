package data

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

type Method struct {
	Func     Value
	Instance value.Value
}

func NewMethod(f Value, i value.Value) *Method {
	return &Method{
		Func:     f,
		Instance: i,
	}
}

func (m *Method) LLVal(block *ir.Block) value.Value {
	return m.Func.LLVal(block)
}

func (m *Method) TType() Type {
	return m
}

func (m *Method) Type() types.Type {
	return m.Func.Type()
}

func (m *Method) TypeData() *TypeData {
	d := *m.Func.TypeData()
	d.AddFlag("method")
	return &d
}

func (m *Method) InstanceV() value.Value {
	return m.Instance
}

func (m *Method) Equals(t Type) bool {
	return m.Func.Type().Equal(t.Type())
}

func (m *Method) Default() constant.Constant {
	return nil
}

func (m *Method) TypeSize() uint64 {
	return 8
}
