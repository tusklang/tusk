package ast

import (
	"strconv"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

type Compiler struct {
	Module        *ir.Module                                                       //llvm module
	ValidTypes    map[string]types.Type                                            //a list of valid types
	tmpVarCnt     uint64                                                           //integer ID to use for temp vars (each use results in increment)
	OS, ARCH      string                                                           //operating system and architecture
	Operations    map[string]func(value.Value, value.Value, *ir.Block) value.Value //map of all operations
	StaticGlobals map[string]*ir.Global                                            //all global static variables
	InitBlock     *ir.Block                                                        //function that runs before main (used to initialize globals and stuff)
}

func (c *Compiler) TmpVar() string {
	c.tmpVarCnt++
	return "tv_" + strconv.FormatUint(c.tmpVarCnt, 10)
}
