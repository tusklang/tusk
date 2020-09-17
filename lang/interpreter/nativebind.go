package interpreter

import (
	. "github.com/tusklang/tusk/lang/types"
)

type TuskGoFunc struct {
	Function func(args []*TuskType, stacktrace []string, line uint64, file string, instance *Instance) *TuskType
}

func (ogf TuskGoFunc) Format() string {
	return "{ native go func }"
}

func (ogf TuskGoFunc) Type() string {
	return "native_func"
}

func (ogf TuskGoFunc) TypeOf() string {
	return ogf.Type()
}

func (ogf TuskGoFunc) Deallocate() {}

//Range ranges over an tusk native function
func (ogf TuskGoFunc) Range(fn func(val1, val2 *TuskType) Returner) *Returner {
	return nil
}

func nativeinit() {

	//init the simple native values first
	for k, v := range native {
		var gofunc TuskType = TuskGoFunc{
			Function: v,
		}

		Native["$"+k] = &gofunc
	}

}
