package types

type Overload struct {
	Params  []string
	Types   []string
	Body    []Action
	VarRefs []string //variables that this function use
}

type OmmFunc struct {
	Overloads []Overload
	Instance  *Instance
}

func (f OmmFunc) Format() string {
	return "(function) { ... }"
}

func (f OmmFunc) Type() string {
	return "function"
}

func (f OmmFunc) TypeOf() string {
	return f.Type()
}

func (f OmmFunc) Deallocate() {}

//Range ranges over a function
func (f OmmFunc) Range(fn func(val1, val2 *OmmType) Returner) *Returner {
	return nil
}
