package main

import (
	"fmt"
	"unsafe"

	"github.com/omm-lang/goat"
	"github.com/omm-lang/omm/lang/interpreter"
	"github.com/omm-lang/omm/lang/types"
	oatenc "github.com/omm-lang/omm/oat/encoding"
)
import "C"

//export CallFunc
func CallFunc(oatf *C.char, fname *C.char, args []unsafe.Pointer) (unsafe.Pointer, *C.char) {

	//decompile the oat
	data, e := oatenc.OatDecode(C.GoString(oatf), 0)

	if e != nil {
		return nil, C.CString(e.Error())
	}

	//create the instance and get the variable
	var instance = goat.NewInstance(data)
	var fn = instance.Fetch(C.GoString(fname))

	//if it does not exist, error
	if fn == nil {
		return nil, C.CString(fmt.Sprintf("Variable %s does not exist", fname))
	}

	//if it is not a function, error
	if (*fn.Value).Type() != "function" {
		return nil, C.CString(fmt.Sprintf("Variable %s is not a function", fname))
	}

	//create the argv as an omm array
	var argv types.OmmArray

	for _, v := range args {
		argv.PushBack(*(*types.OmmType)(v))
	}

	//call the func
	return unsafe.Pointer(interpreter.Operations["function <- array"](*fn.Value, argv, instance, []string{"at goat caller"}, 0, "goat caller", 0)), C.CString("")
}

func main() {}
