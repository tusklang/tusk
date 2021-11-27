package ast

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	"github.com/tusklang/tusk/data"
	"github.com/tusklang/tusk/errhandle"
	"github.com/tusklang/tusk/tokenizer"
)

type Array struct {
	Siz *ASTNode
	Typ *ASTNode
	Arr []*ASTNode

	//if the array is used as an index (var a = [2]i32{1, 2}; a[0];)
	//the first statement is being used an expression
	//the second statement is being used an index
	useAsIndex bool

	//if the array has an initializer list ({1, 2, 3})
	hasInit bool

	//tokens used
	btok, typtok, arrtok tokenizer.Token
	/////////////

	//used during compiling
	csiz data.Value
	ctyp data.Type
}

func (a *Array) Parse(lex []tokenizer.Token, i *int, stopAt []string) *errhandle.TuskError {

	a.btok = lex[*i]

	sizl := braceMatcher(lex, i, []string{"[", "{", "("}, []string{"]", "}", ")"}, true, "")
	*i++
	sizg, e := grouper(sizl)

	if e != nil {
		return e
	}

	siz, e := groupsToAST(sizg)

	if e != nil {
		return e
	}

	if len(siz) == 1 {
		a.Siz = siz[0]
	} else if len(siz) != 0 {
		//error
		//size can't be multiple statements
	}

	if *i >= len(lex) {
		return nil
	}

	//parse the typename
	if lex[*i].Type == "(" || lex[*i].Type == "varname" {
		a.typtok = lex[*i]
		typg, e := groupSpecific(lex, i, nil, 1)

		if e != nil {
			return e
		}

		typ, e := groupsToAST(typg)

		if e != nil {
			return e
		}

		a.Typ = typ[0]
	}

	if *i >= len(lex) {
		return nil
	}

	//if the stopAt includes an opening brace, we can't read the array contents
	//(they are likely meant to be something else, like a function body)
	for _, v := range stopAt {
		if v == "{" {
			*i--
			return nil
		}
	}

	//arrays are written like:
	//var v: []Type = []Type{}
	//if there is no { after the type, then it is being used as the var type
	if lex[*i].Type == "{" {
		a.typtok = lex[*i]
		arrl := braceMatcher(lex, i, []string{"{"}, []string{"}"}, true, "")
		arrg, e := grouper(arrl)

		if e != nil {
			return e
		}

		arr, e := groupsToAST(arrg)

		if e != nil {
			return e
		}

		a.Arr = arr
		a.hasInit = true
	} else {
		//there was no array content
		*i--
	}

	return nil
}

func (a *Array) GetMTok() tokenizer.Token {
	return a.btok
}

func (a *Array) CompileSlice(compiler *Compiler, class *data.Class, node *ASTNode, function *data.Function) data.Value {

	var decl value.Value

	if a.hasInit {
		block := function.ActiveBlock
		tdecl := block.NewAlloca(types.NewPointer(a.ctyp.Type()))
		tdecl.Align = 8
		decl = tdecl

		mallocf := compiler.LinkedFunctions["malloc"]
		mallocc := block.NewCall(mallocf, block.NewMul(
			constant.NewInt(types.I32, int64(len(a.Arr))),
			constant.NewInt(types.I32, int64(a.ctyp.Alignment())),
		))

		block.NewStore(block.NewBitCast(
			mallocc,
			types.NewPointer(a.ctyp.Type()),
		), decl)

		loadeddecl := block.NewLoad(types.NewPointer(a.ctyp.Type()), decl)
		for k, v := range a.Arr {
			vc := v.Group.Compile(compiler, class, v, function)
			gep := block.NewGetElementPtr(a.ctyp.Type(), loadeddecl, constant.NewInt(types.I32, int64(k)))
			block.NewStore(vc.LLVal(function), gep)
		}
	}

	return data.NewSliceArray(a.ctyp, decl, nil)
}

func (a *Array) CompileFixedArray(compiler *Compiler, class *data.Class, node *ASTNode, function *data.Function) data.Value {

	sizi := a.csiz.(*data.Integer).GetInt()
	var decl value.Value
	var curlen value.Value

	if a.hasInit {

		block := function.ActiveBlock
		arrtyp := types.NewArray(uint64(sizi), a.ctyp.Type())
		tdecl := block.NewAlloca(arrtyp)
		tdecl.Align = ir.Align(16)

		decl = tdecl

		//fill the array with the values needed
		for k, v := range a.Arr {
			vc := v.Group.Compile(compiler, class, v, function)
			gep := block.NewGetElementPtr(arrtyp, decl, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, int64(k)))
			block.NewStore(vc.LLVal(function), gep)
		}

		curlen := block.NewAlloca(types.I32)
		block.NewStore(constant.NewInt(types.I32, sizi), curlen)
	}

	return data.NewFixedArray(a.ctyp, decl, curlen, uint64(sizi))
}

func (a *Array) CompileVariedLengthArray(compiler *Compiler, class *data.Class, node *ASTNode, function *data.Function) data.Value {

	if function == nil || function.ActiveBlock == nil || function == compiler.InitFunc {
		//error
		//cannot use varied length arrays outside of a function
		compiler.AddError(errhandle.NewCompileErrorFTok(
			"invalid varied array",
			"cannot use a varied array outside of a function body",
			a.btok,
		))
		return nil
	}

	sizi := compiler.CastStore.RunCast(true, data.NewPrimitive(types.I64), a.csiz, compiler, function, class) //force the length to an i64
	var curlen value.Value
	var alc *ir.InstAlloca
	block := function.ActiveBlock
	sizill := sizi.LLVal(function)

	if a.hasInit {

		alc = block.NewAlloca(a.ctyp.Type())
		alc.NElems = sizill
		alc.Align = ir.Align(16)

		tcurlen := block.NewAlloca(sizill.Type())
		tcurlen.Align = ir.Align(4)
		block.NewStore(sizill, tcurlen)
		curlen = tcurlen

		for k, v := range a.Arr {
			vc := v.Group.Compile(compiler, class, v, function)
			gep := block.NewGetElementPtr(a.ctyp.Type(), alc, constant.NewInt(types.I32, int64(k)))
			block.NewStore(vc.LLVal(function), gep)
		}
	}

	return data.NewVariedLengthArray(a.ctyp, alc, curlen, sizill)
}

func (a *Array) Compile(compiler *Compiler, class *data.Class, node *ASTNode, function *data.Function) data.Value {

	if a.Siz == nil {
		//it's a slice array
		a.ctyp = a.Typ.Group.Compile(compiler, class, a.Typ, function).(data.Type)
		return a.CompileSlice(compiler, class, node, function)
	} else {
		a.csiz = a.Siz.Group.Compile(compiler, class, a.Siz, function)

		if a.useAsIndex {
			return a.csiz
		}

		a.ctyp = a.Typ.Group.Compile(compiler, class, a.Typ, function).(data.Type)

		switch a.csiz.(type) {
		case *data.Integer:
			//it's a fixed array
			return a.CompileFixedArray(compiler, class, node, function)
		default:
			//it's a varied length array
			return a.CompileVariedLengthArray(compiler, class, node, function)
		}
	}
}
