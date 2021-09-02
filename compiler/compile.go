package compiler

import (
	"runtime"

	"github.com/llir/llvm/ir"
	"github.com/tusklang/tusk/ast"
	"github.com/tusklang/tusk/initialize"
)

func Compile(prog *initialize.Program, outfile string) {
	var compiler = &ast.Compiler{
		Module: ir.NewModule(),
		OS:     runtime.GOOS,
		ARCH:   runtime.GOARCH,
	}

	_ = compiler

	for _, v := range prog.Packages {
		for _, vv := range v.Files {
			_ = vv
		}
	}
}
