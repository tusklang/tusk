package data

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

type Class struct {
	Name     string
	SType    *types.StructType
	Instance map[string]*Variable
	Static   map[string]*Variable

	ParentPackage *Package
}

func NewClass(name string, st *types.StructType, parent *Package) *Class {
	return &Class{
		Name:          name,
		SType:         st,
		ParentPackage: parent,
	}
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
