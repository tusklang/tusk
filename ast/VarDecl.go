package ast

import (
	"errors"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
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

func (vd *VarDecl) Compile(compiler *Compiler, class *types.StructType, node *ASTNode, block *ir.Block) value.Value {
	return nil
}

//get the default value for x type
func getDefault(typ types.Type) constant.Constant {
	switch typ.(type) {
	case *types.IntType: //integers default to 0
		return constant.NewInt(types.I32, 0)
	default: //everything else defaults to null
		return &constant.Null{}
	}
}

//used specifically for global variable declarations
func (vd *VarDecl) CompileGlobal(compiler *Compiler, class *types.StructType, static bool) error {

	val := vd.Value.Group.Compile(compiler, class, vd.Value, compiler.InitBlock)

	vtype := val.Type()

	if static {
		name := class.Name() + "_" + vd.Name
		gbl := compiler.Module.NewGlobal(name, vtype)

		gbl.Init = getDefault(vtype)

		compiler.InitBlock.NewStore((val), gbl)

		compiler.StaticGlobals[name] = gbl
		return nil
	}

	class.Fields = append(class.Fields, vtype)

	return nil
}
