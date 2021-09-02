package ast

import "github.com/llir/llvm/ir"

type Compiler struct {
	Module   *ir.Module //llvm mdoule
	OS, ARCH string     //operating system and architecture
}
