package data

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

type String struct {
	CharArray []byte
	gd        *ir.Global
}

func NewString(s []byte) *String {
	return &String{
		CharArray: s,
	}
}

func (s *String) Init(gd *ir.Global) {
	s.gd = gd
}

func (s *String) LLVal(block *ir.Block) value.Value {

	//https://github.com/nektro/slate/blob/master/pgk/parse/llvm/llvm.go#L49
	//took a lot from there >_>
	gep := constant.NewGetElementPtr(
		types.NewArray(uint64(len(s.CharArray)), types.I8),
		s.gd,
		constant.NewInt(types.I32, 0),
		constant.NewInt(types.I32, 0),
	)
	gep.InBounds = true
	return gep
}

func (s *String) TType() Type {
	return NewPointer(NewPrimitive(types.I8))
}

func (s *String) Type() types.Type {
	return types.I8Ptr
}

func (s *String) TypeData() *TypeData {
	return s.TType().TypeData()
}

func (s *String) InstanceV() value.Value {
	return nil
}
