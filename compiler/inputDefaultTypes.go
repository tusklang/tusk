package compiler

import (
	"github.com/llir/llvm/ir/types"
	"github.com/tusklang/tusk/ast"
)

//all of the basic types
func inputDefaultTypes(compiler *ast.Compiler) {

	compiler.ValidTypes = make(map[string]types.Type)

	compiler.ValidTypes["i64"] = types.I64
	compiler.ValidTypes["i32"] = types.I32
	compiler.ValidTypes["i16"] = types.I16
	compiler.ValidTypes["i8"] = types.I8

	compiler.ValidTypes["f64"] = types.Double
	compiler.ValidTypes["f32"] = types.Float
}
