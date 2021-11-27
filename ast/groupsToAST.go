package ast

import "github.com/tusklang/tusk/errhandle"

//used to group operations into a tree
/*
pub fn int main() {
	var a = 1 + 3 + 4;
}
becomes:

pub fn main
-> =
	-> var a
	-> +
		-> 1
		-> +
			-> 3
			-> 4
*/

func groupsToAST(items []Group) ([]*ASTNode, *errhandle.TuskError) {

	var opList = []map[string]func(exp []Group, index int) ([]*ASTNode, *errhandle.TuskError){
		{
			";":   termOpHandle,
			",":   termOpHandle,
			"ltr": nil, //ltr associativity
		},
		{
			"=":   defaultOperationHandle,
			"ltr": nil,
		},
		{
			":":   defaultOperationHandle,
			"ltr": nil,
		},
		{
			"&":   defaultOperationHandle,
			"ltr": nil,
		},
		{
			"|":   defaultOperationHandle,
			"ltr": nil,
		},
		{
			"^":   defaultOperationHandle,
			"ltr": nil,
		},
		{
			"==":  defaultOperationHandle,
			"!=":  defaultOperationHandle,
			"ltr": nil,
		},
		{
			"<=":  defaultOperationHandle,
			"<":   defaultOperationHandle,
			">":   defaultOperationHandle,
			">=":  defaultOperationHandle,
			"ltr": nil,
		},
		{
			"+":   defaultOperationHandle,
			"-":   defaultOperationHandle,
			"ltr": nil,
		},
		{
			"*":   defaultOperationHandle,
			"/":   defaultOperationHandle,
			"ltr": nil,
		},
		{
			".":   defaultOperationHandle,
			"[]":  nil, //array index
			"()":  nil, //function call
			"ltr": nil,
		},
		{
			"~":   defaultOperationHandle,
			"rtl": nil,
		},
		{
			"@":   defaultOperationHandle,
			"#":   defaultOperationHandle,
			"rtl": nil,
		},
		{
			"->":  defaultOperationHandle,
			"ltr": nil,
		},
		//lower on this list means greater precedence
	}

	//go through all the operation groups
	for _, v := range opList {

		var start int
		var cond func(int) bool
		var inc func(*int)

		if _, ltr := v["ltr"]; ltr {
			//it has ltr assoc
			start = len(items) - 1
			cond = func(i int) bool {
				return i >= 0
			}
			inc = func(i *int) {
				*i--
			}
		} else if _, rtl := v["rtl"]; rtl {
			//it has rtl assoc
			start = 0
			cond = func(i int) bool {
				return i < len(items)
			}
			inc = func(i *int) {
				*i++
			}
		} else {
			//error
			//internal tusk malfunction
		}

		for i := start; cond(i); inc(&i) {

			for k, vv := range v {

				if k == "ltr" || k == "rtl" {
					//associativity props
					continue
				}

				switch g := items[i].(type) {
				case *Operation:
					if g.GetMTok().Name == k {
						return vv(items, i)
					}
				case *Array:

					if k == "[]" && g.Siz != nil && g.Arr == nil && g.Typ == nil {

						if len(items) == 1 {
							return []*ASTNode{{
								Group: items[i],
							}}, nil
						}

						return arrIndexHandle(items, i)
					}

				case *Block:
					//for function calls

					//if the blocktype is a (
					if k == "()" && g.BlockType == "(" {

						if i-1 >= 0 {

							//if the item prior is a function, variable, or block
							/*
								a();
								fn() {}();
								(fn() {})();
							*/
							switch items[i-1].(type) {
							case *VarRef:
								return funcCallHandle(items, i)
							case *Function:
								return funcCallHandle(items, i)
							case *Block:
								return funcCallHandle(items, i)
							}

						}

					}
				}
			}

		}
	}

	var ret = make([]*ASTNode, len(items))

	for k, v := range items {
		ret[k] = &ASTNode{}
		ret[k].Group = v
	}

	//it must be a single, since there is no operation
	return ret, nil
}
