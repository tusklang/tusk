package varprocessor

import (
	"github.com/tusklang/tusk/ast"
	"github.com/tusklang/tusk/initialize"
)

//decl structure used to store variable declarations
type decl struct {
	nname  string //new name
	static bool
}

//util function to merge two varmaps
func mergemap(m1, m2 map[string]decl) (fin map[string]decl) {

	fin = make(map[string]decl)

	for k, v := range m1 {
		fin[k] = v
	}

	for k, v := range m2 {
		fin[k] = v
	}

	return
}

func (p *VarProcessor) process(tree []*ast.ASTNode, declared map[string]decl) {

	var curscope = make(map[string]decl)

	for _, v := range tree {

		switch g := v.Group.(type) {
		case *ast.VarDecl:

			if _, exists := mergemap(declared, curscope)[g.Name]; exists {
				//error
				//variable with that name has already been declared
			}

			nname := p.nextvar()
			curscope[g.Name] = decl{
				nname:  nname,
				static: false,
			}
			g.Name = nname
		case *ast.VarRef:

			//check both the outer declarations and current scope for the variable reference
			d, ex1 := declared[g.Name]
			cs, ex2 := curscope[g.Name]

			if !(ex1 || ex2) {
				//error
				//there isn't a variable declared with that name
			}

			//if the outer scope doesn't include the var ref, it's in the current scope
			if !ex1 {
				d = cs
			}

			g.Name = d.nname //rename the variable in the ast
		case *ast.Block:
			p.process(g.Sub, mergemap(declared, curscope))
		case *ast.Function:
			p.process(g.Body.Sub, mergemap(declared, curscope))
		case *ast.Operation:
			m := mergemap(declared, curscope)
			p.process(v.Left, m)
			p.process(v.Right, m)
		}
	}
}

func (p *VarProcessor) ProcessVars(file *initialize.File) {

	var globals = make(map[string]decl)

	for _, v := range file.Globals {
		//add all the globals
		globals[v.Value.Name] = decl{
			nname:  v.Value.Name,
			static: v.IsStatic,
		}
	}

	for _, v := range file.Globals {
		p.process([]*ast.ASTNode{v.Value.Value}, globals) //process the declaration's assigned value
	}

}
