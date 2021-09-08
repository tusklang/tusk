package ast

import (
	"strconv"

	"github.com/llir/llvm/ir"
	"github.com/tusklang/tusk/data"
)

type Compiler struct {
	Module         *ir.Module                //llvm module
	tmpVarCnt      uint64                    //integer ID to use for temp vars (each use results in increment)
	OS, ARCH       string                    //operating system and architecture
	InitBlock      *ir.Block                 //function that runs before main (used to initialize globals and stuff)
	OperationStore *OperationStore           //list of all operations
	VarMap         map[string]*data.Variable //map of all the variables declared
}

func (c *Compiler) TmpVar() string {
	c.tmpVarCnt++
	return "tv_" + strconv.FormatUint(c.tmpVarCnt, 10)
}

func (c *Compiler) AddVar(name string, v *data.Variable) {
	c.VarMap[name] = v
}

func (c *Compiler) FetchVar(name string) *data.Variable {
	return c.VarMap[name]
}
