package data

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

type ClassField struct {
	Index int64
	Type  Type

	/*
		access types
		0 - public
		1 - protected
		2 - private
	*/
	Access int
	Value  Value
}
type Class struct {
	Name           string
	SType          *types.StructType
	Instance       map[string]*ClassField
	Static         map[string]*ClassField
	Methods        map[string]*ClassField
	Construct      *Function
	ConstructAlloc value.Value

	ParentPackage *Package

	TypSiz uint64 //size of the type, in bytes

	curInstCnt int64 //current index we're on in the instance count (temporary while adding items to the instance map)
}

func NewClass(name string, st *types.StructType, parent *Package) *Class {
	return &Class{
		Name:          name,
		SType:         st,
		ParentPackage: parent,
		Instance:      make(map[string]*ClassField),
		Static:        make(map[string]*ClassField),
		Methods:       make(map[string]*ClassField),
	}
}

func (c *Class) AppendInstance(name string, typ Type, access int) {
	c.AddInstanceItem(name, typ, c.nextInstanceIdx(), access)
}

func (c *Class) AddInstanceItem(name string, typ Type, idx int64, access int) {
	c.Instance[name] = &ClassField{
		Index:  idx,
		Type:   typ,
		Access: access,
	}
}

func (c *Class) AppendStatic(name string, val Value, typ Type, access int) {
	c.Static[name] = &ClassField{
		Type:   typ,
		Value:  val,
		Access: access,
	}
}

func (c *Class) nextInstanceIdx() int64 {
	idx := c.curInstCnt
	c.curInstCnt++
	return idx
}

func (c *Class) NewMethod(name string, fn *Function, access int) {
	c.Methods[name] = &ClassField{
		Type:   fn.TType(),
		Access: access,
		Value:  fn,
	}
}

func (c *Class) LLVal(block *ir.Block) value.Value {
	return nil
}

func (c *Class) Default() constant.Constant {
	return constant.NewNull(types.NewPointer(c.SType))
}

func (c *Class) TType() Type {
	return NewInstance(c)
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

func (c *Class) InstanceV() value.Value {
	return nil
}

func (c *Class) Equals(t Type) bool {
	return c.Type().Equal(t.Type())
}

func (c *Class) FullName() string {
	return c.ParentPackage.FullName + "." + c.Name
}

func (c *Class) TypeSize() uint64 {
	return c.TypSiz
}
