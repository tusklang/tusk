package ast

import (
	"github.com/llir/llvm/ir"
	"github.com/tusklang/tusk/data"
)

type castdef struct {
	toType, fromType string
	handler          func(fromData data.Value, compiler *Compiler, block *ir.Block, class *data.Class) data.Value
	auto             bool
}

type CastStore struct {
	casts []castdef
}

func NewCastStore() *CastStore {
	return &CastStore{}
}

func (cs *CastStore) NewCast(auto bool, toType string, fromType string, handler func(fromData data.Value, compiler *Compiler, block *ir.Block, class *data.Class) data.Value) {
	cs.casts = append(cs.casts, castdef{
		auto:     auto,
		toType:   toType,
		fromType: fromType,
		handler:  handler,
	})
}

func (cs *CastStore) RunCast(auto bool, toType string, fromData data.Value, compiler *Compiler, block *ir.Block, class *data.Class) data.Value {

	for _, v := range cs.casts {
		if toType == v.toType && fromData.TypeData().Name() == v.fromType && (!auto || v.auto) {
			return v.handler(fromData, compiler, block, class)
		}
	}

	//there isn't an operation matching the given types
	return nil
}
