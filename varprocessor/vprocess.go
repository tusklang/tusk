package varprocessor

import (
	"github.com/tusklang/tusk/ast"
	"github.com/tusklang/tusk/initialize"
)

//decl structure used to store variable declarations
type decl struct {
	nname  string       //new name
	macro  *ast.ASTNode //or if we replace the old name with a macro
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

	for k, v := range tree {

		switch g := v.Group.(type) {
		case *ast.VarDecl:

			m := mergemap(declared, curscope)

			if _, exists := m[g.Name]; exists {
				//error
				//variable with that name has already been declared
			}

			if g.Type != nil {
				p.process([]*ast.ASTNode{g.Type}, m)
			}
			if g.Value != nil {
				p.process([]*ast.ASTNode{g.Value}, m)
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

			if d.macro != nil {
				tree[k] = d.macro
			} else {
				g.Name = d.nname //rename the variable in the ast
			}

		case *ast.Block:
			p.process(g.Sub, mergemap(declared, curscope))
		case *ast.Function:

			m := mergemap(declared, curscope)

			for _, v := range g.Params {
				m[v.Name] = decl{
					nname: p.nextvar(),
				}
				v.Name = m[v.Name].nname
			}

			if g.Body != nil {
				p.process(g.Body.Sub, m)
			}
		case *ast.Operation:
			m := mergemap(declared, curscope)
			p.process(v.Left, m)

			//if it's the dot operator, only check the left side
			if g.OpType != "." {
				p.process(v.Right, m)
			}
		case *ast.IfStatement:

			merged := mergemap(declared, curscope)

			p.process(g.Condition, merged)
			p.process(g.Body, merged)
			p.process(g.ElseBody, merged)

		case *ast.WhileStatement:

			merged := mergemap(declared, curscope)

			p.process(g.Condition, merged)
			p.process(g.Body, merged)

		}
	}
}

func (p *VarProcessor) ProcessVars(file *initialize.File) {

	var globals = make(map[string]decl)

	for _, v := range file.Globals {

		if v.CRel == 2 {
			//link variable
			globals[v.Link.TName] = decl{
				nname: v.Link.TName,
			}
			continue
		}

		//add all the globals
		globals[v.Value.Name] = decl{
			nname:  v.Value.Name,
			static: v.CRel == 1,
		}
	}

	for _, v := range file.Globals {

		if v.Value == nil || v.Value.Value == nil {
			continue
		}

		p.process([]*ast.ASTNode{v.Value.Value}, mergemap(p.predecl, globals)) //process the declaration's assigned value
	}

}
