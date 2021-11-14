package ast

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/tusklang/tusk/data"
	"github.com/tusklang/tusk/tokenizer"
)

type Link struct {
	stname /*<- stored tname after varname mangling*/, TName, AName string
	DType                                                           *ASTNode
	Access                                                          int
}

func (l *Link) Parse(lex []tokenizer.Token, i *int) error {

	//format looks like
	//	link fn tusk_name() -> asm_name

	*i++

	if lex[*i].Name != "fn" {
		//error
	}

	fnd := groupSpecific(lex, i, nil, 1)
	dtype, e := groupsToAST(fnd)

	if e != nil {
		return e
	}

	if lex[*i].Name != "->" {
		//error
	}

	*i++

	aname := lex[*i].Name

	l.TName = fnd[0].(*Function).Name
	l.stname = l.TName
	l.AName = aname
	l.DType = dtype[0]
	l.Access = 2 //access is private by default

	fnd[0].(*Function).Name = "" //remove the name, explained in Function.go

	return nil
}

func (l *Link) addToClass(lf *ir.Func, compiler *Compiler, dtype data.Value, class *data.Class) data.Value {
	tfd := data.NewLinkedFunc(lf, dtype.(*data.Function).RetType())
	tfd.SetLName(l.stname)
	compiler.AddVar(l.TName, tfd)

	class.AppendStatic(l.stname, tfd, tfd.TType(), l.Access)
	return nil
}

func (l *Link) Compile(compiler *Compiler, class *data.Class, node *ASTNode, function *data.Function) data.Value {

	aname := l.AName //name in the linked binary
	dtype := l.DType.Group.Compile(compiler, class, l.DType, function)

	if dtype.TypeData().Name() != "func" {
		//error
		//linked values must be functions
	}

	if lf, exists := compiler.LinkedFunctions[aname]; exists {
		return l.addToClass(lf, compiler, dtype, class)
	}

	dfunc := dtype.(*data.Function).LLFunc

	dfunc.SetName(aname)
	dfunc.Params = nil
	dfunc.Sig.Variadic = true        //make it a variadic function in case it is declared elsewhere
	dfunc.Sig.RetType = types.I64Ptr //make the return an i64 pointer, this way we can return any value, and cast it appropriately when called. the appropriate cast is provided by dtype.(*data.Function).RetType()

	compiler.LinkedFunctions[l.AName] = dfunc

	//add the linked function to the llvm ir
	//the function compiler only adds it to the ir if it has a body
	//and linked functions do not
	//so we add it manually here
	compiler.Module.Funcs = append(compiler.Module.Funcs, dfunc)

	return l.addToClass(dfunc, compiler, dtype, class)
}
