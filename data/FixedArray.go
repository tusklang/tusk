package data

import (
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

type FixedArray struct {
	atype  Type
	decl   value.Value
	curlen value.Value
	length uint64
}

func NewFixedArray(atype Type, decl, curlen value.Value, length uint64) *FixedArray {
	return &FixedArray{
		atype:  atype,
		decl:   decl,
		curlen: curlen,
		length: length,
	}
}

func (a *FixedArray) ValType() Type {
	return a.atype
}

func (a *FixedArray) Length() uint64 {
	return a.length
}

func (a *FixedArray) LLVal(function *Function) value.Value {

	//we load decl and return it
	//if you are wondering why, a variable uses a value's LLVal() to store
	//but if we just give decl, the value and type will be a pointer
	//e.g. we want to store [3]i32{1, 2, 3} into variable 'a'
	//tusk code would look like:
	/*
		var a = [3]i32{1, 2, 3};
	*/
	//the `[3]i32{1, 2, 3}` is the fixed array, seen here
	//that is stored within the `decl` part of this struct
	//so it's really more like this under the hood:
	/*
		var decl = &[3]i32{1, 2, 3};
		var a = &decl;
	*/
	//we don't want this, so we load decl and give it to the `a` variable
	//and it becomes something like this:
	/*
		var decl = &[3]i32{1, 2, 3};
		var loaded = *decl;
		var a = &loaded;
	*/

	//why not just make `a` an **[3]i32 pointer?
	//well it comes down to globals
	//take a look at the following tusk code
	/*
		var a: [3]i32 = [3]i32{1, 2, 3};

		pub stat fn main() {
			a[0];
		};
	*/
	//simple enough, right? just a global fixed array and an indexing operations on it in `main`
	//this breaks down to something like this under the hood, in pseudocode
	/*
		declare a global called a of type *[3]i32

		create a function called init
			create a variable called decl of type [3]i32
			store the value 1 in the first index of decl
			store the value 2 in the second index of decl
			store the value 3 in the third index of decl
			store the variable decl in the global a

		create a function called main
			load the value of the global a and store it in a variable called loaded
			find the first index of the variable load
	*/
	//this works yeah? well, it doesn't thanks to the stack
	//the stack pops `decl` off at the end of the `init` function
	//but decl's pointer is still inside `a`
	//so if we dereference `a` for it's stored value, we get a pointer to `decl` which isn't valid after `init` ends...

	//funnily enough, this code does work, but is **extremely** unsafe
	//when we pop off the stack, we're not removing any data, but instead we're moving the stack pointer
	//meaning the decl pointer still has the data of decl
	//but if we allocate more data, we can override `decl`'s pointer's value
	//i was kinda confused on that, so thanks @Whimpers#3099 for clearing that up to me :)

	return function.ActiveBlock.NewLoad(a.Type(), a.decl)
}

func (a *FixedArray) TType() Type {
	return a
}

func (a *FixedArray) Type() types.Type {
	return types.NewArray(a.Length(), a.ValType().Type())
}

func (a *FixedArray) TypeData() *TypeData {
	td := NewTypeData("fixed")
	td.AddFlag("array")
	td.AddFlag("type")
	td.AddOtherDat("valtyp", a.ValType().(Value))
	td.AddOtherDat("length", NewInteger(constant.NewInt(types.I32, int64(a.Length()))))
	return td
}

func (a *FixedArray) InstanceV() value.Value {
	return nil
}

func (a *FixedArray) Default() constant.Constant {
	return constant.NewUndef(a.Type())
}

func (a *FixedArray) Equals(t Type) bool {
	switch c := t.(type) {
	case *FixedArray:
		return a.ValType().Equals(c.ValType()) && a.Length() == c.Length()
	}
	return false
}

func (a *FixedArray) Alignment() uint64 {
	return 16
}
