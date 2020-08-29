package native

import (
	"reflect"

	"github.com/omm-lang/omm/lang/types"
)

/*
	GoatProtoIndex is the default <custom go struct> :: string operation
	To use it, you can create an operation function:
	func(val1, val2 types.OmmType, instance *types.Instance, stacktrace []string, line uint64, file string, stacksize uint) *types.OmmType {
		asserted := val1.(<your type>)
		return native.GoatProtoIndex(reflect.ValueOf(&asserted), val2.(types.OmmString), stacktrace, line, file)
	}
*/
func GoatProtoIndex(val1 reflect.Value, val2 types.OmmString, stacktrace []string, line uint64, file string) *types.OmmType {

	key := val2.ToGoType()
	field := val1.Elem().FieldByName(key)

	if !field.IsValid() {
		OmmPanic("Type does not contain field "+key, line, file, stacktrace)
	}

	var ommtype = field.Interface().(types.OmmType)
	return &ommtype
}
