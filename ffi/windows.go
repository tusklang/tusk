//+build windows
package ommffi

//#include "wincall.h"
import "C"
import (
	"errors"
	"unsafe"
)

type OMM_C_MODULE **C.struct_HINSTANCE__

func LoadLib(filename string) (OMM_C_MODULE, error) {
	lib := C.loadlib(C.CString(filename))

	if lib.error == C.bool(true) {
		var tmp OMM_C_MODULE
		return tmp, errors.New("Could not load library: " + filename)
	}

	return lib.module, nil
}

func CallProc(lib OMM_C_MODULE, procname string, argv []unsafe.Pointer, argc int) unsafe.Pointer {

	var called *unsafe.Pointer

	if len(argv) == 0 {
		var tmp *unsafe.Pointer
		called = C.callproc((**C.struct_HINSTANCE__)(lib), C.CString(procname), tmp, 0)
	} else {
		called = C.callproc((**C.struct_HINSTANCE__)(lib), C.CString(procname), &argv[0], C.int(argc))
	}

	return *called
}
