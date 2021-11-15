package ast

import (
	"errors"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	"github.com/tusklang/tusk/data"
	"github.com/tusklang/tusk/tokenizer"
)

type VarDecl struct {
	Name  string
	Type  *ASTNode
	Value *ASTNode

	//for globals
	declaration value.Value
	decltyp     data.Type
}

func (vd *VarDecl) Parse(lex []tokenizer.Token, i *int) error {

	*i++

	if lex[*i].Type != "varname" {
		return errors.New("expected a variable name")
	}

	vd.Name = lex[*i].Name

	*i++

	//has a specified type
	if lex[*i].Name == ":" {
		*i++
		t, e := groupsToAST(groupSpecific(lex, i, []string{"=", ";"}, -1))
		if e != nil {
			return e
		}
		vd.Type = t[0]
	}

	//has a value assigned to it
	if lex[*i].Name == "=" {
		*i++
		v, e := groupsToAST(grouper(braceMatcher(lex, i, allopeners, allclosers, false, "terminator")))
		if e != nil {
			return e
		}
		vd.Value = v[0]
	}

	*i-- //the outer loop will incremenet for us

	return nil
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
	default:
		//error
	}

	return vtype
}

func (vd *VarDecl) Compile(compiler *Compiler, class *data.Class, node *ASTNode, function *data.Function) data.Value {

	var vtype data.Type
	var varval data.Value

	if vd.Value == nil && vd.Type != nil {
		//var varname: typename
		vtype = vd.getDeclType(compiler, class, function)
		varval = data.NewVariable(vtype.Default(), vtype)
	} else if vd.Value != nil && vd.Type == nil {
		//var varname = value
		varval = vd.Value.Group.Compile(compiler, class, vd.Value, function)
		vtype = varval.TType()
	} else if vd.Value != nil && vd.Type != nil {
		//var varname: typename = value;
		vtype = vd.getDeclType(compiler, class, function)
		varval = vd.Value.Group.Compile(compiler, class, vd.Value, function)
	} else {
		//var varname
		//error
		//cannot determine what type this variable is
	}

	decl := function.ActiveBlock.NewAlloca(vtype.Type())
	decl.Align = ir.Align(vtype.TypeSize())

	if !vtype.Equals(varval.TType()) {
		//compiler error
		//variable value type doesn't match inputted type
	}

	if llv := varval.LLVal(function.ActiveBlock); llv != nil {
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
		class.TypSiz += vtype.TypeSize()

		class.AppendInstance(vd.Name, vtype, access)
	}
	return nil
}

//used specifically for global variable declarations
func (vd *VarDecl) CompileGlobal(compiler *Compiler, class *data.Class, function *data.Function) {

	//if the value of the global is nil, we don't need to assign any value that it's defaulted to
	if vd.Value == nil {
		return
	}

	val := vd.Value.Group.Compile(compiler, class, vd.Value, compiler.InitFunc)

	if !vd.decltyp.Equals(val.TType()) {
		//compiler error
		//variable value type doesn't match inputted type
	}

	//it's an instance variable
	//so we make the GEP to the init malloc here
	if vd.declaration == nil {
		//create a new GEP instruction to initialize the struct
		gep := class.Construct.ActiveBlock.NewGetElementPtr(class.SType, class.ConstructAlloc, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, int64(len(class.SType.Fields)-1)))
		vd.declaration = gep
	}

	function.ActiveBlock.NewStore(val.LLVal(function.ActiveBlock), vd.declaration)
}
