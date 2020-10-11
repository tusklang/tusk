package native

import (
	. "github.com/tusklang/tusk/lang/types"
)

type TuskGoFunc struct {
	Function   func(args []*TuskType, stacktrace []string, line uint64, file string, instance *Instance) (*TuskType, *TuskError)
	Signatures [][]string
}

func (tgf TuskGoFunc) Format() string {
	return "{ native go func }"
}

func (tgf TuskGoFunc) Type() string {
	return "native_func"
}

func (tgf TuskGoFunc) TypeOf() string {
	return tgf.Type()
}

func (tgf TuskGoFunc) Deallocate() {}

//Range ranges over an tusk native function
func (ogf TuskGoFunc) Range(fn func(val1, val2 *TuskType) (Returner, *TuskError)) (*Returner, *TuskError) {
	return nil, nil
}

func init() {

	//init the simple native values first
	for k, v := range NativeFuncs {
		var tmp TuskType = v
		Native[k] = &tmp
	}

}
