package varprocessor

import (
	"github.com/tusklang/tusk/ast"
	"github.com/tusklang/tusk/errhandle"
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

func (p *VarProcessor) process(tree []*ast.ASTNode, declared map[string]decl, instatic bool) {

	var curscope = make(map[string]decl)

	for k, v := range tree {

		switch g := v.Group.(type) {
		case *ast.VarDecl:

			m := mergemap(declared, curscope)

			if _, exists := m[g.Name]; exists {
				//error
				//variable with that name has already been declared
				p.compiler.AddError(errhandle.NewTuskErrorFTok(
					"duplicated varname",
					"try renaming this variable",
					g.GetMTok(),
				))
			}

			if g.Type != nil {
				p.process([]*ast.ASTNode{g.Type}, m, instatic)
			}
			if g.Value != nil {
				p.process([]*ast.ASTNode{g.Value}, m, instatic)
			}

			nname := p.nextvar()
			curscope[g.Name] = decl{
				nname:  nname,
				static: true,
			}
			g.Name = nname
		case *ast.VarRef:

			//check both the outer declarations and current scope for the variable reference
			d, ex1 := declared[g.Name]
			cs, ex2 := curscope[g.Name]

			if !(ex1 || ex2) {
				//error
				//there isn't a variable declared with that name
				p.compiler.AddError(errhandle.NewTuskErrorFTok(
					"undefined variable",
					"",
					g.GetMTok(),
				))
				continue
			}

			//if the outer scope doesn't include the var ref, it's in the current scope
			if !ex1 {
				d = cs
			}

			if !d.static && instatic {
				p.compiler.AddError(errhandle.NewTuskErrorFTok(
					"accessing an instance member from a static member",
					"try making this global's declaration static",
					g.GetMTok(),
				))
			}

			if d.macro != nil {
				tree[k] = d.macro
			} else {
				g.Name = d.nname //rename the variable in the ast
			}

		case *ast.Block:
			p.process(g.Sub, mergemap(declared, curscope), instatic)
		case *ast.Function:

			m := mergemap(declared, curscope)

			for _, v := range g.Params {
				m[v.Name] = decl{
					nname: p.nextvar(),
				}
				v.Name = m[v.Name].nname
			}

			if g.Body != nil {
				p.process(g.Body.Sub, m, instatic)
			}
		case *ast.Operation:
			m := mergemap(declared, curscope)
			p.process(v.Left, m, instatic)

			//if it's the dot operator, only check the left side
			if g.OpType != "." {
				p.process(v.Right, m, instatic)
			}
		case *ast.IfStatement:

			merged := mergemap(declared, curscope)

			p.process(g.Condition, merged, instatic)
			p.process(g.Body, merged, instatic)
			p.process(g.ElseBody, merged, instatic)

		case *ast.WhileStatement:

			merged := mergemap(declared, curscope)

			p.process(g.Condition, merged, instatic)
			p.process(g.Body, merged, instatic)

		case *ast.Array:
			merged := mergemap(declared, curscope)

			if g.Siz != nil {
				p.process([]*ast.ASTNode{g.Siz}, merged, instatic)
			}
			if g.Typ != nil {
				p.process([]*ast.ASTNode{g.Typ}, merged, instatic)
			}
			p.process(g.Arr, merged, instatic)

		case *ast.Return:

			if g.Val != nil {
				p.process([]*ast.ASTNode{g.Val}, mergemap(declared, curscope), instatic)
			}

		}
	}
}

func (p *VarProcessor) ProcessVars(file *initialize.File) {

	var globals = make(map[string]decl)

	for _, v := range file.Globals {

		if v.CRel == 2 {
			//link variable
			globals[v.Link.TName] = decl{
				nname:  v.Link.TName,
				static: true,
			}
			continue
		}

		var nam string

		if v.Value != nil {
			nam = v.Value.Name
		} else if v.Func != nil {
			nam = v.Func.Name
		}

		//add all the globals
		globals[nam] = decl{
			nname:  nam,
			static: v.CRel != 0,
		}
	}

	for _, v := range file.Globals {

		if v.Func != nil {

			if v.Func.Body.Sub == nil {
				//function has no body
				continue
			}

			m := mergemap(p.predecl, globals)

			//process the function params
			for _, vv := range v.Func.Params {
				m[vv.Name] = decl{
					nname: p.nextvar(),
				}
				vv.Name = m[vv.Name].nname
			}

			p.process(v.Func.Body.Sub, m, v.CRel == 1) //process the function body
		} else if v.Value != nil {

			if v.Value.Value == nil {
				//global declared with no value
				//e.g. var a: i32
				//vs
				//var a: i32 = 0;
				continue
			}

			p.process([]*ast.ASTNode{v.Value.Value}, mergemap(p.predecl, globals), v.CRel == 1) //process the declaration's assigned value
		}

	}

}
