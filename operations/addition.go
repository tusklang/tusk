package operations

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/value"
	"github.com/tusklang/tusk/ast"
)

func initIntOps(compiler *ast.Compiler) {

	compiler.Operations["i64 + i64"] = func(s1, s2 value.Value, block *ir.Block) value.Value {
		return block.NewAdd(s1, s2)
	}
	compiler.Operations["i32 + i32"] = func(s1, s2 value.Value, block *ir.Block) value.Value {
		return block.NewAdd(s1, s2)
	}
	compiler.Operations["i16 + i16"] = func(s1, s2 value.Value, block *ir.Block) value.Value {
		return block.NewAdd(s1, s2)
	}
	compiler.Operations["i8 + i8"] = func(s1, s2 value.Value, block *ir.Block) value.Value {
		return block.NewAdd(s1, s2)
	}

	compiler.Operations["f64 + f64"] = func(s1, s2 value.Value, block *ir.Block) value.Value {
		return block.NewAdd(s1, s2)
	}
	compiler.Operations["f32 + f32"] = func(s1, s2 value.Value, block *ir.Block) value.Value {
		return block.NewAdd(s1, s2)
	}

}
