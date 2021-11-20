package data

import (
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

type UntypeIntType struct {
	it *Integer
}

type UntypeFloatType struct {
	ft *Float
}

func NewUntypeIntType(it *Integer) *UntypeIntType {
	return &UntypeIntType{
		it: it,
	}
}

func (i *UntypeIntType) TType() Type {
	return i
}

func (i *UntypeIntType) Type() types.Type {
	return types.I32
}

func (i *UntypeIntType) TypeData() *TypeData {
	return NewTypeData("untypedint")
}

func (i *UntypeIntType) Default() constant.Constant {
	return constant.NewInt(types.I32, 0)
}

func (i *UntypeIntType) Equals(other Type) bool {
	return other.TypeData().Name() == i.TypeData().Name()
}

func (i *UntypeIntType) InstanceV() value.Value {
	return nil
}

func (i *UntypeIntType) TypeSize() uint64 {
	return 4
}

func NewUntypeFloatType(ft *Float) *UntypeFloatType {
	return &UntypeFloatType{
		ft: ft,
	}
}

func (f *UntypeFloatType) TType() Type {
	return f
}

func (f *UntypeFloatType) Type() types.Type {
	return types.Double
}

func (f *UntypeFloatType) TypeData() *TypeData {
	return NewTypeData("untypedfloat")
}

func (f *UntypeFloatType) Default() constant.Constant {
	return constant.NewFloat(types.Double, 0)
}

func (f *UntypeFloatType) Equals(other Type) bool {
	return other.TypeData().Name() == f.TypeData().Name()
}

func (f *UntypeFloatType) InstanceV() value.Value {
	return nil
}

func (f *UntypeFloatType) TypeSize() uint64 {
	return 8
}
