package data

import (
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/value"

	"github.com/llir/llvm/ir/types"
)

type Integer struct {
	value *constant.Int

	//untyped
	untyped bool
	UTypVal int64
}

func NewInteger(i *constant.Int) *Integer {
	return &Integer{
		value: i,
	}
}

func NewUntypedInteger(v int64) *Integer {
	return &Integer{
		untyped: true,
		UTypVal: v,
	}
}

func (i *Integer) GetInt() int64 {
	if i.untyped {
		return i.UTypVal
	}

	return i.value.X.Int64()
}

func (i *Integer) LLVal(function *Function) value.Value {

	if i.untyped {
		return constant.NewInt(types.I32, i.UTypVal)
	}

	return i.value
}

func (i *Integer) TType() Type {

	if i.untyped {
		return NewUntypeIntType(i)
	}

	return NewPrimitive(i.Type())
}

func (i *Integer) Type() types.Type {

	if i.untyped {
		return nil
	}

	return i.value.Type()
}

func (i *Integer) TypeData() *TypeData {
	return i.TType().TypeData()
}

func (i *Integer) InstanceV() value.Value {
	return nil
}
