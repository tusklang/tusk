package ast

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
)

type Compiler struct {
	Module        *ir.Module
	ValidTypes    map[string]types.Type //a list of valid types
	tmpVarCnt     int                   //integer ID to use for temp vars (each use results in increment)
	OS, ARCH      string                //operating system and architecture
	StaticGlobals map[string]*ir.Global
}

func (c *Compiler) TmpVar() int {
	c.tmpVarCnt++
	return c.tmpVarCnt
}
