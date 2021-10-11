package compiler

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
)

type builtin interface {
	FetchName() string
}

type builtinClass struct {
	exported map[int64]string
	name     string

	structT     *types.StructType
	constructor *ir.Func

	helpers []*ir.Func

	nomangle []string
}

func newBuiltinClass() *builtinClass {
	return &builtinClass{
		exported: make(map[int64]string),
	}
}

func (c *builtinClass) FetchName() string {
	return c.name
}

type builtinFunc struct {
	name string

	exported *ir.Func

	helpers  []*ir.Func
	nomangle []string
}

func newBuiltinFunc() *builtinFunc {
	return &builtinFunc{}
}

func (f *builtinFunc) FetchName() string {
	return f.name
}
