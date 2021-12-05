package ast

import (
	"fmt"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	"github.com/tusklang/tusk/data"
	"github.com/tusklang/tusk/errhandle"
	"github.com/tusklang/tusk/tokenizer"
)

type VarDecl struct {
	Name  string
	Type  *ASTNode
	Value *ASTNode

	//tokens in the decl
	kwtok,
	vnametok,
	typtok,
	valtok tokenizer.Token
	////////////////////

	//for globals
	declaration value.Value
	decltyp     data.Type
	globalerr   bool
}

func (vd *VarDecl) Parse(lex []tokenizer.Token, i *int, stopAt []string) *errhandle.TuskError {

	vd.kwtok = lex[*i]

	*i++

	if lex[*i].Type != "varname" {
		return errhandle.NewParseErrorFTok(
			"expected variable name",
			"",
			lex[*i],
		)
	}

	vd.vnametok = lex[*i]

	vd.Name = lex[*i].Name

	*i++

	//has a specified type
	if lex[*i].Name == ":" {
		*i++
		vd.typtok = lex[*i]
		tg, e := groupSpecific(lex, i, []string{"=", ";"}, -1)
		if e != nil {
			return e
		}
		t, e := groupsToAST(tg)
		if e != nil {
			return e
		}
		vd.Type = t[0]
	}

	//has a value assigned to it
	if lex[*i].Name == "=" {
		*i++
		vd.valtok = lex[*i]
		vgbm, e := braceMatcher(lex, i, allopeners, allclosers, false, "terminator")
		if e != nil {
			return e
		}
		vg, e := grouper(vgbm)
		if e != nil {
			return e
		}
		v, e := groupsToAST(vg)
		if e != nil {
			return e
		}
		vd.Value = v[0]
	}

	*i-- //the outer loop will incremenet for us

	return nil
}

func (vd *VarDecl) GetMTok() tokenizer.Token {
	return vd.kwtok
}

func (vd *VarDecl) getDeclType(compiler *Compiler, class *data.Class, function *data.Function) data.Type {

	//if there is no type supplied, return nil
	if vd.Type == nil {
		return nil
	}

	pvtype := vd.Type.Group.Compile(compiler, class, vd.Type, function)

	var vtype data.Type

	switch vt := pvtype.(type) {
	case *data.Class:
		vtype = data.NewInstance(vt)
	case *data.Primitive:
		vtype = vt
	case *data.Pointer:
		vtype = vt
	case *data.Function:
		vtype = vt
	case *data.FixedArray:
		vtype = vt
	case *data.SliceArray:
		vtype = vt
	case *data.VariedLengthArray:
		vtype = vt
	default:
		//error
		compiler.AddError(errhandle.NewCompileErrorFTok(
			"invalid type",
			"value cannot be used as a type",
			vd.typtok,
		))
	}

	return vtype
}

