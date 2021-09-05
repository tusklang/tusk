package ast

import (
	"errors"
	"strconv"

	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/tusklang/tusk/tokenizer"
)

type Function struct {
	Params  []*VarDecl //parameter list
	RetType *ASTNode   //return type
	Body    *Block     //function body
}

func (f *Function) Parse(lex []tokenizer.Token, i *int) (e error) {

	if lex[*i].Type != "fn" {
		return errors.New("was not given a function")
	}

	*i++

	//read the return type
	//fn int() {}
	//will also work, because if no braces are present, the next token is returned, and the brace matcher exits
	//if the next value is a variable name, then we know it's a void return type
	//so we will skip the return type

	if lex[*i].Type != "(" {
		rt, e := groupsToAST(groupSpecific(lex, 1, i))
		f.RetType = rt[0]
		if e != nil {
			return e
		}
	}

	if lex[*i].Type != "(" { //it has to be a parenthesis for the paramlist
		return errors.New("functions require a parameter list")
	}

	p, e := groupsToAST(grouper(braceMatcher(lex, i, []string{"("}, []string{")"}, false, "")))
	sub := p[0].Group.(*Block).Sub
	plist := make([]*VarDecl, len(sub))

	for k, v := range sub {

		switch g := v.Group.(type) {
		case *Operation:
			if g.OpType != ":" {
				return errors.New("invalid syntax: named parameters must have a type")
			}

			plist[k] = &VarDecl{
				Name: v.Left[0].Group.(*VarRef).Name,
				Type: v.Right[0],
			}
		default:
			plist[k] = &VarDecl{
				Type: v,
			}

		}
	}

	f.Params = plist

	if e != nil {
		return e
	}

	*i++

	if lex[*i].Type != "{" {
		*i-- //move back because there is no brace
		return nil
	}

	f.Body = grouper(braceMatcher(lex, i, []string{"{"}, []string{"}"}, false, "terminator"))[0].(*Block)

	if e != nil {
		return e
	}

	return nil
}

func (f *Function) Compile(compiler *Compiler, class *types.StructType, node *ASTNode) constant.Constant {

	var rt types.Type = types.Void //defaults to void

	if f.RetType != nil {
		var e error
		rt, e = compiler.FetchType(class, f.RetType.Group)
		_ = e
	}

	rf := compiler.Module.NewFunc("f"+strconv.Itoa(compiler.TmpVar()), rt)

	if f.Body != nil {
		fblock := rf.NewBlock("")
		fblock.NewRet(nil)
		_ = fblock
	}

	return rf
}
