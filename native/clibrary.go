package native

import "unsafe"
import "github.com/tusklang/tusk/lang/types"

//#cgo !windows LDFLAGS: -ldl
//#include "openlib.h"
import "C"

func makeCargs(vals []*types.TuskType) ([]unsafe.Pointer, error) {
	var cargs = make([]unsafe.Pointer, int(C.MAX_SYS_ARGC))

	i := 0

	for ; i < len(vals) && i < int(C.MAX_SYS_ARGC); i++ {
		var e error
		cargs[i], e = makectype(vals[i])
		if e != nil {
			return nil, e
		}
	}

	for ; i < int(C.MAX_SYS_ARGC); i++ {
		cnum := C.longlong(0)
		cargs[i] = C.makeunsafell(cnum)
	}

	return cargs, nil
}

type CLibrary struct {
	lib C.struct_TUSK_LIB
}

func LoadLib(file string) CLibrary {
	lib := C.loadlib(C.CString(file))

	var libtype CLibrary
	libtype.lib = lib
	return libtype
}

func (l CLibrary) GetProc(name string) CProc {
	var ret CProc
	ret.proc = C.loadproc(l.lib, C.CString(name))
	return ret
}

func (l CLibrary) Format() string {
	return "{c library}"
}

func (l CLibrary) Type() string {
	return "clib"
}

func (l CLibrary) TypeOf() string {
	return l.Type()
}

func (l CLibrary) Deallocate() {
	C.closelib(l.lib)
}

//Clone clones the value into a new pointer
func (l CLibrary) Clone() *types.TuskType {
	return nil
}

//Range ranges over a c library
func (l CLibrary) Range(fn func(val1, val2 *types.TuskType) (types.Returner, *types.TuskError)) (*types.Returner, *types.TuskError) {
	return nil, nil
}

type CProc struct {
	proc C.struct_TUSK_CPROC
}

func (p CProc) Call(vals []*types.TuskType) (int64, error) {
	cargs, e := makeCargs(vals)

	if e != nil {
		return -1, e
	}

	ret := C.callproc(
		p.proc,
		cargs[0],
		cargs[1],
		cargs[2],
		cargs[3],
		cargs[4],
		cargs[5],
		cargs[6],
		cargs[7],
		cargs[8],
		cargs[9],
		cargs[10],
		cargs[11],
		cargs[12],
		cargs[13],
		cargs[14],
		cargs[15],
		cargs[16],
		cargs[17],
		cargs[18],
		cargs[19],
		cargs[20],
	)

	return int64(ret), nil
}

func (p CProc) Format() string {
	return "{c procedure}"
}

func (p CProc) Type() string {
	return "cproc"
}

func (p CProc) TypeOf() string {
	return p.Type()
}

func (p CProc) Deallocate() {}

//Clone clones the value into a new pointer
func (p CProc) Clone() *types.TuskType {
	return nil
}

//Range ranges over a c library function
func (p CProc) Range(fn func(val1, val2 *types.TuskType) (types.Returner, *types.TuskError)) (*types.Returner, *types.TuskError) {
	return nil, nil
}
