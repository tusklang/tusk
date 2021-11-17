package compiler

import (
	"github.com/llir/llvm/ir"
	"github.com/tusklang/tusk/ast"
	"github.com/tusklang/tusk/data"
)

func initDefaultCasts(compiler *ast.Compiler) {
	compiler.CastStore = ast.NewCastStore()

	compiler.CastStore.NewCast(true, "i64", "i32", func(fromData data.Value, compiler *ast.Compiler, block *ir.Block, class *data.Class) data.Value {
		return data.NewInstVariable(block.NewSExt(fromData.LLVal(block), numtypes["i64"].Type()), numtypes["i64"])
	})
}
