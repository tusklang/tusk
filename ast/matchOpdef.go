package ast

import (
	"strings"

	"github.com/tusklang/tusk/data"
)

//the operation definition has a few special keywords to match a group of types
//this function is able to match the typename to the opdef's typename, which could be a keyword
func matchOpdef(val data.Value, opdef string) bool {

	if opdef == "class" && strings.HasPrefix(val.TypeString(), "class ") {
		return true
	}

	if opdef == "instance" && strings.HasPrefix(val.TypeString(), "instance ") {
		return true
	}

	return val.TypeString() == opdef //for now this is the only one we need, but there will be more later ;)
}
