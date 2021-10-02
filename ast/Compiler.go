package ast

import (
	"github.com/llir/llvm/ir"
	"github.com/tusklang/tusk/data"
)

type Compiler struct {
	Module         *ir.Module            //llvm module
	OS, ARCH       string                //operating system and architecture
	InitBlock      *ir.Block             //function that runs before main (used to initialize globals and stuff)
	OperationStore *OperationStore       //list of all operations
	VarMap         map[string]data.Value //map of all the variables declared
	NewString      *ir.Func              //create a new string class in tusk
}

func (c *Compiler) AddVar(name string, v data.Value) {
	c.VarMap[name] = v
}

func (c *Compiler) FetchVar(name string) data.Value {
	return c.VarMap[name]
}
