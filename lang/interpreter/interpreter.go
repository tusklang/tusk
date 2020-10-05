package interpreter

import (
	. "github.com/tusklang/tusk/lang/types"
	. "github.com/tusklang/tusk/native"
)

const MAX_STACKSIZE = 100001

func dealloc(ins *Instance, varnames []string, value *TuskType) { //function to remove the variables declared in that scope
	for _, v := range varnames { //dealloc from the stack (deffered garbage collection)
		ins.Deallocate(v)
	}
}

//Interpreter starts the Tusk runtime with a given instance, and actions tree
func Interpreter(ins *Instance, actions []Action, stacktrace []string, stacksize uint, varnames []string /* varnames to deallocate */, expReturn bool) Returner {

	if stacksize > MAX_STACKSIZE {
		TuskPanic("Stack size was exceeded", 0, "none", stacktrace)
	}

	for _, v := range actions {
		switch v.Type {

		case "var":

			interpreted := Interpreter(ins, v.ExpAct, stacktrace, stacksize+1, nil, true)

			varnames = append(varnames, v.Name)
			ins.Allocate(v.Name, interpreted.Exp)

			if expReturn {
				variable := interpreted.Exp
				defer dealloc(ins, varnames, variable)
				return Returner{
					Type: "expression",
					Exp:  variable,
				}
			}

		case "ovld":

			interpreted := Interpreter(ins, v.ExpAct, stacktrace, stacksize+1, nil, true)

			if (*(*ins.Fetch(v.Name)).Value).Type() != "function" {
				*(*ins.Fetch(v.Name)).Value = TuskFunc{}
			}

			var appended_ovld TuskType = TuskFunc{
				Overloads: append((*(*ins.Fetch(v.Name)).Value).(TuskFunc).Overloads, (*interpreted.Exp).(TuskFunc).Overloads[0]),
			}

			ins.Allocate(v.Name, &appended_ovld)

			if expReturn {
				variable := &appended_ovld
				defer dealloc(ins, varnames, variable)
				return Returner{
					Type: "expression",
					Exp:  variable,
				}
			}

		case "declare":

			var tmpundef TuskType = undef

			varnames = append(varnames, v.Name)
			ins.Allocate(v.Name, &tmpundef)

			if expReturn {
				variable := ins.Fetch(v.Name).Value
				defer dealloc(ins, varnames, variable)
				return Returner{
					Type: "expression",
					Exp:  variable,
				}
			}

		case "let":

			interpreted := *Interpreter(ins, v.ExpAct, stacktrace, stacksize+1, nil, true).Exp

			variable := Interpreter(ins, v.First, stacktrace, stacksize+1, nil, true)

			*variable.Exp = interpreted

			if expReturn {
				defer dealloc(ins, varnames, variable.Exp)
				return Returner{
					Type: "expression",
					Exp:  variable.Exp,
				}
			}

		//all of the types
		case "string":
			fallthrough
		case "rune":
			fallthrough
		case "number":
			fallthrough
		case "bool":
			fallthrough
		case "undef":
			fallthrough
		case "c-array":
			fallthrough //compile-time calculated array
		case "c-hash":
			fallthrough //compile-time calculated hash
		case "thread":

			if expReturn {
				defer dealloc(ins, varnames, &v.Value)
				return Returner{
					Type: "expression",
					Exp:  &v.Value,
				}
			}

		case "function":
			//for a function, add the instance
			var nf = v.Value.(TuskFunc)
			nf.Instance = ins
			var tusktype TuskType = nf
			if expReturn {
				defer dealloc(ins, varnames, &tusktype)
				return Returner{
					Type: "expression",
					Exp:  &tusktype,
				}
			}

		//arrays, hashes are a bit different
		case "r-array":

			var nArr = make([]*TuskType, len(v.Array))

			for k, i := range v.Array {
				nArr[k] = Interpreter(ins, i, stacktrace, stacksize+1, nil, true).Exp
			}

			var kaType TuskType = TuskArray{
				Array:  nArr,
				Length: uint64(len(v.Array)),
			}

			if expReturn {
				defer dealloc(ins, varnames, &kaType)
				return Returner{
					Type: "expression",
					Exp:  &kaType,
				}
			}

		case "r-hash":

			var nHash = make(map[string]*TuskType)

			for _, i := range v.Hash {
				nHash[(*Interpreter(ins, i[0], stacktrace, stacksize+1, nil, true).Exp).Format()] = Interpreter(ins, i[1], stacktrace, stacksize+1, nil, true).Exp
			}

			var kaType TuskType = TuskHash{
				Hash:   nHash,
				Length: uint64(len(v.Hash)),
			}

			if expReturn {
				defer dealloc(ins, varnames, &kaType)
				return Returner{
					Type: "expression",
					Exp:  &kaType,
				}
			}

		////////////////////////////////////

		case "variable":

			_fetched := ins.Fetch(v.Name)

			if _fetched == nil {
				//if it is a nil pointer (only happens because tusk does not support closures)
				TuskPanic("Invalid memory address (nil pointer reference)", v.Line, v.File, stacktrace)
			}

			fetched := _fetched.Value
			if expReturn {
				defer dealloc(ins, varnames, fetched)
				return Returner{
					Type: "expression",
					Exp:  fetched,
				}
			}

		case "{":

			groupRet := Interpreter(ins, v.ExpAct, stacktrace, stacksize+1, nil, false)

			if expReturn {
				defer dealloc(ins, varnames, groupRet.Exp)
				return Returner{
					Type: "expression",
					Exp:  groupRet.Exp,
				}
			}

		case "cast":

			casted := Cast(*Interpreter(ins, v.ExpAct, stacktrace, stacksize+1, nil, true).Exp, v.Name, stacktrace, v.Line, v.File)

			if expReturn {
				defer dealloc(ins, varnames, casted)
				return Returner{
					Type: "expression",
					Exp:  casted,
				}
			}

		//operations
		case "+":
			fallthrough
		case "-":
			fallthrough
		case "*":
			fallthrough
		case "/":
			fallthrough
		case "%":
			fallthrough
		case "^":
			fallthrough
		case "==":
			fallthrough
		case "!=":
			fallthrough
		case ">":
			fallthrough
		case "<":
			fallthrough
		case ">=":
			fallthrough
		case "<=":
			fallthrough
		case "!":
			fallthrough
		case "::":
			fallthrough
		case "|":
			fallthrough
		case "&":
			fallthrough
		case ":":
			fallthrough
		case "?":

			var firstInterpreted Returner
			var secondInterpreted Returner

			//& and | can be a bit different for bool and bool
			if v.Type == "&" || v.Type == "|" {
				firstInterpreted = Interpreter(ins, v.First, stacktrace, stacksize+1, nil, true)
				var computed TuskType

				if (*firstInterpreted.Exp).Type() == "bool" {

					var assumeVal bool

					if v.Type == "&" {
						assumeVal = !isTruthy(*firstInterpreted.Exp)
					} else {
						assumeVal = isTruthy(*firstInterpreted.Exp)
					}

					if assumeVal {
						var retVal = v.Type == "|"
						computed = TuskBool{
							Boolean: &retVal,
						}
					} else {
						secondInterpreted = Interpreter(ins, v.Second, stacktrace, stacksize+1, nil, true)

						if (*secondInterpreted.Exp).Type() == "bool" {
							var gobool = isTruthy(*secondInterpreted.Exp)
							computed = TuskBool{
								Boolean: &gobool,
							}
						} else {
							goto and_or_skip
						}

					}

					if expReturn {
						defer dealloc(ins, varnames, &computed)
						return Returner{
							Type: "expression",
							Exp:  &computed,
						}
					}
				}
			}
			//////////////////////////////////////////////////'

			firstInterpreted = Interpreter(ins, v.First, stacktrace, stacksize+1, nil, true)
			secondInterpreted = Interpreter(ins, v.Second, stacktrace, stacksize+1, nil, true)

		and_or_skip:
			//check typeof first
			operationFunc, exists := Operations[(*firstInterpreted.Exp).TypeOf()+" "+v.Type+" "+(*secondInterpreted.Exp).TypeOf()]

			if !exists { //if it does not exist, also check the type
				operationFunc, exists = Operations[(*firstInterpreted.Exp).Type()+" "+v.Type+" "+(*secondInterpreted.Exp).Type()]
			}

			if !exists { //if there is no operation for that type, panic
				TuskPanic("Could not find "+v.Type+" operator for types "+(*firstInterpreted.Exp).TypeOf()+" and "+(*secondInterpreted.Exp).TypeOf(), v.Line, v.File, stacktrace)
			}

			computed := operationFunc(*firstInterpreted.Exp, *secondInterpreted.Exp, ins, stacktrace, v.Line, v.File, stacksize+1)

			if expReturn {
				defer dealloc(ins, varnames, computed)
				return Returner{
					Type: "expression",
					Exp:  computed,
				}
			}

		////////////

		case "break":
			fallthrough
		case "continue":

			return Returner{
				Type: v.Type,
			}

		case "return":

			value := Interpreter(ins, v.ExpAct, stacktrace, stacksize+1, nil, true).Exp

			defer dealloc(ins, varnames, value)

			return Returner{
				Type: "return",
				Exp:  value,
			}

		case "defer":

			actions := v.ExpAct
			defer func() { Interpreter(ins, actions, stacktrace, stacksize+1, nil, true) }()

		case "condition":

			for _, v := range v.ExpAct {

				truthy := true

				if v.Type == "if" {
					condition := Interpreter(ins, v.First, stacktrace, stacksize+1, nil, true)
					truthy = isTruthy(*condition.Exp)
				}

				if truthy {
					interpreted := Interpreter(ins, v.ExpAct, stacktrace, stacksize+1, nil, true)

					if interpreted.Type == "return" || interpreted.Type == "break" || interpreted.Type == "continue" {
						return Returner{
							Type: interpreted.Type,
							Exp:  interpreted.Exp,
						}
					}

					break
				}
			}

		case "while":

			cond := Interpreter(ins, v.First, stacktrace, stacksize+1, nil, true)

			for ; isTruthy(*cond.Exp); cond = Interpreter(ins, v.First, stacktrace, stacksize+1, nil, true) {

				interpreted := Interpreter(ins, v.ExpAct, stacktrace, stacksize+1, nil, true)

				if interpreted.Type == "return" {
					return Returner{
						Type: interpreted.Type,
						Exp:  interpreted.Exp,
					}
				}

				if interpreted.Type == "break" {
					break
				}
				if interpreted.Type == "continue" {
					continue
				}
			}

		case "each":

			it := *Interpreter(ins, []Action{v.First[0]}, stacktrace, stacksize+1, nil, true).Exp
			keyName := v.First[1].Name //get name of key
			valName := v.First[2].Name //get name of val

			it.Range(func(key, val *TuskType) Returner {

				ins.Allocate(keyName, key)
				ins.Allocate(valName, val)

				interpreted := Interpreter(ins, v.ExpAct, stacktrace, stacksize+1, nil, true)

				//free the key and val spaces
				ins.Deallocate(keyName)
				ins.Deallocate(valName)
				/////////////////////////////

				if interpreted.Type == "return" || interpreted.Type == "continue" || interpreted.Type == "break" {
					return Returner{
						Type: interpreted.Type,
						Exp:  interpreted.Exp,
					}
				}

				var undefval TuskType = undef

				return Returner{
					Type: "none",
					Exp:  &undefval,
				}
			})

		case "++":

			variable := Interpreter(ins, v.First, stacktrace, stacksize+1, nil, true)

			operationFunc, exists := Operations[(*variable.Exp).Type()+" + number"]

			if !exists { //if there is no operation for that type, panic
				TuskPanic("Could not find + operation for types "+(*variable.Exp).Type()+" and number", v.Line, v.File, stacktrace)
			}

			var onetype TuskType = one
			*variable.Exp = *operationFunc(*variable.Exp, onetype, ins, stacktrace, v.Line, v.File, stacksize)

			if expReturn {
				defer dealloc(ins, varnames, variable.Exp)
				return Returner{
					Type: "expression",
					Exp:  variable.Exp,
				}
			}

		case "--":

			variable := Interpreter(ins, v.First, stacktrace, stacksize+1, nil, true)

			operationFunc, exists := Operations[(*variable.Exp).Type()+" - number"]

			if !exists { //if there is no operation for that type, panic
				TuskPanic("Could not find - operation for types "+(*variable.Exp).Type()+" and number", v.Line, v.File, stacktrace)
			}

			var onetype TuskType = one
			*variable.Exp = *operationFunc(*variable.Exp, onetype, ins, stacktrace, v.Line, v.File, stacksize)

			if expReturn {
				defer dealloc(ins, varnames, variable.Exp)
				return Returner{
					Type: "expression",
					Exp:  variable.Exp,
				}
			}

		case "+=":
			fallthrough
		case "-=":
			fallthrough
		case "*=":
			fallthrough
		case "/=":
			fallthrough
		case "%=":
			fallthrough
		case "^=":

			variable := Interpreter(ins, v.First, stacktrace, stacksize+1, nil, true)
			interpreted := *Interpreter(ins, v.Second, stacktrace, stacksize+1, nil, true).Exp

			operationFunc, exists := Operations[(*variable.Exp).Type()+" "+string(v.Type[0])+" "+interpreted.Type()]

			if !exists { //if there is no operation for that type, panic
				TuskPanic("Could not find "+string(v.Type[0])+" operation for types "+(*variable.Exp).Type()+" and "+interpreted.Type(), v.Line, v.File, stacktrace)
			}

			*variable.Exp = *operationFunc(*variable.Exp, interpreted, ins, stacktrace, v.Line, v.File, stacksize)

			if expReturn {
				defer dealloc(ins, varnames, variable.Exp)
				return Returner{
					Type: "expression",
					Exp:  variable.Exp,
				}
			}

		}
	}

	var undefval TuskType = undef

	defer dealloc(ins, varnames, &undefval)

	return Returner{
		Type: "none",
		Exp:  &undefval,
	}
}
