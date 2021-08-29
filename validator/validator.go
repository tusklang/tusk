package validator

import "github.com/tusklang/tusk/operations"

//this package is used to validate variable types, legal usages, and rename local variables to not have duplicates

type Validator struct {
}

func Validate(ops []*operations.Operation) *Validator {
	return nil
}
