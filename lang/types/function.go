package types

type Overload struct {
	Params  []string
	Types   []string
	Body    []Action
	VarRefs []string
}

type TuskFunc struct {
	Overloads []Overload
	Instance  *Instance
}

func (f TuskFunc) Format() string {
	return "(function) { ... }"
}

func (f TuskFunc) Type() string {
	return "function"
}

func (f TuskFunc) TypeOf() string {
	return f.Type()
}

func (f TuskFunc) Deallocate() {}

//Clone cannot be used on a TuskFunction
func (f TuskFunc) Clone() *TuskType {
	return nil
}

//Range ranges over a function
func (f TuskFunc) Range(fn func(val1, val2 *TuskType) (Returner, *TuskError)) (*Returner, *TuskError) {
	return nil, nil
}
