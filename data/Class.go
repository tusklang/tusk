package data

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

type InstanceVar struct {
	Index int64
	Type  Type
}
type Class struct {
	Name           string
	SType          *types.StructType
	Instance       map[string]*InstanceVar
	Static         map[string]*Variable
	Construct      *Function
	ConstructAlloc value.Value

	ParentPackage *Package

	curInstCnt int64 //current index we're on in the instance count (temporary while adding items to the instance map)
}

func NewClass(name string, st *types.StructType, parent *Package) *Class {
	return &Class{
		Name:          name,
		SType:         st,
		ParentPackage: parent,
		Instance:      make(map[string]*InstanceVar),
		Static:        make(map[string]*Variable),
	}
}

func (c *Class) AppendInstance(name string, typ Type) {
	c.AddInstanceItem(name, typ, c.nextInstanceIdx())
}

func (c *Class) AddInstanceItem(name string, typ Type, idx int64) {
	c.Instance[name] = &InstanceVar{
		Index: idx,
		Type:  typ,
	}
}

func (c *Class) nextInstanceIdx() int64 {
	idx := c.curInstCnt
	c.curInstCnt++
	return idx
}

func (c *Class) LLVal(block *ir.Block) value.Value {
	return nil
}

func (c *Class) Default() constant.Constant {
	return constant.NewNull(types.NewPointer(c.SType))
}

func (c *Class) TType() Type {
	return c
}

func (c *Class) Type() types.Type {
	return types.NewPointer(c.SType)
}

func (c *Class) TypeData() *TypeData {

	td := NewTypeData(c.Name)
	td.AddFlag("class")
	td.AddFlag("type")

	return td
}

func (c *Class) Equals(t Type) bool {
	return c.Type().Equal(t.Type())
}
