package types

type OmmObject struct {
	Name       string
	Instance   Instance
	AccessList map[string][]string
}

func (o OmmObject) Get(field, file string) (*OmmType, error) {

	var mappedvars = make(map[string]*OmmType)

	for k, v := range o.Instance.vars {
		mappedvars[k] = v.Value
	}

	return getfield(mappedvars, field, o.AccessList, file)
}

func (o OmmObject) Format() string {
	return "{" + o.Name[1:] + "}"
}

func (o OmmObject) Type() string {
	return "object"
}

func (o OmmObject) TypeOf() string {
	return o.Name[1:]
}

func (_ OmmObject) Deallocate() {}
