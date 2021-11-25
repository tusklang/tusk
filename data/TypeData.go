package data

type TypeData struct {
	nam      string
	flags    map[string]bool
	otherdat map[string]Value
}

func NewTypeData(name string) *TypeData {
	return &TypeData{
		nam:      name,
		flags:    make(map[string]bool),
		otherdat: make(map[string]Value),
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

func (td TypeData) GetOtherDat(name string) Value {
	return td.otherdat[name]
}

func (td *TypeData) AddOtherDat(name string, value Value) {
	td.otherdat[name] = value
}
