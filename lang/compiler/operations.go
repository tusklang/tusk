package compiler

type Operation struct {
	Type   string
	File   string
	Line   uint64
	Left   *Operation
	Right  *Operation
	Degree *Operation
	Item   Item //in case there is no operation
}

var operations []map[string]func(exp []Item, index int, opType string) (Operation, error)

func operationIncludes(group []Item) bool {
	for _, v := range operations {
		for k := range v {
			for _, i := range group {
				if i.Token.Name == k {
					return true
				}
			}
		}
	}

	return false
}

func indexOper_ltr(group []Item, opers []string) (int, int) {

	for k, v := range group {
		for k2, v2 := range opers {
			if v.Token.Name == v2 {
				return k, k2
			}
		}
	}

	return -1, -1
}

func indexOper(group []Item, opers []string) (int, int) {

	for i := len(group) - 1; i >= 0; i-- {
		for k, v := range opers {
			if group[i].Token.Name == v {
				return i, k
			}
		}
	}

	return -1, -1
}

//operation function used for most operators (except assignment, not gate, function calls, etc...)
func normalOpFunc(exp []Item, index int, opType string) (Operation, error) {

	var (
		left  = exp[:index]
		right = exp[index+1:]
	)

	if len(left) == 0 {
		return Operation{}, makeCompilerErr("Must have a value to the left of the "+opType+" operator", exp[index].File, exp[index].Line)
	}
	if len(right) == 0 {
		return Operation{}, makeCompilerErr("Must have a value to the right of the "+opType+" operator", exp[index].File, exp[index].Line)
	}

	leftop, el := makeOperations([][]Item{left})
	rightop, er := makeOperations([][]Item{right})

	if el != nil {
		return Operation{}, el
	}
	if er != nil {
		return Operation{}, er
	}

	return Operation{
		Type:  opType,
		File:  exp[index].File,
		Line:  exp[index].Line,
		Left:  &leftop[0],
		Right: &rightop[0],
	}, nil
}

//ODO is
// ~, ?, =>, and =
// break and continue
// boolean operations (except not gate)
// comparisons
// assignment operations
// exponentiation
// mult, div, modulo
// add, subtract
// index operations and function calls
// not gate
// cast operation

func makeOperations(groups [][]Item) ([]Operation, error) {

	operations = []map[string]func(exp []Item, index int, opType string) (Operation, error){
		map[string]func(exp []Item, index int, opType string) (Operation, error){ //these ones start from left to right
			"~":  normalOpFunc,
			"=":  normalOpFunc,
			":=": normalOpFunc,
			"=>": normalOpFunc,
		},
		map[string]func(exp []Item, index int, opType string) (Operation, error){
			"break": func(exp []Item, index int, opType string) (Operation, error) {
				return Operation{
					Type: opType,
				}, nil
			},
			"continue": func(exp []Item, index int, opType string) (Operation, error) {
				return Operation{
					Type: opType,
				}, nil
			},
		},
		map[string]func(exp []Item, index int, opType string) (Operation, error){
			"&": normalOpFunc,
			"|": normalOpFunc,
		},
		map[string]func(exp []Item, index int, opType string) (Operation, error){
			"==": normalOpFunc,
			"!=": normalOpFunc,
			">":  normalOpFunc,
			">=": normalOpFunc,
			"<":  normalOpFunc,
			"<=": normalOpFunc,
		},
		map[string]func(exp []Item, index int, opType string) (Operation, error){
			"++": func(exp []Item, index int, opType string) (Operation, error) {

				if len(exp[:index]) == 0 {
					return Operation{}, makeCompilerErr("Must have a value to the left of the "+opType+" operator", exp[index].File, exp[index].Line)
				}

				left, e := makeOperations([][]Item{exp[:index]})

				if e != nil {
					return Operation{}, e
				}

				return Operation{
					Type: opType,
					File: exp[index].File,
					Line: exp[index].Line,
					Left: &left[0],
				}, nil
			},
			"--": func(exp []Item, index int, opType string) (Operation, error) {

				if len(exp[:index]) == 0 {
					return Operation{}, makeCompilerErr("Must have a value to the left of the "+opType+" operator", exp[index].File, exp[index].Line)
				}

				left, e := makeOperations([][]Item{exp[:index]})

				if e != nil {
					return Operation{}, e
				}

				return Operation{
					Type: opType,
					File: exp[index].File,
					Line: exp[index].Line,
					Left: &left[0],
				}, nil
			},
			"+=": normalOpFunc,
			"-=": normalOpFunc,
			"*=": normalOpFunc,
			"/=": normalOpFunc,
			"%=": normalOpFunc,
			"^=": normalOpFunc,
		},
		map[string]func(exp []Item, index int, opType string) (Operation, error){
			"^": normalOpFunc,
		},
		map[string]func(exp []Item, index int, opType string) (Operation, error){
			"*": normalOpFunc,
			"/": normalOpFunc,
			"%": normalOpFunc,
		},
		map[string]func(exp []Item, index int, opType string) (Operation, error){
			"+": normalOpFunc,
			"-": normalOpFunc,
		},
		map[string]func(exp []Item, index int, opType string) (Operation, error){
			"::": normalOpFunc,
			"<-": normalOpFunc,
			"<~": normalOpFunc,
		},
		map[string]func(exp []Item, index int, opType string) (Operation, error){
			"!": func(exp []Item, index int, opType string) (Operation, error) {

				if len(exp[index+1:]) == 0 {
					return Operation{}, makeCompilerErr("Must have a value to the right of the "+opType+" operator", exp[index].File, exp[index].Line)
				}

				right, e := makeOperations([][]Item{exp[index+1:]})

				if e != nil {
					return Operation{}, e
				}

				return Operation{
					Type:  opType,
					File:  exp[index].File,
					Line:  exp[index].Line,
					Right: &right[0],
				}, nil
			},
		},
		map[string]func(exp []Item, index int, opType string) (Operation, error){
			"->": normalOpFunc,
		},
	}

	var newGroups []Operation

	for _, v := range groups {

		if !operationIncludes(v) {
			newGroups = append(newGroups, Operation{
				Type: "none",
				File: v[0].File,
				Line: v[0].Line,
				Item: v[0],
			})
		}

		for k, val := range operations {

			var opers []string
			var funcs []func(exp []Item, index int, opType string) (Operation, error)

			for oper, function := range val {
				opers = append(opers, oper)
				funcs = append(funcs, function)
			}

			var indexOfOper int
			var operNum int

			if k == 0 { //the first one should go from left to right
				indexOfOper, operNum = indexOper_ltr(v, opers)
			} else {
				indexOfOper, operNum = indexOper(v, opers)
			}

			if indexOfOper != -1 {
				operation, e := funcs[operNum](v, indexOfOper, opers[operNum])
				if e != nil {
					return []Operation{}, e
				}
				newGroups = append(newGroups, operation)
				goto breakOuter
			}

		}

	breakOuter:
	}

	return newGroups, nil
}
