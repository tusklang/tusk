package data

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

type Package struct {
	PackageName   string
	Classes       map[string]*Class
	ParentPackage *Package
}

func NewPackage(name string, parent *Package) *Package {
	return &Package{
		PackageName:   name,
		Classes:       make(map[string]*Class),
		ParentPackage: parent,
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
