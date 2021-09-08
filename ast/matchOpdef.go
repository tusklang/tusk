package ast

//the operation definition has a few special keywords to match a group of types
//this function is able to match the typename given as a string to the opdef's typename, which could be a keyword
func matchOpdef(typename, opdef string) bool {
	return typename == opdef //for now this is the only one we need, but there will be more later ;)
}
