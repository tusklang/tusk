package initialize

type Package struct {
	Name   string
	Files  []*File
	parent *Package
}

func (p Package) Parent() *Package {
	return p.parent
}

func (p Package) FullName() string {
	//returns the full name, with all the parents
	if p.parent == nil {
		return p.Name
	}
	return p.parent.FullName() + p.Name + "."
}
