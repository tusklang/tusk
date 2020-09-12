package native

import (
	"reflect"

	"omm/lang/types"
	. "omm/lang/types"
)

type OmmGoFunc struct {
	Function func(args []*OmmType, stacktrace []string, line uint64, file string, instance *Instance) *OmmType
}

func (ogf OmmGoFunc) Format() string {
	return "{ native go func }"
}

func (ogf OmmGoFunc) Type() string {
	return "native_func"
}

func (ogf OmmGoFunc) TypeOf() string {
	return ogf.Type()
}

func (ogf OmmGoFunc) Deallocate() {}

//Range ranges over an omm native function
func (ogf OmmGoFunc) Range(fn func(val1, val2 *OmmType) Returner) *Returner {
	return nil
}

var native = make(map[string]*OmmType)
var operations = make(map[string]func(val1, val2 OmmType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *OmmType)

//the simple native functions are stored in lang/interpreter/nativefuncs.go

func putfunc(name string, fn *func([]*OmmType, []string, uint64, string, *Instance) *OmmType) {
	var ogf OmmType = OmmGoFunc{
		Function: *fn,
	}
	native[name] = &ogf
}

func putoper(name string, fn func(val1, val2 types.OmmType, instance *types.Instance, stacktrace []string, line uint64, file string, stacksize uint) *types.OmmType) {
	operations[name] = fn
}

//GetStd returns all of the more complex natives (they are not all functions)
func GetStd() (map[string]*OmmType, map[string]func(val1, val2 OmmType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *OmmType) {

	//return 1 is the native vars, and return 2 is the native operations (operations that the native vars make)

	putfunc("url.request", &urlrequest)
	putoper("http_response :: string", func(val1, val2 types.OmmType, instance *types.Instance, stacktrace []string, line uint64, file string, stacksize uint) *types.OmmType {
		asserted := val1.(OmmURLResp)
		return gostructprotoindex(reflect.ValueOf(&asserted), val2.(types.OmmString), stacktrace, line, file)
	})
	putfunc("runtime.load", &loadlibrary)
	putoper("runtime_lib :: string", func(val1, val2 types.OmmType, instance *types.Instance, stacktrace []string, line uint64, file string, stacksize uint) *types.OmmType {
		var varname = "$" + val2.(OmmString).ToGoType()
		var lib = val1.(OmmRuntimeLibrary)

		v, e := lib.Get(varname)
		if e != nil {
			OmmPanic(e.Error(), line, file, stacktrace)
		}

		return v
	})
	putfunc("instance.setprec", &setprec)
	putfunc("files.open", &omm_newfile)
	putfunc("files.create", &omm_createfile)
	putoper("file :: string", func(val1, val2 types.OmmType, instance *types.Instance, stacktrace []string, line uint64, file string, stacksize uint) *types.OmmType {
		asserted := val1.(OmmFile)
		return gostructprotoindex(reflect.ValueOf(&asserted), val2.(types.OmmString), stacktrace, line, file)
	})

	return native, operations
}
