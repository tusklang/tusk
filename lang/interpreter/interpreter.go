package interpreter

import (
	"strconv"

	. "github.com/tusklang/tusk/lang/types"
	"github.com/tusklang/tusk/native"
	. "github.com/tusklang/tusk/native"
)

const MAX_STACKSIZE = 100001

var tmpret uint64 = 0 //return values are stored in temporary new variables, and this is the counter for those

func bindReturn(ins *Instance, ret *TuskType, bindto string) {

	switch (*ret).TypeOf() {
	case "function":

		for _, v := range (*ret).(TuskFunc).Overloads {
			for _, vv := range v.VarRefs {
				ins.Bind(vv, bindto)
			}
		}

	case "array":

		(*ret).(TuskArray).Range(func(k, v *TuskType) (Returner, *TuskError) {
			bindReturn(ins, v, bindto)
			return Returner{}, nil
		})

	case "hash":

		(*ret).(TuskHash).Range(func(k, v *TuskType) (Returner, *TuskError) {
			bindReturn(ins, v, bindto)
			return Returner{}, nil
		})

	}

}

//Interpreter starts the Tusk runtime with a given instance, and actions tree
func Interpreter(ins *Instance, actions []Action, stacktrace []string, stacksize uint, expReturn bool, namespace string) (Returner, *TuskError) {

	var allocated []string

	if stacksize > MAX_STACKSIZE {
		TuskPanic("Stack size was exceeded", 0, "none", stacktrace, native.ErrCodes["NOMEM"])
	}

	for _, v := range actions {
		switch v.Type {

		case "var":

			interpreted, e := Interpreter(ins, v.ExpAct, stacktrace, stacksize+1, true, namespace)

			if e != nil {
				return Returner{}, e
			}

			name := v.Name
			ins.Allocate(name, interpreted.Exp)
			allocated = append(allocated, name)

			if expReturn {
				variable := interpreted.Exp
				return Returner{
					Type: "expression",
					Exp:  variable,
				}, nil
			}

		case "ovld":

			interpreted, e := Interpreter(ins, v.ExpAct, stacktrace, stacksize+1, true, namespace)

			if e != nil {
				return Returner{}, e
			}

			if (*(*ins.Fetch(v.Name)).Value).Type() != "function" {
				*(*ins.Fetch(v.Name)).Value = TuskFunc{}
			}

			var appended_ovld TuskType = TuskFunc{
				Overloads: append((*(*ins.Fetch(v.Name)).Value).(TuskFunc).Overloads, (*interpreted.Exp).(TuskFunc).Overloads[0]),
			}

			ins.Allocate(v.Name, &appended_ovld)

			if expReturn {
				variable := &appended_ovld
				return Returner{
					Type: "expression",
					Exp:  variable,
				}, nil
			}

		case "declare":

			var tmpundef TuskType = undef

			name := v.Name
			ins.Allocate(name, &tmpundef)
			allocated = append(allocated, name)

			if expReturn {
				variable := ins.Fetch(v.Name).Value
				return Returner{
					Type: "expression",
					Exp:  variable,
				}, nil
			}

		case "let":

			pinterpreted, e := Interpreter(ins, v.ExpAct, stacktrace, stacksize+1, true, namespace)

			if e != nil {
				return Returner{}, e
			}

			interpreted := *pinterpreted.Exp

			variable, e := Interpreter(ins, v.First, stacktrace, stacksize+1, true, namespace)

			if e != nil {
				return Returner{}, e
			}

			*variable.Exp = interpreted

			if expReturn {
				return Returner{
					Type: "expression",
					Exp:  variable.Exp,
				}, nil
			}

		//all of the types
		case "string":
			fallthrough
		case "rune":
			fallthrough
		case "int":
			fallthrough
		case "float":
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

			val := v.Value

			if expReturn {
				if cloned := val.Clone(); cloned != nil { //if it can be cloned
					return Returner{
						Type: "expression",
						Exp:  cloned, //clone it and return
					}, nil
				}

				//otherwise, dont clone it
				return Returner{
					Type: "expression",
					Exp:  &val,
				}, nil
			}

		case "function":
			//for a function, add the instance
			var nf = v.Value.(TuskFunc)
			nf.Instance = ins
			var tusktype TuskType = nf
			if expReturn {
				return Returner{
					Type: "expression",
					Exp:  &tusktype,
				}, nil
			}

		//arrays, hashes are a bit different
		case "r-array":

			var nArr = make([]*TuskType, len(v.Array))

			for k, i := range v.Array {
				tmp, e := Interpreter(ins, i, stacktrace, stacksize+1, true, namespace)
				if e != nil {
					return Returner{}, e
				}
				nArr[k] = tmp.Exp
			}

			var kaType TuskType = TuskArray{
				Array:  nArr,
				Length: uint64(len(v.Array)),
			}

			if expReturn {
				return Returner{
					Type: "expression",
					Exp:  &kaType,
				}, nil
			}

		case "r-hash":

			var nHash = make(map[string]*TuskType)

			for _, i := range v.Hash {

				keyi, e := Interpreter(ins, i[0], stacktrace, stacksize+1, true, namespace)

				if e != nil {
					return Returner{}, e
				}

				vali, e := Interpreter(ins, i[1], stacktrace, stacksize+1, true, namespace)

				if e != nil {
					return Returner{}, e
				}

				nHash[(*keyi.Exp).Format()] = vali.Exp
			}

			var kaType TuskType = TuskHash{
				Hash:   nHash,
				Length: uint64(len(v.Hash)),
			}

			if expReturn {
				return Returner{
					Type: "expression",
					Exp:  &kaType,
				}, nil
			}

		////////////////////////////////////

		case "variable":

			_fetched := ins.Fetch(v.Name)

			if _fetched == nil {
				//if it is a nil pointer
				return Returner{}, TuskPanic("Invalid memory address referenced", v.Line, v.File, stacktrace, native.ErrCodes["NILPTR"])
			}

			fetched := _fetched.Value
			if expReturn {
				return Returner{
					Type: "expression",
					Exp:  fetched,
				}, nil
			}

		case "{":

			groupRet, e := Interpreter(ins, v.ExpAct, stacktrace, stacksize+1, false, namespace)

			if e != nil {
				return Returner{}, e
			}

			if expReturn {
				return Returner{
					Type: groupRet.Type,
					Exp:  groupRet.Exp,
				}, nil
			}

		case "cast":

			cv, e := Interpreter(ins, v.ExpAct, stacktrace, stacksize+1, true, namespace)
			if e != nil {
				return Returner{}, e
			}
			casted, e := Cast(*cv.Exp, v.Name, stacktrace, v.Line, v.File)

			if e != nil {
				return Returner{}, e
			}

			if expReturn {
				return Returner{
					Type: "expression",
					Exp:  casted,
				}, nil
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
		case "//":
			fallthrough
		case "%":
			fallthrough
		case "**":
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
		case "||":
			fallthrough
		case "&&":
			fallthrough
		case "~":
			fallthrough
		case "&":
			fallthrough
		case "|":
			fallthrough
		case "^":
			fallthrough
		case ">>":
			fallthrough
		case "<<":
			fallthrough
		case ":":
			fallthrough
		case "?":

			//this is kinda spaghetti ngl
			//but it works for now

			var firstInterpreted Returner
			var secondInterpreted Returner

			var computed *TuskType
			var retaddr string
			var e *TuskError

			//& and | can be a bit different for bool and bool
			if v.Type == "&&" || v.Type == "||" {
				firstInterpreted, e = Interpreter(ins, v.First, stacktrace, stacksize+1, true, namespace)

				if e != nil {
					return Returner{}, e
				}

				var computed TuskType

				if (*firstInterpreted.Exp).Type() == "bool" {

					var assumeVal bool

					if v.Type == "&&" {
						assumeVal = !isTruthy(*firstInterpreted.Exp)
					} else {
						assumeVal = isTruthy(*firstInterpreted.Exp)
					}

					if assumeVal {
						var retVal = v.Type == "||"
						computed = TuskBool{
							Boolean: &retVal,
						}
					} else {
						secondInterpreted, e = Interpreter(ins, v.Second, stacktrace, stacksize+1, true, namespace)

						if e != nil {
							return Returner{}, e
						}

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
						return Returner{
							Type: "expression",
							Exp:  &computed,
						}, nil
					}
				}
			}
			//////////////////////////////////////////////////'

			firstInterpreted, e = Interpreter(ins, v.First, stacktrace, stacksize+1, true, namespace)
			if e != nil {
				return Returner{}, e
			}
			secondInterpreted, e = Interpreter(ins, v.Second, stacktrace, stacksize+1, true, namespace)
			if e != nil {
				return Returner{}, e
			}

			if v.Type == "==" || v.Type == "!=" {
				comparedVal := ((*firstInterpreted.Exp).TypeOf() == (*secondInterpreted.Exp).TypeOf()) == (v.Type == "==")

				if !comparedVal {
					var tmp TuskType = falsev
					computed = &tmp
					goto comparsion_false
				}
			}

		and_or_skip:
			{
				//check typeof first
				operationFunc, exists := Operations[(*firstInterpreted.Exp).TypeOf()+" "+v.Type+" "+(*secondInterpreted.Exp).TypeOf()]

				if !exists { //if it does not exist, also check the type
					operationFunc, exists = Operations[(*firstInterpreted.Exp).Type()+" "+v.Type+" "+(*secondInterpreted.Exp).Type()]
				}

				if !exists { //if there is no operation for that type, panic
					return Returner{}, TuskPanic("Could not find "+v.Type+" operator for types "+(*firstInterpreted.Exp).TypeOf()+" and "+(*secondInterpreted.Exp).TypeOf(), v.Line, v.File, stacktrace, native.ErrCodes["OPNOTFOUND"])
				}

				computed, e, retaddr = operationFunc(*firstInterpreted.Exp, *secondInterpreted.Exp, ins, stacktrace, v.Line, v.File, stacksize+1, namespace)

				if retaddr != "" {
					allocated = append(allocated, retaddr)
				}

				if e != nil {
					return Returner{}, e
				}
			}

		comparsion_false:

			if expReturn {
				return Returner{
					Type: "expression",
					Exp:  computed,
				}, nil
			}

		////////////

		case "break":
			fallthrough
		case "continue":

			return Returner{
				Type: v.Type,
			}, nil

		case "return":

			pvalue, e := Interpreter(ins, v.ExpAct, stacktrace, stacksize+1, true, namespace)

			if e != nil {
				return Returner{}, e
			}

			value := pvalue.Exp

			tmpName := "ret" + strconv.FormatUint(tmpret, 10)
			tmpret++

			bindReturn(ins, value, tmpName)
			ins.Allocate(tmpName, value)

			for _, v := range allocated { //remove the scope binds of all the local variables
				ins.DeBind(v, "scope")
			}

			return Returner{
				Type:       "return",
				Exp:        value,
				ReturnAddr: tmpName,
			}, nil

		case "defer":

			actions := v.ExpAct
			defer Interpreter(ins, actions, stacktrace, stacksize+1, true, namespace) //call this defer before the garbage collection, because the garbage collection ones are FIFO, but regular is LIFO, so this will be executed first

		case "condition":

			for _, v := range v.ExpAct {

				truthy := true

				if v.Type == "if" {
					condition, e := Interpreter(ins, v.First, stacktrace, stacksize+1, true, namespace)
					if e != nil {
						return Returner{}, e
					}
					truthy = isTruthy(*condition.Exp)
				}

				if truthy {
					interpreted, e := Interpreter(ins, v.ExpAct, stacktrace, stacksize+1, true, namespace)

					if e != nil {
						return Returner{}, e
					}

					if interpreted.Type == "return" || interpreted.Type == "break" || interpreted.Type == "continue" {
						return Returner{
							Type: interpreted.Type,
							Exp:  interpreted.Exp,
						}, nil
					}

					break
				}
			}

		case "try":

			_, e := Interpreter(ins, v.Second, stacktrace, stacksize+1, false, namespace)

			if e != nil {
				evar := v.First[1].Name
				svar := v.First[2].Name

				//allocate the error and stacktrace variables
				{
					var estr TuskString
					estr.FromGoType(e.Err)
					var tusktype TuskType = estr
					ins.Allocate(evar, &tusktype)
				}

				{
					var sarray TuskArray

					for _, v := range e.Stacktrace {
						var curstack TuskString
						curstack.FromGoType(v)
						sarray.PushBack(curstack)
					}

					var tusktype TuskType = sarray
					ins.Allocate(svar, &tusktype)
				}
				/////////////////////////////////////////////

				//call the catch block
				Interpreter(ins, v.ExpAct, stacktrace, stacksize+1, false, namespace)
			}

		case "while":

			cond, e := Interpreter(ins, v.First, stacktrace, stacksize+1, true, namespace)

			if e != nil {
				return Returner{}, e
			}

			for ; isTruthy(*cond.Exp); cond, e = Interpreter(ins, v.First, stacktrace, stacksize+1, true, namespace) {

				if e != nil {
					return Returner{}, e
				}

				interpreted, e := Interpreter(ins, v.ExpAct, stacktrace, stacksize+1, true, namespace)

				if e != nil {
					return Returner{}, e
				}

				if interpreted.Type == "return" {
					return Returner{
						Type: interpreted.Type,
						Exp:  interpreted.Exp,
					}, nil
				}

				if interpreted.Type == "break" {
					break
				}
				if interpreted.Type == "continue" {
					continue
				}
			}

		case "each":

			pit, e := Interpreter(ins, []Action{v.First[0]}, stacktrace, stacksize+1, true, namespace)
			if e != nil {
				return Returner{}, e
			}
			it := *pit.Exp
			keyName := v.First[1].Name //get name of key
			valName := v.First[2].Name //get name of val

			v, e := it.Range(func(key, val *TuskType) (Returner, *TuskError) {

				ins.Allocate(keyName, key)
				ins.Allocate(valName, val)

				interpreted, e := Interpreter(ins, v.ExpAct, stacktrace, stacksize+1, true, namespace)

				if e != nil {
					return Returner{}, e
				}

				Interpreter(ins, []Action{v.First[0]}, stacktrace, stacksize+1, true, namespace)

				if interpreted.Type == "return" || interpreted.Type == "continue" || interpreted.Type == "break" {
					return Returner{
						Type: interpreted.Type,
						Exp:  interpreted.Exp,
					}, nil
				}

				var undefval TuskType = undef

				return Returner{
					Type: "none",
					Exp:  &undefval,
				}, nil
			})

			if e != nil {
				return Returner{}, e
			}

			if v != nil && (expReturn || ((*v).Type == "return" || (*v).Type == "break" || (*v).Type == "continue")) {
				return *v, nil
			}

		case "++":

			variable, e := Interpreter(ins, v.First, stacktrace, stacksize+1, true, namespace)

			if e != nil {
				return Returner{}, e
			}

			operationFunc, exists := Operations[(*variable.Exp).Type()+" + int"]

			if !exists { //if there is no operation for that type, panic
				return Returner{}, TuskPanic("Could not find + operation for types "+(*variable.Exp).Type()+" and number", v.Line, v.File, stacktrace, native.ErrCodes["OPNOTFOUND"])
			}

			var onetype TuskInt
			onetype.FromGoType(1)
			tmp, e, _ := operationFunc(*variable.Exp, onetype, ins, stacktrace, v.Line, v.File, stacksize, namespace)
			if e != nil {
				return Returner{}, e
			}
			*variable.Exp = *tmp

			if expReturn {
				return Returner{
					Type: "expression",
					Exp:  variable.Exp,
				}, nil
			}

		case "--":

			variable, e := Interpreter(ins, v.First, stacktrace, stacksize+1, true, namespace)

			if e != nil {
				return Returner{}, e
			}

			operationFunc, exists := Operations[(*variable.Exp).Type()+" - int"]

			if !exists { //if there is no operation for that type, panic
				return Returner{}, TuskPanic("Could not find - operation for types "+(*variable.Exp).Type()+" and number", v.Line, v.File, stacktrace, native.ErrCodes["OPNOTFOUND"])
			}

			var onetype TuskInt
			onetype.FromGoType(1)
			tmp, e, _ := operationFunc(*variable.Exp, onetype, ins, stacktrace, v.Line, v.File, stacksize, namespace)
			if e != nil {
				return Returner{}, e
			}
			*variable.Exp = *tmp

			if expReturn {
				return Returner{
					Type: "expression",
					Exp:  variable.Exp,
				}, nil
			}

		case "+=":
			fallthrough
		case "-=":
			fallthrough
		case "*=":
			fallthrough
		case "/=":
			fallthrough
		case "//=":
			fallthrough
		case "%=":
			fallthrough
		case "**=":

			variable, e := Interpreter(ins, v.First, stacktrace, stacksize+1, true, namespace)
			if e != nil {
				return Returner{}, e
			}
			pinterpreted, e := Interpreter(ins, v.Second, stacktrace, stacksize+1, true, namespace)
			if e != nil {
				return Returner{}, e
			}
			interpreted := *pinterpreted.Exp

			operation := v.Type[:len(v.Type)-1]
			operationFunc, exists := Operations[(*variable.Exp).Type()+" "+operation+" "+interpreted.Type()]

			if !exists { //if there is no operation for that type, panic
				return Returner{}, TuskPanic("Could not find "+operation+" operation for types "+(*variable.Exp).Type()+" and "+interpreted.Type(), v.Line, v.File, stacktrace, native.ErrCodes["OPNOTFOUND"])
			}

			calc, e, _ := operationFunc(*variable.Exp, interpreted, ins, stacktrace, v.Line, v.File, stacksize, namespace)
			if e != nil {
				return Returner{}, e
			}
			*variable.Exp = *calc

			if expReturn {
				return Returner{
					Type: "expression",
					Exp:  variable.Exp,
				}, nil
			}

		}
	}

	var undefval TuskType = undef

	return Returner{
		Type: "none",
		Exp:  &undefval,
	}, nil
}
