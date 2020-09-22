package native

import (
	"unsafe"

	"github.com/tusklang/tusk/lang/types"
)

//#include <stdbool.h>
//#include "cload.h"
//#include "cload_types.h"
import "C"

//TuskCLib represents a c (.dll or .so) in tusk
type TuskCLib struct {
	lib C.struct_CLib
}

func (lib TuskCLib) Format() string {
	return "{ c-library }"
}

func (lib TuskCLib) Type() string {
	return "clib"
}

func (lib TuskCLib) TypeOf() string {
	return lib.Type()
}

//Deallocate frees a C Library
func (lib TuskCLib) Deallocate() {
	C.freelib(lib.lib)
}

func (lib TuskCLib) Range(fn func(val1, val2 *types.TuskType) types.Returner) *types.Returner {
	return nil
}

//TuskCProc represents a c (.dll or .so) function in tusk
type TuskCProc struct {
	proc C.struct_CProc
}

func (proc TuskCProc) Format() string {
	return "{ c-proc }"
}

func (proc TuskCProc) Type() string {
	return "cproc"
}

func (proc TuskCProc) TypeOf() string {
	return proc.Type()
}

func (proc TuskCProc) Deallocate() {}

//Range is a dummy range function in order to make TuskCProc implement TuskType
func (proc TuskCProc) Range(fn func(val1, val2 *types.TuskType) types.Returner) *types.Returner {
	return nil
}

func cload(args []*types.TuskType, stacktrace []string, line uint64, file string, instance *types.Instance) *types.TuskType {

	if len(args) != 1 || (*args[0]).Type() != "string" {
		TuskPanic("Function cload requires the signature (string)", line, file, stacktrace)
	}

	cname := (*args[0]).(types.TuskString).ToGoType()
	lib := C.cloadlib(C.CString(cname))

	if lib.err {
		TuskPanic("Could not open "+cname, line, file, stacktrace)
	}

	var clib TuskCLib
	clib.lib = lib

	var tusktype types.TuskType = clib
	return &tusktype
}

//CloadGetProc is an operation function that represents a clib::cproc
func CloadGetProc(val1, val2 types.TuskType, instance *types.Instance, stacktrace []string, line uint64, file string, stacksize uint) *types.TuskType {
	lib := val1.(TuskCLib).lib
	name := val2.(types.TuskString).ToGoType()

	cproc := C.cgetproc(lib, C.CString(name))

	var tuskcproc TuskCProc
	tuskcproc.proc = cproc
	var tusktype types.TuskType = tuskcproc
	return &tusktype
}

//CloadCallProc is an operation function that represents a cproc:array
func CloadCallProc(val1, val2 types.TuskType, instance *types.Instance, stacktrace []string, line uint64, file string, stacksize uint) *types.TuskType {

	//convert the second param to a c pointer array

	length := val2.(types.TuskArray).Length
	tmp, e := cConvertType(val2)

	if e != nil {
		TuskPanic(e.Error(), line, file, stacktrace)
	}

	gotypePtr := (*unsafe.Pointer)(tmp)
	ret := float64(C.ccallproc(val1.(TuskCProc).proc, gotypePtr, C.int(length)))

	//return as a number
	var tusknum types.TuskNumber
	tusknum.FromGoType(ret)
	var tusktype types.TuskType = tusknum
	return &tusktype
}

func initcload() {

}
