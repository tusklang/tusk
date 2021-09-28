package data

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

type Null struct{}

func NewNull() *Null {
	return &Null{}
}

func (n *Null) LLVal(block *ir.Block) value.Value {
	return nil
}

func (n *Null) TType() Type {
	return nil
}

func (n *Null) Type() types.Type {
	return nil
}

func (n *Null) TypeString() string {
	return "null"
}
