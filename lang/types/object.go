package types

type KaObject struct {
	Name       string
	Instance   Instance
	AccessList map[string][]string
}

func (o KaObject) Get(field, file string) (*KaType, error) {

	var mappedvars = make(map[string]*KaType)

	for k, v := range o.Instance.vars {
		mappedvars[k] = v.Value
	}

	return getfield(mappedvars, field, o.AccessList, file)
}

func (o KaObject) Format() string {
	return "{" + o.Name[1:] + "}"
}

func (o KaObject) Type() string {
	return "object"
}

func (o KaObject) TypeOf() string {
	return o.Name[1:]
}

func (o KaObject) Deallocate() {}

//Range ranges over an object
func (o KaObject) Range(fn func(val1, val2 *KaType) Returner) *Returner {
	return nil
}
