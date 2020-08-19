package interpreter

import (
	"fmt"
	"os"

	. "github.com/omm-lang/omm/lang/types"
)

const MAX_STACKSIZE = 100001

//export OmmPanic
func OmmPanic(err string, line uint64, file string, stacktrace []string) {
	fmt.Println("Panic on line", line, "file", file)
	fmt.Println(err)
	fmt.Println("\nWhen the error was thrown, this was the stack:")
	fmt.Println("  at line", line, "in file", file)
	for i := len(stacktrace) - 1; i >= 0; i-- { //print the stacktrace
		fmt.Println("  " + stacktrace[i])
	}
	os.Exit(1)
}

func dealloc(ins *Instance, varnames []string, value *OmmType) { //function to remove the variables declared in that scope
	for _, v := range varnames {
		if (*value).Type() == "function" { //if it is a curryed function, make sure it does not garbage collect the vars used

			for _, vv := range (*value).(OmmFunc).Overloads {
				for _, vvv := range vv.VarRefs {
					if vvv == v {
						goto nodealloc
					}
				}
			}

		}

		ins.Deallocate(v)
	nodealloc:
	}
}

func Interpreter(ins *Instance, actions []Action, stacktrace []string, stacksize uint) Returner {

	if stacksize > MAX_STACKSIZE {
		OmmPanic("Stack size was exceeded", 0, "none", stacktrace)
	}

	var varnames []string //deallocate these names

	var expReturn = false //if it is inside an expression

	if len(actions) == 1 {
		expReturn = true
	}

	for _, v := range actions {
		switch v.Type {

		case "var":

			interpreted := Interpreter(ins, v.ExpAct, stacktrace, stacksize+1)

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

			interpreted := Interpreter(ins, v.ExpAct, stacktrace, stacksize+1)

			if (*(*ins.Fetch(v.Name)).Value).Type() != "function" {
				*(*ins.Fetch(v.Name)).Value = OmmFunc{}
			}

			var appended_ovld OmmType = OmmFunc{
				Overloads: append((*(*ins.Fetch(v.Name)).Value).(OmmFunc).Overloads, (*interpreted.Exp).(OmmFunc).Overloads[0]),
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

			var tmpundef OmmType = undef

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

			interpreted := *Interpreter(ins, v.ExpAct, stacktrace, stacksize+1).Exp

			variable := Interpreter(ins, v.First, stacktrace, stacksize+1)

			*variable.Exp = interpreted

			if expReturn {
				defer dealloc(ins, varnames, variable.Exp)
				return Returner{
					Type: "expression",
					Exp:  variable.Exp,
				}
			}

		case "log":
			interpreted := Interpreter(ins, v.ExpAct, stacktrace, stacksize+1)
			fmt.Println((*interpreted.Exp).Format())
		case "print":
			interpreted := Interpreter(ins, v.ExpAct, stacktrace, stacksize+1)
			fmt.Print((*interpreted.Exp).Format())

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
		case "proto":
			fallthrough
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
			var nf = v.Value.(OmmFunc)
			nf.Instance = ins
			var ommtype OmmType = nf
			if expReturn {
				defer dealloc(ins, varnames, &ommtype)
				return Returner{
					Type: "expression",
					Exp:  &ommtype,
				}
			}

		//arrays, hashes are a bit different
		case "r-array":

			var nArr = make([]*OmmType, len(v.Array))

			for k, i := range v.Array {
				nArr[k] = Interpreter(ins, i, stacktrace, stacksize+1).Exp
			}

			var ommType OmmType = OmmArray{
				Array:  nArr,
				Length: uint64(len(v.Array)),
			}

			if expReturn {
				defer dealloc(ins, varnames, &ommType)
				return Returner{
					Type: "expression",
					Exp:  &ommType,
				}
			}

		case "r-hash":

			var nHash = make(map[string]*OmmType)

			for _, i := range v.Hash {
				nHash[(*Interpreter(ins, i[0], stacktrace, stacksize+1).Exp).Format()] = Interpreter(ins, i[1], stacktrace, stacksize+1).Exp
			}

			var ommType OmmType = OmmHash{
				Hash:   nHash,
				Length: uint64(len(v.Hash)),
			}

			if expReturn {
				defer dealloc(ins, varnames, &ommType)
				return Returner{
					Type: "expression",
					Exp:  &ommType,
				}
			}

		////////////////////////////////////

		case "variable":

			if expReturn {
				fetched := ins.Fetch(v.Name).Value
				defer dealloc(ins, varnames, fetched)
				return Returner{
					Type: "expression",
					Exp:  fetched,
				}
			}

		case "{":
			fallthrough
		case "(":

			groupRet := Interpreter(ins, v.ExpAct, stacktrace, stacksize+1)

			if expReturn {
				defer dealloc(ins, varnames, groupRet.Exp)
				return Returner{
					Type: "expression",
					Exp:  groupRet.Exp,
				}
			}

		case "cast":

			casted := cast(*Interpreter(ins, v.ExpAct, stacktrace, stacksize+1).Exp, v.Name, stacktrace, v.Line, v.File)

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
		case "=>":
			fallthrough //this is probably not necessary, but i just left it here
		case "<-":
			fallthrough
		case "<~":

			var firstInterpreted Returner
			var secondInterpreted Returner

			//& and | can be a bit different for bool and bool
			if v.Type == "&" || v.Type == "|" {
				firstInterpreted = Interpreter(ins, v.First, stacktrace, stacksize+1)
				var computed OmmType

				if (*firstInterpreted.Exp).Type() == "bool" {

					var assumeVal bool

					if v.Type == "&" {
						assumeVal = !isTruthy(*firstInterpreted.Exp)
					} else {
						assumeVal = isTruthy(*firstInterpreted.Exp)
					}

					if assumeVal {
						var retVal = v.Type == "|"
						computed = OmmBool{
							Boolean: &retVal,
						}
					} else {
						secondInterpreted = Interpreter(ins, v.Second, stacktrace, stacksize+1)

						if (*secondInterpreted.Exp).Type() == "bool" {
							var gobool = isTruthy(*secondInterpreted.Exp)
							computed = OmmBool{
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
			//////////////////////////////////////////////////

			firstInterpreted = Interpreter(ins, v.First, stacktrace, stacksize+1)
			secondInterpreted = Interpreter(ins, v.Second, stacktrace, stacksize+1)

		and_or_skip:
			//check typeof first
			operationFunc, exists := Operations[(*firstInterpreted.Exp).TypeOf()+" "+v.Type+" "+(*secondInterpreted.Exp).TypeOf()]

			if !exists { //if it does not exist, also check the type
				operationFunc, exists = Operations[(*firstInterpreted.Exp).Type()+" "+v.Type+" "+(*secondInterpreted.Exp).Type()]
			}

			if !exists { //if there is no operation for that type, panic
				OmmPanic("Could not find "+v.Type+" operator for types "+(*firstInterpreted.Exp).TypeOf()+" and "+(*secondInterpreted.Exp).TypeOf(), v.Line, v.File, stacktrace)
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

		case "await":

			interpreted := Interpreter(ins, v.ExpAct, stacktrace, stacksize+1).Exp
			var awaited OmmType

			switch (*interpreted).(type) {
			case OmmThread:

				//put the new value back into the given interpreted pointer
				thread := (*interpreted).(OmmThread)
				thread.Join()
				*interpreted = thread
				///////////////////////////////////////////////////////////

				awaited = thread
			default:
				awaited = *interpreted
			}

			if expReturn {
				defer dealloc(ins, varnames, &awaited)
				return Returner{
					Type: "expression",
					Exp:  &awaited,
				}
			}

		case "break":
			fallthrough
		case "continue":

			return Returner{
				Type: v.Type,
			}

		case "return":

			value := Interpreter(ins, v.ExpAct, stacktrace, stacksize+1).Exp

			defer dealloc(ins, varnames, value)

			return Returner{
				Type: "return",
				Exp:  value,
			}

		case "defer":

			actions := v.ExpAct
			defer func() { Interpreter(ins, actions, stacktrace, stacksize+1) }()

		case "condition":

			for _, v := range v.ExpAct {

				truthy := true

				if v.Type == "if" {
					condition := Interpreter(ins, v.First, stacktrace, stacksize+1)
					truthy = isTruthy(*condition.Exp)
				}

				if truthy {
					interpreted := Interpreter(ins, v.ExpAct, stacktrace, stacksize+1)

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

			cond := Interpreter(ins, v.First, stacktrace, stacksize+1)

			for ; isTruthy(*cond.Exp); cond = Interpreter(ins, v.First, stacktrace, stacksize+1) {

				interpreted := Interpreter(ins, v.ExpAct, stacktrace, stacksize+1)

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

			it := *Interpreter(ins, []Action{v.First[0]}, stacktrace, stacksize+1).Exp
			keyName := v.First[1].Name //get name of key
			valName := v.First[2].Name //get name of val

			switch it.(type) {
			case OmmHash:

				for key, val := range it.(OmmHash).Hash {

					ommtypeKeyn := OmmString{}
					ommtypeKeyn.FromGoType(key)

					var ommtypeKey OmmType = ommtypeKeyn

					ins.Allocate(keyName, &ommtypeKey)
					ins.Allocate(valName, val)

					interpreted := Interpreter(ins, v.ExpAct, stacktrace, stacksize+1)

					//free the key and val spaces
					ins.Deallocate(keyName)
					ins.Deallocate(valName)
					/////////////////////////////

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

			case OmmArray:

				for key, val := range it.(OmmArray).Array {

					ommtypeKeyn := OmmNumber{}
					ommtypeKeyn.FromGoType(float64(key))

					var ommtypeKey OmmType = ommtypeKeyn

					ins.Allocate(keyName, &ommtypeKey)
					ins.Allocate(valName, val)

					interpreted := Interpreter(ins, v.ExpAct, stacktrace, stacksize+1)

					//free the key and val spaces
					ins.Deallocate(keyName)
					ins.Deallocate(valName)
					/////////////////////////////

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

			case OmmString:

				for key, val := range it.(OmmString).ToGoType() {

					ommtypeKeyn := OmmNumber{}
					ommtypeKeyn.FromGoType(float64(key))
					ommtypeValr := OmmRune{}
					ommtypeValr.FromGoType(val)

					var ommtypeKey OmmType = ommtypeKeyn
					var ommtypeVal OmmType = ommtypeValr

					ins.Allocate(keyName, &ommtypeKey)
					ins.Allocate(valName, &ommtypeVal)

					interpreted := Interpreter(ins, v.ExpAct, stacktrace, stacksize+1)

					//free the key and val spaces
					ins.Deallocate(keyName)
					ins.Deallocate(valName)
					/////////////////////////////

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

			}

		case "++":

			variable := Interpreter(ins, v.First, stacktrace, stacksize+1)

			operationFunc, exists := Operations[(*variable.Exp).Type()+" + number"]

			if !exists { //if there is no operation for that type, panic
				OmmPanic("Could not find + operation for types "+(*variable.Exp).Type()+" and number", v.Line, v.File, stacktrace)
			}

			var onetype OmmType = one
			*variable.Exp = *operationFunc(*variable.Exp, onetype, ins, stacktrace, v.Line, v.File, stacksize+1)

			if expReturn {
				defer dealloc(ins, varnames, variable.Exp)
				return Returner{
					Type: "expression",
					Exp:  variable.Exp,
				}
			}

		case "--":

			variable := Interpreter(ins, v.First, stacktrace, stacksize+1)

			operationFunc, exists := Operations[(*variable.Exp).Type()+" - number"]

			if !exists { //if there is no operation for that type, panic
				OmmPanic("Could not find - operation for types "+(*variable.Exp).Type()+" and number", v.Line, v.File, stacktrace)
			}

			var onetype OmmType = one
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

			variable := Interpreter(ins, v.First, stacktrace, stacksize+1)
			interpreted := *Interpreter(ins, v.Second, stacktrace, stacksize+1).Exp

			operationFunc, exists := Operations[(*variable.Exp).Type()+" "+string(v.Type[0])+" "+interpreted.Type()]

			if !exists { //if there is no operation for that type, panic
				OmmPanic("Could not find "+string(v.Type[0])+" operation for types "+(*variable.Exp).Type()+" and "+interpreted.Type(), v.Line, v.File, stacktrace)
			}

			*variable.Exp = *operationFunc(*variable.Exp, interpreted, ins, stacktrace, v.Line, v.File, stacksize+1)

			if expReturn {
				defer dealloc(ins, varnames, variable.Exp)
				return Returner{
					Type: "expression",
					Exp:  variable.Exp,
				}
			}

		}
	}

	var undefval OmmType = undef

	defer dealloc(ins, varnames, &undefval)

	return Returner{
		Type: "none",
		Exp:  &undefval,
	}
}
