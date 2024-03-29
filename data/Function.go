package data

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

type Function struct {
	LLFunc      *ir.Func
	ActiveBlock *ir.Block
	ParamTypes  []Type
	ret         Type

	IsPure bool

	nam string

	//list of terminators to append
	//for example:
	/*
		if (true) {
			if (false) {
				//some task in the second if
			}
		}

		//some code here

		would become:

		mainbody:
			br i32 1, label %if1-true, label %if1-false     ; if (true) {}

		if1-true:
			br i32 0, label %if2-true, label %if2-false     ; if (false) {}

		if1-false:
			br label %rest-of-main-body						; go straight to rest-of-main-body since there is no `else`

		if2-true:
			; some task in the second if
			br label %if1-after

		if2-false:
			br label %if1-after

		if1-after:
			br label %rest-of-main-body						; go to the rest of the main body (outside of the first if)

		rest-of-main-body:
			; some code here

		in this example, how would if1-after know to go to rest-of-main-body?
		this stack stores that
	*/
	todoTerms []ir.Terminator

	//if the function is linked to a binary, instead of being declared in tusk
	linkedfunc bool

	//stuff for methods

	//if the given function is a method
	IsMethod bool

	Instance value.Value

	//given functions's class, if it is a method
	MethodClass *Class
	///////////////////

	AllocatedObjs []value.Value //objects allocated (on the heap) in this function
}

func NewFunc(f *ir.Func, ret Type) *Function {
	return &Function{
		LLFunc: f,
		ret:    ret,
	}
}

func NewLinkedFunc(f *ir.Func, ret Type) *Function {
	rf := NewFunc(f, ret)
	rf.linkedfunc = true
	return rf
}

func CloneFunc(f *Function) *Function {
	var ret = new(Function)
	ret.LLFunc = f.LLFunc
	ret.ActiveBlock = f.ActiveBlock
	ret.ParamTypes = f.ParamTypes
	ret.ret = f.ret
	ret.nam = f.nam
	ret.linkedfunc = f.linkedfunc
	ret.IsMethod = f.IsMethod
	ret.MethodClass = f.MethodClass
	return ret
}

func (f *Function) LLVal(function *Function) value.Value {
	return f.LLFunc
}

func (f *Function) RetType() Type {
	return f.ret
}

func (f *Function) Default() constant.Constant {
	return constant.NewNull(f.LLFunc.Typ)
}

func (f *Function) TType() Type {
	return f
}

func (f *Function) Type() types.Type {
	return f.LLFunc.Type()
}

func (f *Function) TypeData() *TypeData {
	td := NewTypeData("func")
	if f.linkedfunc {
		td.AddFlag("linked")
	}
	if f.IsMethod {
		td.AddFlag("method")
	}
	return td
}

func (f *Function) InstanceV() value.Value {
	return f.Instance
}

func (f *Function) Equals(t Type) bool {
	return f.LLFunc.Type().Equal(t.Type())
}

func (f *Function) PopTermStack() ir.Terminator {
	r := f.LastTermStack()

	if r == nil {
		return r
	}

	f.todoTerms = f.todoTerms[:len(f.todoTerms)-1]
	return r
}

func (f *Function) PushTermStack(v ir.Terminator) {
	f.todoTerms = append(f.todoTerms, v)
}

func (f *Function) LastTermStack() ir.Terminator {

	if len(f.todoTerms) == 0 {
		return nil
	}

	return f.todoTerms[len(f.todoTerms)-1]
}

func (f *Function) SetLName(n string) {
	f.nam = n
}

func (f *Function) GetLName() string {
	return f.nam
}

func (f *Function) Alignment() uint64 {
	return 8
}
