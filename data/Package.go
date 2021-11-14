package data

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

type Package struct {
	PackageName, FullName string
	Classes               map[string]*Class
	ChildPacks            map[string]*Package

	parent *Package
}

func NewPackage(name, fullname string, parent *Package) *Package {
	return &Package{
		PackageName: name,
		FullName:    fullname,
		Classes:     make(map[string]*Class),
		ChildPacks:  make(map[string]*Package),
		parent:      parent,
	}
}

func (p *Package) AddClass(name string, class *Class) {
	p.Classes[name] = class
}

func (p *Package) RemParent() {
	p.parent = nil
}

//list all the parents
func (p *Package) ReferenceFromStart() []*Package {
	if p.parent == nil {
		return []*Package{p}
	}

	return append(p.parent.ReferenceFromStart(), p)
}

func (p *Package) LLVal(block *ir.Block) value.Value {
	return nil
}

func (p *Package) TType() Type {
	return nil
}

func (p *Package) Type() types.Type {
	return nil
}

func (p *Package) TypeData() *TypeData {
	return NewTypeData("package")
}

func (p *Package) InstanceV() value.Value {
	return nil
}
