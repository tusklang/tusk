package operations

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/value"
	"github.com/tusklang/tusk/ast"
)

func initDefaultCasts(compiler *ast.Compiler) {

	compiler.Operations["i64 -> i32"] = func(s1, s2 value.Value, block *ir.Block) value.Value {
		return block.NewSExt(s2, s1.Type())
	}

}
