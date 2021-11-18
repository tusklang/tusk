package compiler

import (
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/value"
	"github.com/tusklang/tusk/ast"
	"github.com/tusklang/tusk/data"
)

var (
	inttypes   = []string{"i8", "i16", "i32", "i64", "i128"}
	uinttypes  = []string{"u8", "u16", "u32", "u64", "u128"}
	floattypes = []string{"f32", "f64"}
)

func addCastArray(compiler *ast.Compiler, typArr []string, fromType string, fn func(tname string, fromData data.Value, compiler *ast.Compiler, block *ir.Block, class *data.Class) data.Value) {
	for _, _v := range typArr {
		v := _v
		compiler.CastStore.NewCast(true, v, fromType, func(fromData data.Value, compiler *ast.Compiler, block *ir.Block, class *data.Class) data.Value {
			return fn(v, fromData, compiler, block, class)
		})
	}
}

func addXCasts2(auto, slice bool, compiler *ast.Compiler, fromArr []string, toArr []string, fn func(fromData data.Value, compiler *ast.Compiler, block *ir.Block, class *data.Class, typ types.Type) value.Value) {
	for k, _v := range fromArr {
		var v = _v

		sl := 0

		if slice {
			sl = k
		}

		for _, _vv := range toArr[sl:] {
			var vv = _vv

			if v == vv {
				continue
			}

			compiler.CastStore.NewCast(auto, vv, v, func(fromData data.Value, compiler *ast.Compiler, block *ir.Block, class *data.Class) data.Value {
				return data.NewInstVariable(fn(fromData, compiler, block, class, numtypes[vv].Type()), numtypes[vv])
			})

		}
	}
}

func addXCasts(auto bool, compiler *ast.Compiler, typArr []string, fn func(fromData data.Value, compiler *ast.Compiler, block *ir.Block, class *data.Class, typ types.Type) value.Value) {
	addXCasts2(auto, true, compiler, typArr, typArr, fn) //just do this function, but the outer and inner loops params are equal
}

//reverse a string array (type arrays)
func reverseStrArr(a []string) []string {
	var fin = make([]string, len(a))
	for k, v := range a {
		fin[len(a)-k-1] = v
	}
	return fin
}

func initDefaultCasts(compiler *ast.Compiler) {
	compiler.CastStore = ast.NewCastStore()

	//add upcasts
	addXCasts(true, compiler, inttypes, func(fromData data.Value, compiler *ast.Compiler, block *ir.Block, class *data.Class, typ types.Type) value.Value {
		return block.NewSExt(fromData.LLVal(block), typ)
	})
	addXCasts(true, compiler, uinttypes, func(fromData data.Value, compiler *ast.Compiler, block *ir.Block, class *data.Class, typ types.Type) value.Value {
		return block.NewZExt(fromData.LLVal(block), typ)
	})
	addXCasts(true, compiler, floattypes, func(fromData data.Value, compiler *ast.Compiler, block *ir.Block, class *data.Class, typ types.Type) value.Value {
		return block.NewFPExt(fromData.LLVal(block), typ)
	})

	//add downcasts
	addXCasts(false, compiler, reverseStrArr(inttypes), func(fromData data.Value, compiler *ast.Compiler, block *ir.Block, class *data.Class, typ types.Type) value.Value {
		return block.NewTrunc(fromData.LLVal(block), typ)
	})
	addXCasts(false, compiler, reverseStrArr(uinttypes), func(fromData data.Value, compiler *ast.Compiler, block *ir.Block, class *data.Class, typ types.Type) value.Value {
		return block.NewTrunc(fromData.LLVal(block), typ)
	})
	addXCasts(false, compiler, reverseStrArr(floattypes), func(fromData data.Value, compiler *ast.Compiler, block *ir.Block, class *data.Class, typ types.Type) value.Value {
		return block.NewFPTrunc(fromData.LLVal(block), typ)
	})

	//add casts between int/uint/float types
	addXCasts2(false, true, compiler, inttypes, uinttypes, func(fromData data.Value, compiler *ast.Compiler, block *ir.Block, class *data.Class, typ types.Type) value.Value {
		return block.NewSExt(fromData.LLVal(block), typ)
	})
	addXCasts2(false, false, compiler, inttypes, floattypes, func(fromData data.Value, compiler *ast.Compiler, block *ir.Block, class *data.Class, typ types.Type) value.Value {
		return block.NewSIToFP(fromData.LLVal(block), typ)
	})
	addXCasts2(false, true, compiler, uinttypes, inttypes, func(fromData data.Value, compiler *ast.Compiler, block *ir.Block, class *data.Class, typ types.Type) value.Value {
		return block.NewZExt(fromData.LLVal(block), typ)
	})
	addXCasts2(false, false, compiler, uinttypes, floattypes, func(fromData data.Value, compiler *ast.Compiler, block *ir.Block, class *data.Class, typ types.Type) value.Value {
		return block.NewUIToFP(fromData.LLVal(block), typ)
	})
	addXCasts2(false, false, compiler, floattypes, inttypes, func(fromData data.Value, compiler *ast.Compiler, block *ir.Block, class *data.Class, typ types.Type) value.Value {
		return block.NewFPToSI(fromData.LLVal(block), typ)
	})
	addXCasts2(false, false, compiler, floattypes, uinttypes, func(fromData data.Value, compiler *ast.Compiler, block *ir.Block, class *data.Class, typ types.Type) value.Value {
		return block.NewFPToUI(fromData.LLVal(block), typ)
	})

	//and also for downcasts
	addXCasts2(false, true, compiler, reverseStrArr(inttypes), reverseStrArr(uinttypes), func(fromData data.Value, compiler *ast.Compiler, block *ir.Block, class *data.Class, typ types.Type) value.Value {
		return block.NewTrunc(fromData.LLVal(block), typ)
	})
	addXCasts2(false, true, compiler, reverseStrArr(uinttypes), reverseStrArr(inttypes), func(fromData data.Value, compiler *ast.Compiler, block *ir.Block, class *data.Class, typ types.Type) value.Value {
		return block.NewTrunc(fromData.LLVal(block), typ)
	})

	//add casts from untyped numeric vals
	addCastArray(compiler, append(inttypes, uinttypes...), "untypedint", func(tname string, fromData data.Value, compiler *ast.Compiler, block *ir.Block, class *data.Class) data.Value {
		return data.NewInstVariable(constant.NewInt(numtypes[tname].Type().(*types.IntType), fromData.(*data.Integer).UTypVal), numtypes[tname])
	})

	addCastArray(compiler, floattypes, "untypedint", func(tname string, fromData data.Value, compiler *ast.Compiler, block *ir.Block, class *data.Class) data.Value {
		return data.NewInstVariable(constant.NewFloat(numtypes[tname].Type().(*types.FloatType), float64(fromData.(*data.Integer).UTypVal)), numtypes[tname])
	})

	addCastArray(compiler, append(inttypes, uinttypes...), "untypedfloat", func(tname string, fromData data.Value, compiler *ast.Compiler, block *ir.Block, class *data.Class) data.Value {
		return data.NewInstVariable(constant.NewInt(numtypes[tname].Type().(*types.IntType), int64(fromData.(*data.Float).UTypVal)), numtypes[tname])
	})

	addCastArray(compiler, floattypes, "untypedfloat", func(tname string, fromData data.Value, compiler *ast.Compiler, block *ir.Block, class *data.Class) data.Value {
		return data.NewInstVariable(constant.NewFloat(numtypes[tname].Type().(*types.FloatType), fromData.(*data.Float).UTypVal), numtypes[tname])
	})

}
