package compiler

import "github.com/llir/llvm/ir"

type Compiler struct {
	module   *ir.Module
	OS, ARCH string //operating system and architecture
}
