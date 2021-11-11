package ast

import (
	"github.com/llir/llvm/ir"
	"github.com/tusklang/tusk/data"
)

type Compiler struct {
	Module          *ir.Module            //llvm module
	OS, ARCH        string                //operating system and architecture
	InitFunc        *data.Function        //function that runs before main (used to initialize globals and stuff)
	OperationStore  *OperationStore       //list of all operations
	VarMap          map[string]data.Value //map of all the variables declared
	LinkedFunctions map[string]*ir.Func   //map of all the linked functions used (e.g. glibc functions)
}

func (c *Compiler) AddVar(name string, v data.Value) {
	c.VarMap[name] = v
}

func (c *Compiler) FetchVar(name string) data.Value {
	return c.VarMap[name]
}
