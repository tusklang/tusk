package ast

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/tusklang/tusk/data"
	"github.com/tusklang/tusk/errhandle"
)

var Booltype = data.NewNamedPrimitive(types.I1, "bool")

var InttypeV = map[string]data.Type{
	"i128": data.NewPrimitive(types.I128),
	"i64":  data.NewPrimitive(types.I64),
	"i32":  data.NewPrimitive(types.I32),
	"i16":  data.NewPrimitive(types.I16),
	"i8":   data.NewPrimitive(types.I8),
}

var UinttypeV = map[string]data.Type{
	"u128": data.NewNamedPrimitive(types.I128, "u128"),
	"u64":  data.NewNamedPrimitive(types.I64, "u64"),
	"u32":  data.NewNamedPrimitive(types.I32, "u32"),
	"u16":  data.NewNamedPrimitive(types.I16, "u16"),
	"u8":   data.NewNamedPrimitive(types.I8, "u8"),
}

var FloattypeV = map[string]data.Type{
	"f64": data.NewNamedPrimitive(types.Double, "f64"),
	"f32": data.NewNamedPrimitive(types.Float, "f32"),
}

//list of all the numerical types
var Numtypes = map[string]data.Type{}

type Compiler struct {
	Module          *ir.Module             //llvm module
	OS, ARCH        string                 //operating system and architecture
	InitFunc        *data.Function         //function that runs before main (used to initialize globals and stuff)
	OperationStore  *OperationStore        //list of all operations
	CastStore       *CastStore             //list of all typecasts
	VarMap          map[string]data.Value  //map of all the variables declared
	LinkedFunctions map[string]*ir.Func    //map of all the linked functions used (e.g. glibc functions)
	Errors          []*errhandle.TuskError //list of all errors generated while compiling a program
	lambdacnt       uint64                 //amount of lambdas used in the program
}

func (c *Compiler) AddVar(name string, v data.Value) {
	c.VarMap[name] = v
}

func (c *Compiler) FetchVar(name string) data.Value {
	return c.VarMap[name]
}

func (c *Compiler) AddError(e *errhandle.TuskError) {
	c.Errors = append(c.Errors, e)
}

func (c *Compiler) LambdaInc() uint64 {
	c.lambdacnt++
	return c.lambdacnt
}
