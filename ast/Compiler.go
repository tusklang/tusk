package ast

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
)

type Compiler struct {
	Module     *ir.Module
	ValidTypes map[string]types.Type //a list of valid types
	OS, ARCH   string                //operating system and architecture
}
