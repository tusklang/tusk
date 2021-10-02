package data

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

type String struct {
	CharArray     []byte
	gd            *ir.Global
	newTuskString *ir.Func //function to init a new string struct type
}

func NewString(s []byte) *String {
	return &String{
		CharArray: s,
	}
}

func (s *String) Init(gd *ir.Global, ns *ir.Func) {
	s.gd = gd
	s.newTuskString = ns
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
	return block.NewCall(s.newTuskString, gep, constant.NewInt(types.I32, int64(len(s.CharArray))))
}

func (s *String) TType() Type {
	return NewPrimitive(types.I8Ptr)
}

func (s *String) Type() types.Type {
	return types.I8Ptr
}

func (s *String) TypeData() *TypeData {
	return NewTypeData("string")
}
