package native

import (
	"fmt"
	"unsafe"

	"github.com/tusklang/tusk/lang/types"
)

//#include <stdbool.h>
//#include "cload_types.h"
import "C"

//converting tusk types to c types
/////////////////////////////////////////////////////////////////////////

//	---------------------------------------
//		string -> char*
//		bool   -> bool
//		number -> double
//		rune   -> int
//		array  -> void**
//	---------------------------------------
//	hashes, protos, and objects are not supported

func cConvertType(val types.TuskType) (unsafe.Pointer, error) {

	switch val.(type) {
	case types.TuskString:
		gostr := val.(types.TuskString).ToGoType()
		ccstr := C.CString(gostr)
		return unsafe.Pointer(&ccstr), nil
	case types.TuskBool:
		cbool := C.bool(val.(types.TuskBool).ToGoType())
		return unsafe.Pointer(&cbool), nil
	case types.TuskNumber:
		gonum := val.(types.TuskNumber).ToGoType()
		cdouble := C.double(gonum)
		return unsafe.Pointer(&cdouble), nil
	case types.TuskRune:
		gorune := val.(types.TuskRune).ToGoType()
		cint := C.int(gorune)
		return unsafe.Pointer(&cint), nil
	case types.TuskArray:

		tuskarr := val.(types.TuskArray)
		var carr = C.allocarr(C.ulonglong(tuskarr.Length))

		var err error

		tuskarr.Range(func(k, v *types.TuskType) types.Returner {
			kint := int((*k).(types.TuskNumber).ToGoType())
			cval, e := cConvertType(*v)

			if e != nil { //if there is an error, give it back to the err var
				err = e
				return types.Returner{
					Type: "break",
				}
			}

			C.setindex(carr, C.int(kint), cval)
			return types.Returner{}
		})

		if err != nil {
			return nil, err
		}

		return unsafe.Pointer(carr), nil
	}

	return nil, fmt.Errorf("Cannot convert a type %s to a C type", val.Type())
}
