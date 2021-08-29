package validator

import "github.com/tusklang/tusk/ast"

//this package is used to validate variable types, legal usages, and rename local variables to not have duplicates

type Validator struct {
}

func Validate(ops []*ast.ASTNode) *Validator {
	return nil
}
