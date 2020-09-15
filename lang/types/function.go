package types

type Overload struct {
	Params  []string
	Types   []string
	Body    []Action
	VarRefs []string //variables that this function use
}

type KaFunc struct {
	Overloads []Overload
	Instance  *Instance
}

func (f KaFunc) Format() string {
	return "(function) { ... }"
}

func (f KaFunc) Type() string {
	return "function"
}

func (f KaFunc) TypeOf() string {
	return f.Type()
}

func (f KaFunc) Deallocate() {}

//Range ranges over a function
func (f KaFunc) Range(fn func(val1, val2 *KaType) Returner) *Returner {
	return nil
}
