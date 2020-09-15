package interpreter

import (
	. "ka/lang/types"
)

type KaGoFunc struct {
	Function func(args []*KaType, stacktrace []string, line uint64, file string, instance *Instance) *KaType
}

func (ogf KaGoFunc) Format() string {
	return "{ native go func }"
}

func (ogf KaGoFunc) Type() string {
	return "native_func"
}

func (ogf KaGoFunc) TypeOf() string {
	return ogf.Type()
}

func (ogf KaGoFunc) Deallocate() {}

//Range ranges over an ka native function
func (ogf KaGoFunc) Range(fn func(val1, val2 *KaType) Returner) *Returner {
	return nil
}

func nativeinit() {

	//init the simple native values first
	for k, v := range native {
		var gofunc KaType = KaGoFunc{
			Function: v,
		}

		Native["$"+k] = &gofunc
	}

}
