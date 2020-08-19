package types

import (
	"unsafe"

	ommffi "github.com/omm-lang/omm/ffi"
)

type OmmLibrary struct {
	Module ommffi.OMM_C_MODULE
}

func (l OmmLibrary) CallProc(procname string, args []*OmmType) *OmmType {
	var argv = make([]unsafe.Pointer, len(args))

	for k, v := range args {
		argv[k] = unsafe.Pointer(v)
	}

	return (*OmmType)(ommffi.CallProc(l.Module, procname, argv, len(args)))
}

func (l OmmLibrary) Format() string {
	return "{ native library }"
}

func (l OmmLibrary) Type() string {
	return "nativelib"
}

func (l OmmLibrary) TypeOf() string {
	return l.Type()
}

func (_ OmmLibrary) Deallocate() {}
