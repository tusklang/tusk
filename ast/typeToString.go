package ast

import (
	"github.com/llir/llvm/ir/types"
)

func typeToString(t types.Type) string {
	switch t.String() {
	case "float":
		return "f32"
	case "double":
		return "f64"
	default:
		return t.String()
	}
}
