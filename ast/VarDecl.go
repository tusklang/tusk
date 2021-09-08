package ast

import (
	"errors"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/tusklang/tusk/data"
	"github.com/tusklang/tusk/tokenizer"
)

type VarDecl struct {
	Name  string
	Type  *ASTNode
	Value *ASTNode
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
		v, e := groupsToAST(grouper(braceMatcher(lex, i, []string{"{", "("}, []string{"}", ")"}, false, "terminator")))
		vd.Value = v[0]
		if e != nil {
			return e
		}
	}

	*i-- //the outer loop will incremenet for us

	return nil
}

func (vd *VarDecl) Compile(compiler *Compiler, class *data.Class, node *ASTNode, block *ir.Block) data.Value {

	varval := vd.Value.Group.Compile(compiler, class, vd.Value, block)
	pvtype := vd.Type.Group.Compile(compiler, class, vd.Type, block)

	var vtype types.Type

	//i'll try and make this less bad later
	switch pvtype.(type) {
	case *data.Class:
		vtype = pvtype.Type()
	case *data.Type:
		vtype = pvtype.Type()
	default:
		//error
	}

	decl := block.NewAlloca(vtype)

	if varval.LLVal(block) != nil {
		block.NewStore(varval.LLVal(block), decl)
	}

	dv := data.NewVariable(data.NewInstruction(decl), data.NewType(vtype), false)

	compiler.AddVar(vd.Name, dv)

	return dv
}

//used specifically for global variable declarations
func (vd *VarDecl) CompileGlobal(compiler *Compiler, class *data.Class, static bool) error {

	val := vd.Value.Group.Compile(compiler, class, vd.Value, compiler.InitBlock)
	pvtype := vd.Type.Group.Compile(compiler, class, vd.Type, compiler.InitBlock)

	var vtype *data.Type

	switch vt := pvtype.(type) {
	case *data.Class:
		vtype = data.NewType(vt.Type())
	case *data.Type:
		vtype = vt
	default:
		//error
	}

	if static {
		class.Static[vd.Name] = data.NewVariable(val, vtype, false)
	} else {
		class.Instance[vd.Name] = data.NewVariable(val, vtype, false)
	}

	return nil
}