func (vd *VarDecl) Compile(compiler *Compiler, class *data.Class, node *ASTNode, function *data.Function) data.Value {

	var vtype data.Type
	var varval data.Value

	if vd.Value == nil && vd.Type != nil {
		//var varname: typename
		vtype = vd.getDeclType(compiler, class, function)
		if vtype == nil {
			//invalid type
			return nil
		}
		varval = data.NewInstVariable(vtype.Default(), vtype)
	} else if vd.Value != nil && vd.Type == nil {
		//var varname = value
		varval = vd.Value.Group.Compile(compiler, class, vd.Value, function)
		vtype = varval.TType()
	} else if vd.Value != nil && vd.Type != nil {
		//var varname: typename = value;
		vtype = vd.getDeclType(compiler, class, function)
		if vtype == nil {
			//invalid type
			return nil
		}
		varval = vd.Value.Group.Compile(compiler, class, vd.Value, function)
	} else {
		//var varname
		//error
		//cannot determine what type this variable is
		compiler.AddError(errhandle.NewCompileErrorFTok(
			"cannot infer type",
			"provide a type or a value for this variable",
			vd.vnametok,
		))
		return nil
	}

	switch varval.(type) {
	case *data.InvalidType:
		compiler.AddError(errhandle.NewCompileErrorFTok(
			"invalid value",
			"",
			vd.valtok,
		))
		return nil
	}

	//untyped values don't exist in llvm, so we force them to doubles/i32
	switch vtype.(type) {
	case *data.UntypeFloatType:
		vtype = data.NewNamedPrimitive(types.Double, "f64")
	case *data.UntypeIntType:
		vtype = data.NewPrimitive(types.I32)
	}

	decl := function.ActiveBlock.NewAlloca(vtype.Type())
	decl.Align = ir.Align(vtype.Alignment())

	if !vtype.Equals(varval.TType()) {

		if cast := compiler.CastStore.RunCast(true, vtype, varval, vd.Value.Group, compiler, function, class); cast != nil {
			varval = cast
		} else {
			//compiler error
			//variable value type doesn't match inputted type
			compiler.AddError(errhandle.NewCompileErrorFTok(
				"mismatched types",
				fmt.Sprintf("expected type %s", varval.TType().TypeData().String()),
				vd.typtok,
			))
			return nil
		}
	}

	if llv := varval.LLVal(function); llv != nil {
		function.ActiveBlock.NewStore(llv, decl)
	}

	dv := data.NewVariable(decl, vtype)

	compiler.AddVar(vd.Name, dv)

	return dv
}

func (vd *VarDecl) DeclareGlobal(name string, compiler *Compiler, class *data.Class, static bool, access int) error {

	vtype := vd.getDeclType(compiler, class, compiler.InitFunc)

	if vtype == nil {
		//error
		//tusk can't assume types of globals
		//(atleast not yet)
		compiler.AddError(errhandle.NewCompileErrorFTok(
			"untyped global",
			"add a type to this global",
			vd.vnametok,
		))
		vd.globalerr = true

		if static {
			class.AppendStatic(vd.Name, data.NewInvalidType(), data.NewInvalidType(), access)
		} else {

			class.SType.Fields = append(class.SType.Fields, types.Void)

			class.AppendInstance(vd.Name, data.NewInvalidType(), access)
		}

		return nil
	}

	vd.decltyp = vtype

	if static {

		//static variable

		var decl = compiler.Module.NewGlobal(name, vtype.Type())
		decl.Init = vtype.Default()

		vd.declaration = decl

		nv := data.NewVariable(vd.declaration, vtype)

		class.AppendStatic(vd.Name, nv, nv.TType(), access)
	} else {

		//instance variable

		class.SType.Fields = append(class.SType.Fields, vtype.Type())
		class.TypSiz += vtype.Alignment()

		class.AppendInstance(vd.Name, vtype, access)
	}
	return nil
}

//used specifically for global variable declarations
func (vd *VarDecl) CompileGlobal(compiler *Compiler, class *data.Class, function *data.Function) {

	if vd.globalerr {
		return
	}

	//if the value of the global is nil, we don't need to assign any value that it's defaulted to
	if vd.Value == nil {
		return
	}

	val := vd.Value.Group.Compile(compiler, class, vd.Value, compiler.InitFunc)

	if !vd.decltyp.Equals(val.TType()) {

		if cast := compiler.CastStore.RunCast(true, vd.decltyp, val, vd.Value.Group, compiler, function, class); cast != nil {
			val = cast
		} else {
			//compiler error
			//variable value type doesn't match inputted type
			compiler.AddError(errhandle.NewCompileErrorFTok(
				"mismatched types",
				fmt.Sprintf("expected type %s", val.TType().TypeData().String()),
				vd.typtok,
			))
			vd.globalerr = true
			return
		}
	}

	//it's an instance variable
	//so we make the GEP to the init malloc here
	if vd.declaration == nil {
		//create a new GEP instruction to initialize the struct
		gep := class.Construct.ActiveBlock.NewGetElementPtr(class.SType, class.ConstructAlloc, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, int64(len(class.SType.Fields)-1)))
		vd.declaration = gep
	}

	function.ActiveBlock.NewStore(val.LLVal(function), vd.declaration)
}
