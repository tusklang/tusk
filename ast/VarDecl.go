package ast

import (
	"errors"

	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
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
		t, e := groupsToAST(groupSpecific(lex, 1, i))
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

func (vd *VarDecl) Compile(compiler *Compiler, class *types.StructType, node *ASTNode) constant.Constant {
	return nil
}

//used specifically for global variable declarations
func (vd *VarDecl) CompileGlobal(compiler *Compiler, class *types.StructType, static bool) error {

	vtype, e := compiler.FetchType(class, vd.Type.Group)

	if e != nil {
		return e
	}

	if static {
		val := vd.Value.Group.Compile(compiler, class, vd.Value)

		name := class.Name() + "_" + vd.Name
		gbl := compiler.Module.NewGlobalDef(name, val)

		compiler.StaticGlobals[name] = gbl
		return nil
	}

	class.Fields = append(class.Fields, vtype)

	return nil
}
