package native

import (
	"reflect"
	"unicode"

	"omm/lang/types"
)

func gostructprotoindex(val1 reflect.Value, val2 types.OmmString, stacktrace []string, line uint64, file string) *types.OmmType {

	key := val2.ToGoType()

	if !unicode.IsUpper(rune(key[0])) { //if it is unexported
		OmmPanic("Cannot access unexported field: "+key, line, file, stacktrace)
	}

	field := val1.Elem().FieldByName(key)

	if !field.IsValid() {
		OmmPanic("Type does not contain field "+key, line, file, stacktrace)
	}

	var ommtype = field.Interface().(types.OmmType)
	return &ommtype
}
