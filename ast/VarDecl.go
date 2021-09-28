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

	declaration value.Value
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
		t, e := groupsToAST(groupSpecific(lex, i, []string{"=", ";"}))
		if e != nil {
			return e
		}
		vd.Type = t[0]
	}

	//has a value assigned to it
	if lex[*i].Name == "=" {
		*i++
		v, e := groupsToAST(grouper(braceMatcher(lex, i, allopeners, allclosers, false, "terminator")))
		vd.Value = v[0]
		if e != nil {
			return e
		}
	}

	*i-- //the outer loop will incremenet for us

	return nil
}

func (vd *VarDecl) getDeclType(compiler *Compiler, class *data.Class, block *ir.Block) data.Type {

	pvtype := vd.Type.Group.Compile(compiler, class, vd.Type, block)

	var vtype data.Type

	switch vt := pvtype.(type) {
	case *data.Class:
		vtype = data.NewInstance(vt)
	case *data.Primitive:
		vtype = vt
	case *data.Function:
		vtype = vt
	default:
		//error
	}

	return vtype
}

func (vd *VarDecl) Compile(compiler *Compiler, class *data.Class, node *ASTNode, block *ir.Block) data.Value {

	varval := vd.Value.Group.Compile(compiler, class, vd.Value, block)

	vtype := vd.getDeclType(compiler, class, block)

	decl := block.NewAlloca(vtype.Type())

	if llv := varval.LLVal(block); llv != nil {
		block.NewStore(llv, decl)
	}

	dv := data.NewVariable(decl, vtype)

	compiler.AddVar(vd.Name, dv)

	return dv
}

func (vd *VarDecl) DeclareGlobal(name string, compiler *Compiler, class *data.Class, static bool) error {

	vtype := vd.getDeclType(compiler, class, compiler.InitBlock)

	if static {

		//static variable

		decl := compiler.Module.NewGlobal(name, vtype.Type())
		decl.Init = vtype.Default()

		vd.declaration = decl

		nv := data.NewVariable(vd.declaration, vtype)

		class.Static[vd.Name] = nv
	} else {

		//instance variable

		class.SType.Fields = append(class.SType.Fields, vtype.Type())

		//create a new GEP to the initialize struct
		gep := class.Construct.NewGetElementPtr(class.SType, class.ConstructAlloc, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, int64(len(class.SType.Fields)-1)))
		vd.declaration = gep

		class.AppendInstance(vd.Name, vtype)
	}

	return nil
}

//used specifically for global variable declarations
func (vd *VarDecl) CompileGlobal(compiler *Compiler, class *data.Class, block *ir.Block) {
	val := vd.Value.Group.Compile(compiler, class, vd.Value, compiler.InitBlock)
	block.NewStore(val.LLVal(block), vd.declaration)
}
