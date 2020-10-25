package types

import "fmt"

//TuskError represents an eror in tusk
type TuskError struct {
	Err        string
	Stacktrace []string
}

//Print prints the error (in a formatted way)
func (e *TuskError) Print() {
	fmt.Println("Uncaught Panic:", e.Err)
	fmt.Println("When the error was thrown, this was the stack:")
	for _, v := range e.Stacktrace {
		fmt.Println("  " + v)
	}
}

//Format returns the error message
func (e TuskError) Format() string {
	return e.Err
}

//Type returns the type of object it is
func (e TuskError) Type() string {
	return "error"
}

//TypeOf is equivalent to Type, but if it is a prototype or object, it gives the prototype or object name
func (e TuskError) TypeOf() string {
	return e.Type()
}

//Deallocate defines extra steps to deallocate the type
func (e TuskError) Deallocate() {}

//Clone cannot be used on a TuskError
func (e TuskError) Clone() *TuskType {
	return nil
}

//Range ranges through the variable, does not work for this type
func (e TuskError) Range(fn func(val1, val2 *TuskType) (Returner, *TuskError)) (*Returner, *TuskError) {
	return nil, nil
}
