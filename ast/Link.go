package ast

import (
	"github.com/tusklang/tusk/data"
	"github.com/tusklang/tusk/tokenizer"
)

type Link struct {
	stname /*stored tname after varname mangling*/, TName, AName string
	DType                                                        *ASTNode
	Access                                                       int
}

func (l *Link) Parse(lex []tokenizer.Token, i *int) error {

	//format looks like
	//	link var tusk_name: fn() -> asm_name

	*i++

	if lex[*i].Name != "var" {
		//error
	}

	*i++

	if lex[*i].Type != "variable" {
		//must be the varname
	}

	tname := lex[*i].Name
	*i++

	if lex[*i].Name != ":" {
		//error
		//must supply type
	}

	*i++

	dtype, e := groupsToAST(groupSpecific(lex, i, nil, 1))

	if e != nil {
		return e
	}

	if lex[*i].Name != "->" {
		//error
	}

	*i++

	aname := lex[*i].Name

	l.TName = tname
	l.stname = tname
	l.AName = aname
	l.DType = dtype[0]
	l.Access = 2 //access is private by default

	return nil
}

func (l *Link) Compile(compiler *Compiler, class *data.Class, node *ASTNode, function *data.Function) data.Value {

	aname := l.AName //name in the linked binary

	dtype := l.DType.Group.Compile(compiler, class, l.DType, function)

	if dtype.TypeData().Name() != "func" {
		//error
		//linked values must be functions
	}

	dfunc := dtype.(*data.Function).LLFunc

	dfunc.SetName(aname)

	tfd := data.NewFunc(dfunc, dtype.(*data.Function).RetType())
	tfd.SetLName(l.stname)
	compiler.AddVar(l.TName, tfd)

	class.AppendStatic(l.stname, tfd, tfd.TType(), l.Access)

	return nil
}
