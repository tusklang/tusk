package data

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

type Class struct {
	Name     string
	SType    *types.StructType
	Instance []*Variable
	Static   []*Variable
}

func NewClass(st *types.StructType) *Class {
	return nil
}

func (c *Class) LLVal(block *ir.Block) value.Value {
	return nil
}

func (c *Class) Type() types.Type {
	return c.SType
}

func (c *Class) TypeString() string {
	return c.Type().LLString()
}
