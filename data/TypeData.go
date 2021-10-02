package data

type TypeData struct {
	nam   string
	flags map[string]bool
}

func NewTypeData(name string) *TypeData {
	return &TypeData{
		nam:   name,
		flags: make(map[string]bool),
	}
}

func (td TypeData) Name() string {
	return td.nam
}

func (td *TypeData) AddFlag(name string) {
	td.flags[name] = true
}

func (td TypeData) HasFlag(name string) bool {
	return td.flags[name]
}
