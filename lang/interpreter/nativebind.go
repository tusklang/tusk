package interpreter

import (
	. "github.com/omm-lang/omm/lang/types"
	"github.com/omm-lang/omm/stdlib/native"
)

func nativebind() {

	//init the simple native values first
	for k, v := range simplenative {
		var gofunc OmmType = OmmGoFunc{
			Function: v,
		}

		Native[k] = &gofunc
	}

	//now do the complex ones
	complexnative, nativeops := native.NativeStd()

	for k, v := range complexnative {
		Native[k] = v
	}
	for k, v := range nativeops {
		Operations[k] = v
	}
}
