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
}

func NewPackage(name, fullname string) *Package {
	return &Package{
		PackageName: name,
		FullName:    fullname,
		Classes:     make(map[string]*Class),
		ChildPacks:  make(map[string]*Package),
	}
}

func (p *Package) AddClass(name string, class *Class) {
	p.Classes[name] = class
}

func (p *Package) LLVal(block *ir.Block) value.Value {
	return nil
}

func (p *Package) Type() types.Type {
	return nil
}

func (p *Package) TypeString() string {
	return "package"
}
