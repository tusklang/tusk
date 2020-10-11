package types

import "fmt"

//TuskError represents an eror in tusk
type TuskError struct {
	Err        string
	Stacktrace []string
}

//Print prints the error (in a formatted way)
func (e *TuskError) Print() {
	fmt.Println(e.Err)
	fmt.Println("When the error was thrown, this was the stack:")
	for k, v := range e.Stacktrace {
		end := "\n"
		if k+1 == len(e.Stacktrace) {
			end = ""
		}
		fmt.Print("  "+v, end)
	}
}

func (e TuskError) Format() string {
	return e.Err
}

func (e TuskError) Type() string {
	return "error"
}

func (e TuskError) TypeOf() string {
	return e.Type()
}

func (e TuskError) Deallocate() {}

func (e TuskError) Range(fn func(val1, val2 *TuskType) (Returner, *TuskError)) (*Returner, *TuskError) {
	return nil, nil
}
