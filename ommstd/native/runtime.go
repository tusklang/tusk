package native

import (
	"errors"

	oatenc "github.com/omm-lang/oat/format/encoding"
	"github.com/omm-lang/omm/lang/types"
)

//OmmRuntimeLibrary stores an oat binary at runtime
type OmmRuntimeLibrary struct {
	instance *types.Instance
}

var loadlibrary = func(args []*types.OmmType, stacktrace []string, line uint64, file string, instance *types.Instance) *types.OmmType {

	if len(args) != 1 || (*args[0]).Type() != "string" {
		OmmPanic("Function runtime.load requires a parameter count of 1 with the type string", line, file, stacktrace)
	}

	filename := (*args[0]).(types.OmmString).ToGoType()
	lib, e := oatenc.OatDecode(filename)

	if e != nil {
		OmmPanic(e.Error(), line, file, stacktrace)
	}

	var libInstance = &types.Instance{}
	libInstance.Params = instance.Params

	for k, v := range lib {
		libInstance.Allocate(k, v)
	}

	var runtimelib types.OmmType = OmmRuntimeLibrary{
		instance: libInstance,
	}
	return &runtimelib
}

func (l OmmRuntimeLibrary) Get(v string) (*types.OmmType, error) {

	vari := (*l.instance).Fetch(v)

	if vari == nil { //does not exist
		return nil, errors.New("Library does not contain global: " + v)
	}

	return vari.Value, nil
}

func (l OmmRuntimeLibrary) Format() string {
	return "{ runtime library }"
}

func (l OmmRuntimeLibrary) Type() string {
	return "runtime_lib"
}

func (l OmmRuntimeLibrary) TypeOf() string {
	return l.Type()
}

func (_ OmmRuntimeLibrary) Deallocate() {}
