package compiler

import (
	"github.com/llir/llvm/ir/types"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/value"
	"github.com/tusklang/tusk/ast"
	"github.com/tusklang/tusk/data"
)

var (
	inttypes   = []string{"bool", "i8", "i16", "i32", "i64", "i128"}
	uinttypes  = []string{"u8", "u16", "u32", "u64", "u128"}
	floattypes = []string{"f32", "f64"}
)

func addUpcasts(compiler *ast.Compiler, typArr []string, fn func(fromData data.Value, compiler *ast.Compiler, block *ir.Block, class *data.Class, typ types.Type) value.Value) {
	for k, _v := range typArr {
		var v = _v
		for _, _vv := range typArr[k:] {
			var vv = _vv

			if v == vv {
				continue
			}

			compiler.CastStore.NewCast(true, vv, v, func(fromData data.Value, compiler *ast.Compiler, block *ir.Block, class *data.Class) data.Value {
				return data.NewInstVariable(fn(fromData, compiler, block, class, numtypes[vv].Type()), numtypes[vv])
			})

		}
	}
}

func initDefaultCasts(compiler *ast.Compiler) {
	compiler.CastStore = ast.NewCastStore()

	addUpcasts(compiler, inttypes, func(fromData data.Value, compiler *ast.Compiler, block *ir.Block, class *data.Class, typ types.Type) value.Value {
		return block.NewSExt(fromData.LLVal(block), typ)
	})
	addUpcasts(compiler, uinttypes, func(fromData data.Value, compiler *ast.Compiler, block *ir.Block, class *data.Class, typ types.Type) value.Value {
		return block.NewZExt(fromData.LLVal(block), typ)
	})
	addUpcasts(compiler, floattypes, func(fromData data.Value, compiler *ast.Compiler, block *ir.Block, class *data.Class, typ types.Type) value.Value {
		return block.NewFPExt(fromData.LLVal(block), typ)
	})
}
