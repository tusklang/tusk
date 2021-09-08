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

func (vd *VarDecl) Compile(compiler *Compiler, class *types.StructType, node *ASTNode, block *ir.Block) data.Value {

	varval := vd.Value.Group.Compile(compiler, class, vd.Value, block)
	pvtype := vd.Type.Group.Compile(compiler, class, vd.Type, block)

	var vtype types.Type

	//i'll try and make this less bad later
	switch pvtype.(type) {
	case *data.Type:
		vtype = pvtype.Type()
	default:
		//also error
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
func (vd *VarDecl) CompileGlobal(compiler *Compiler, class *types.StructType, static bool) error {

	val := vd.Value.Group.Compile(compiler, class, vd.Value, compiler.InitBlock)
	pvtype := vd.Type.Group.Compile(compiler, class, vd.Type, compiler.InitBlock)

	var vtype types.Type

	switch pvtype.(type) {
	case *data.Type:
		vtype = pvtype.Type()
	default:
		//also error
	}

	if static {
		name := class.Name() + "_" + vd.Name
		gbl := compiler.Module.NewGlobal(name, vtype)

		gbl.Init = data.GetDefault(vtype)

		if val.LLVal(compiler.InitBlock) != nil {
			compiler.InitBlock.NewStore(val.LLVal(compiler.InitBlock), gbl)
		}

		compiler.StaticGlobals[name] = gbl
		compiler.AddVar(name, data.NewVariable(data.NewInstruction(gbl), data.NewType(vtype), false))
	} else {
		class.Fields = append(class.Fields, vtype) //append the field to the class if it's not a static field
	}

	return nil
}
