package ast

import (
	"strings"

	"github.com/tusklang/tusk/data"
)

//the operation definition has a few special keywords to match a group of types
//this function is able to match the typename to the opdef's typename, which could be a keyword
func matchOpdef(val data.Value, opdef string) bool {

	splopdef := strings.SplitN(opdef, "&", 2) //split by the & sign (& sign seperates and clauses)

	if len(splopdef) == 2 {
		return matchOpdef(val, splopdef[0]) && matchOpdef(val, splopdef[1])
	}

	if opdef[0] == '!' {
		//inverse the output
		return !matchOpdef(val, opdef[1:])
	}

	if val == nil {

		//if the operand request is an empty operand
		//we return true
		//otherwise it's false
		return opdef == "-"
	}

	//any type
	if opdef == "*" {
		return true
	}

	//special flags
	if (opdef == "class" || opdef == "instance" || opdef == "ptr" || opdef == "type" || opdef == "var") && val.TypeData().HasFlag(opdef) {
		return true
	}

	//arrays
	if (opdef == "slice" || opdef == "fixed" || opdef == "varied") && val.TypeData().HasFlag(opdef) {
		return true
	}

	return val.TypeData().Name() == opdef //for now this is the only one we need, but there will be more later ;)
}
