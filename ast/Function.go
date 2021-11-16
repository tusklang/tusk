package ast

import (
	"errors"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/tusklang/tusk/data"

	"github.com/tusklang/tusk/tokenizer"
)

type Function struct {
	Name     string     //function name
	Params   []*VarDecl //parameter list
	Body     *Block     //function body
	RetType  *ASTNode   //return type
	isMethod bool
}

func (f *Function) Parse(lex []tokenizer.Token, i *int) (e error) {
	*i++ //skip the "fn" token

	if lex[*i].Type == "varname" {
		//read the function name if there is one
		f.Name = lex[*i].Name
		*i++
	}

	if lex[*i].Type != "(" {
		//error
		//functions require a parameter list
	}

	p, e := groupsToAST(grouper(braceMatcher(lex, i, []string{"("}, []string{")"}, false, "")))

	if e != nil {
		return e
	}

	sub := p[0].Group.(*Block).Sub
	plist := make([]*VarDecl, len(sub))

	*i++

	for k, v := range sub {

		switch g := v.Group.(type) {
		case *Operation:

			switch g.OpType {
			case ":":
				plist[k] = &VarDecl{
					Name: v.Left[0].Group.(*VarRef).Name,
					Type: v.Right[0],
				}
			case "*":
				plist[k] = &VarDecl{
					Type: v,
				}
			case ".":
				plist[k] = &VarDecl{
					Type: v,
				}
			default:
				return errors.New("invalid syntax: named parameters must have a type")
			}

		default:

			plist[k] = &VarDecl{
				Type: v,
			}

		}
	}

	f.Params = plist

	if lex[*i].Type != "{" && lex[*i].Type != "terminator" && lex[*i].Type != "operation" {
		//read the return type
		//if there is no body or terminator next, it has to be a return
		rtg := groupSpecific(lex, i, []string{"{", "terminator"}, -1)
		rt, e := groupsToAST(rtg)

		if e != nil {
			return e
		}

		f.RetType = rt[0]
	}

	if lex[*i].Type == "{" {
		f.Body = grouper(braceMatcher(lex, i, []string{"{"}, []string{"}"}, false, ""))[0].(*Block)
		return nil
	}

	*i--
	return nil
}

func (f *Function) CompileSig(compiler *Compiler, class *data.Class, node *ASTNode, function *data.Function) *data.Function {
	var rt data.Type = data.NewPrimitive(types.Void) //defaults to void

	if f.RetType != nil {
		rt = f.RetType.Group.Compile(compiler, class, f.RetType, function).TType()
	}

	var params []*ir.Param

	if f.isMethod {
		//make the first argument the `this` or `self` value
		//methods can use the `this` keyword to access members of the current instance
		//this first parameter will store that
		params = append(
			params,
			ir.NewParam("", types.NewPointer(class.SType)),
		)

	}

	for _, v := range f.Params {
		typ := v.Type.Group.Compile(compiler, class, v.Type, function)
		p := ir.NewParam(
			"",
			typ.Type(),
		)
		params = append(params, p)
		compiler.AddVar(v.Name, data.NewInstVariable(p, typ.TType()))
	}

	if f.Name != "" {
		//error
		//function names are only for global functions
		//we remove the name of the function from the object when compiling a global function
		//if it has a name, it is in a local scope, and it was not declared anonymous
	}

	rf := ir.NewFunc("", rt.Type(), params...)

	ffunc := data.NewFunc(rf, rt)
	ffunc.IsMethod = f.isMethod

	if ffunc.IsMethod {
		ffunc.MethodClass = class
	}

	return ffunc
}

func (f *Function) CompileBody(ffunc *data.Function, compiler *Compiler, class *data.Class, node *ASTNode, function *data.Function) {
	rf := ffunc.LLFunc

	if f.Body != nil {
		fblock := rf.NewBlock("")

		if f.RetType == nil { //if the function returns void, append a `return void` to the term stack
			ffunc.PushTermStack(ir.NewRet(nil))
		}

		ffunc.ActiveBlock = fblock
		f.Body.Compile(compiler, class, nil, ffunc)

		//pop the entire term stack
		for v := ffunc.PopTermStack(); v != nil; v = ffunc.PopTermStack() {
			ffunc.ActiveBlock.Term = v
		}

		//add the function to the actual llvm bytecode (only if it has a body)
		compiler.Module.Funcs = append(compiler.Module.Funcs, rf)
	}
}

func (f *Function) Compile(compiler *Compiler, class *data.Class, node *ASTNode, function *data.Function) data.Value {
	ffunc := f.CompileSig(compiler, class, node, function)
	f.CompileBody(ffunc, compiler, class, node, function)
	return ffunc
}

func (f *Function) DeclareGlobal(compiler *Compiler, class *data.Class, static bool, access int) {
	stoname := f.Name

	//we do this to detect functions declared within the non-global scope that have names
	//if it passes through the global scope, though, it will have a name
	//we store the name above, and remove the name when passing to the Compile() method
	//this way the compiler only throws errors on lambda functions with names
	f.Name = ""

	if !static {
		f.isMethod = true
	}

	var fn = f.CompileSig(compiler, class, nil, nil)
	fn.SetLName(stoname)

	if static {
		class.AppendStatic(stoname, fn, fn.TType(), access)
	} else {
		class.NewMethod(stoname, fn, access)
	}

	f.Name = stoname //put the name back
}

func (f *Function) CompileGlobal(compiler *Compiler, class *data.Class, static bool) {

	if f.Body == nil {
		return
	}

	var fn *data.Function

	if static {
		fn = class.Static[f.Name].Value.(*data.Function)
	} else {
		fn = class.Methods[f.Name].Value.(*data.Function)
	}

	f.CompileBody(fn, compiler, class, nil, nil)
}
