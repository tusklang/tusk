package data

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

type Lambda struct {
	fn   *Function
	ltyp types.Type
}

func NewLambda(fn *Function, ltyp types.Type) *Lambda {
	return &Lambda{
		fn:   fn,
		ltyp: ltyp,
	}
}

func (l *Lambda) Func() *Function {
	return l.fn
}

func (l *Lambda) LLVal(function *Function) value.Value {
	return l.fn.LLFunc
}

func (l *Lambda) RetType() Type {
	return l.fn.ret
}

func (l *Lambda) Default() constant.Constant {
	return constant.NewNull(l.fn.LLFunc.Typ)
}

func (l *Lambda) TType() Type {
	return l
}

func (l *Lambda) Type() types.Type {
	return l.fn.LLFunc.Type()
}

func (l *Lambda) TypeData() *TypeData {
	td := NewTypeData("lambda")
	return td
}

func (l *Lambda) InstanceV() value.Value {
	return l.fn.Instance
}

func (l *Lambda) Equals(t Type) bool {
	return l.fn.LLFunc.Type().Equal(t.Type())
}

func (l *Lambda) PopTermStack() ir.Terminator {
	r := l.fn.LastTermStack()

	if r == nil {
		return r
	}

	l.fn.todoTerms = l.fn.todoTerms[:len(l.fn.todoTerms)-1]
	return r
}

func (l *Lambda) PushTermStack(v ir.Terminator) {
	l.fn.todoTerms = append(l.fn.todoTerms, v)
}

func (l *Lambda) LastTermStack() ir.Terminator {

	if len(l.fn.todoTerms) == 0 {
		return nil
	}

	return l.fn.todoTerms[len(l.fn.todoTerms)-1]
}

func (l *Lambda) Alignment() uint64 {
	return 8
}
