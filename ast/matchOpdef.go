package ast

import (
	"github.com/tusklang/tusk/data"
)

//the operation definition has a few special keywords to match a group of types
//this function is able to match the typename to the opdef's typename, which could be a keyword
func matchOpdef(val data.Value, opdef string) bool {

	if opdef == "*" {
		return true
	}

	if (opdef == "class" || opdef == "instance" || opdef == "ptr") && val.TypeData().HasFlag(opdef) {
		return true
	}

	return val.TypeData().Name() == opdef //for now this is the only one we need, but there will be more later ;)
}
