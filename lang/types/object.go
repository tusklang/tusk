package types

type TuskObject struct {
	Name       string
	Instance   Instance
	AccessList map[string][]string
}

func (o TuskObject) Get(field, file string) (*TuskType, error) {

	var mappedvars = make(map[string]*TuskType)

	for k, v := range o.Instance.vars {
		mappedvars[k] = v.Value
	}

	return getfield(mappedvars, field, o.AccessList, file)
}

func (o TuskObject) Format() string {
	return "{" + o.Name + "}"
}

func (o TuskObject) Type() string {
	return "object"
}

func (o TuskObject) TypeOf() string {
	return o.Name[1:]
}

func (o TuskObject) Deallocate() {}

func (o TuskObject) Clone() *TuskType {
	return nil
}

//Range ranges over an object
func (o TuskObject) Range(fn func(val1, val2 *TuskType) (Returner, *TuskError)) (*Returner, *TuskError) {
	return nil, nil
}
