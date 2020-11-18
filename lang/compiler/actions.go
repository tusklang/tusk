package compiler

import (
	"runtime"

	. "github.com/tusklang/tusk/lang/types"
)

func arraytogroup(arractions []Action) []Action {
	//convert a parenthesis array to a {} group
	//	if (true || false && !true) { ;would be registered as an array
	//		;blah blah
	//	}

	var converted []Action

	if len(arractions) == 0 {
		return arractions
	}

	switch arractions[0].Type {
	case "c-array":

		if arractions[0].Name == "definite-array" {
			return arractions
		}

		//range through it and append to the converted
		arractions[0].Value.(TuskArray).Range(func(_, v *TuskType) (Returner, *TuskError) {
			converted = append(converted, Action{
				Type:  (*v).Type(),
				Value: *v,
			})

			return Returner{}, nil
		})

	case "r-array":

		if arractions[0].Name == "definite-array" {
			return arractions
		}

		for _, v := range arractions[0].Array {
			converted = append(converted, v...)
		}

	default:
		converted = arractions
	}

	return converted
}

func actionizer(operations []Operation) ([]Action, error) {

	var actions []Action
	var e error

	for _, v := range operations {

		var left []Action
		var right []Action

		if v.Left != nil {
			left, e = actionizer([]Operation{*v.Left})

			if e != nil {
				return []Action{}, e
			}
		}
		if v.Right != nil {
			right, e = actionizer([]Operation{*v.Right})

			if e != nil {
				return []Action{}, e
			}
		}

		switch v.Type {
		case "STATE-OP":

			var statements = []string{"var", "if", "elif", "else", "while", "each", "function", "return", "proto", "static", "instance", "build", "ovld", "defer", "access", "try", "catch"} //list of statements

			var hasStatement bool = false

			for _, val := range statements {
				if val == (*v.Left).Item.Token.Name {

					switch val {
					case "function":
						if right[0].Type != "CB-OB" {
							return []Action{}, makeCompilerErr("Functions need a parameter list and a function body", v.File, right[0].Line)
						}

						var typeList []string
						var paramList []string
						var fnparams = arraytogroup(right[0].First)

						for _, p := range fnparams {
							if p.Type == "variable" {
								//automatically infer that it is type "any"
								typeList = append(typeList, "any")
								paramList = append(paramList, p.Name)
								continue
							}

							if p.Type != "cast" || p.ExpAct[0].Type != "variable" {
								return []Action{}, makeCompilerErr("Function parameter lists can only have variables", v.File, right[0].Line)
							}

							typeList = append(typeList, p.Name)
							paramList = append(paramList, p.ExpAct[0].Name)
						}

						actions = append(actions, Action{
							Type: "function",
							Value: TuskFunc{
								Overloads: []Overload{
									Overload{
										Params: paramList,
										Types:  typeList,
										Body:   arraytogroup(right[0].Second),
									},
								},
							},
							File: v.File,
							Line: v.Line,
						})

					case "if":

						if right[0].Type != "CB-OB" {
							return []Action{}, makeCompilerErr("If statements need a condition and a body", v.File, right[0].Line)
						}

						actions = append(actions, Action{
							Type: "condition",
							ExpAct: []Action{Action{
								Type:   "if",
								First:  arraytogroup(right[0].First),
								ExpAct: arraytogroup(right[0].Second),
							}},
							File: v.File,
							Line: v.Line,
						})

					case "elif":

						if right[0].Type != "CB-OB" {
							return []Action{}, makeCompilerErr("Elif statements need a condition and a body", v.File, right[0].Line)
						}

						if len(actions) == 0 || actions[len(actions)-1].Type != "condition" {
							return []Action{}, makeCompilerErr("Unexpected elif statement", v.File, right[0].Line)
						}

						//append to the previous conditional statement
						actions[len(actions)-1].ExpAct = append(actions[len(actions)-1].ExpAct, Action{
							Type:   "if",
							First:  arraytogroup(right[0].First),
							ExpAct: arraytogroup(right[0].Second),
						})

					case "else":

						if len(actions) == 0 || actions[len(actions)-1].Type != "condition" {
							return []Action{}, makeCompilerErr("Unexpected else statement", v.File, right[0].Line)
						}

						//append to the previous conditional statement
						actions[len(actions)-1].ExpAct = append(actions[len(actions)-1].ExpAct, Action{
							Type:   "else",
							ExpAct: arraytogroup(right),
						})

					case "try":

						actions = append(actions, Action{
							Type:   val,
							Second: arraytogroup(right),
							File:   v.File,
							Line:   v.Line,
						})

					case "catch":

						if right[0].Type != "CB-OB" {
							return []Action{}, makeCompilerErr("Catch statements need a catcher and a body", v.File, right[0].Line)
						}

						if len(actions) == 0 || actions[len(actions)-1].Type != "try" {
							return []Action{}, makeCompilerErr("Unexpected catch statement", v.File, right[0].Line)
						}

						var catchervars = right[0].First

						var tmpvar = Action{
							Type: "variable",
							Name: "dv 0", //dummy variable name
							File: v.File,
							Line: v.Line,
						}

						switch len(catchervars) {
						case 0: //both vars are dummies
							catchervars = make([]Action, 2)
							catchervars[0] = tmpvar
							catchervars[1] = tmpvar
						case 1: //only the second one is a dummy
							catchervars = append(catchervars, tmpvar)
						case 2: //none are dummies
						default: //error
							return []Action{}, makeCompilerErr("Catch statement catcher list was given too many parameters", v.File, right[0].Line)
						}

						//append to the previous try statement

						//prepend a dummy value (because in varname check, I am reusing the "each", and each requires an iterator for the first one)
						actions[len(actions)-1].First = append([]Action{Action{}}, catchervars...)
						actions[len(actions)-1].ExpAct = arraytogroup(right[0].Second)
						//////////////////////////////////////

					case "while":
						if right[0].Type != "CB-OB" {
							return []Action{}, makeCompilerErr("While loops need a condition and a body", v.File, right[0].Line)
						}

						actions = append(actions, Action{
							Type:   val,
							First:  arraytogroup(right[0].First),
							ExpAct: arraytogroup(right[0].Second),
							File:   v.File,
							Line:   v.Line,
						})

					case "each":
						if right[0].Type != "CB-OB" {
							return []Action{}, makeCompilerErr("Each loops require an iterator and a body", v.File, right[0].Line)
						}

						iter := arraytogroup(right[0].First)

						if len(iter) != 3 {
							return []Action{}, makeCompilerErr("Each loops must have an iterator, key, and a value", v.File, right[0].Line)
						}

						for _, n := range iter[1:] {
							if n.Type != "variable" {
								return []Action{}, makeCompilerErr("Key or value was not given as a variable", v.File, right[0].Line)
							}
						}

						actions = append(actions, Action{
							Type:   val,
							First:  iter, //because it doesnt matter if they use a { or (
							ExpAct: arraytogroup(right[0].Second),
							File:   v.File,
							Line:   v.Line,
						})

					case "var":

						if right[0].Type == "variable" { //the dev is declaring is like "var a" (meaning declare a)
							actions = append(actions, Action{
								Type: "declare",
								Name: right[0].Name,
								File: v.File,
								Line: v.Line,
							})
						} else {
							if right[0].Type != "let" {
								return []Action{}, makeCompilerErr("Expected a assigner statement after var", v.File, right[0].Line)
							}

							if right[0].First[0].Type != "variable" {
								return []Action{}, makeCompilerErr("Cannot use :: operator in variable declaration", v.File, right[0].Line)
							}

							actions = append(actions, Action{
								Type:   val,
								Name:   arraytogroup(right[0].First)[0].Name,
								ExpAct: arraytogroup(right[0].ExpAct),
								File:   v.File,
								Line:   v.Line,
							})
						}
					case "proto":

						//prototype compilation is **pretty** messy

						if len(right) == 0 {
							return []Action{}, makeCompilerErr("Prototypes require a body", v.File, right[0].Line)
						}

						if right[0].Type != "{" {
							return []Action{}, makeCompilerErr("Prototype bodies can only be curly brace enclosed", v.File, right[0].Line)
						}

						var (
							static   = make(map[string]*TuskType)
							instance = make(map[string]*TuskType)
						)
						var body = right[0].ExpAct //get the struct body

						for i := range body {
							if body[i].Type != "static" && body[i].Type != "instance" { //if it does not name static or instance, automatically make it instance
								body[i] = Action{
									Type:   "instance",
									ExpAct: []Action{body[i]},
									File:   body[i].File,
									Line:   body[i].Line,
								}
							}

							name := arraytogroup(body[i].ExpAct)[0].Name

							if len(name) == 0 {
								return nil, makeCompilerErr("Prototype body variable has no name", v.File, right[0].Line)
							}

							if body[i].ExpAct[0].Type == "var" {

								if len(arraytogroup(body[i].ExpAct)[0].ExpAct) == 0 || arraytogroup(arraytogroup(body[i].ExpAct)[0].ExpAct)[0].Value == nil {
									return []Action{}, makeCompilerErr("Cannot have compound types at the global scope of a prototype", v.File, right[0].Line)
								}

								var current = body[i].ExpAct[0].ExpAct[0].Value

								if body[i].Type == "static" {
									static[name] = &current
								} else {
									instance[name] = &current
								}

							} else if body[i].ExpAct[0].Type == "declare" {

								var tmp TuskType = TuskUndef{}

								if body[i].Type == "static" {
									static[name] = &tmp
								} else {
									instance[name] = &tmp
								}
							} else if body[i].ExpAct[0].Type == "ovld" {

								if body[i].Type == "static" {

									if _, e := static[name]; !e || (*static[name]).Type() != "function" { //if the value does not exist yet (or it is not a function), make it
										var tmp TuskType = TuskFunc{}
										static[name] = &tmp
									}

									var currentfn = (*static[name]).(TuskFunc)
									currentfn.Overloads = append(currentfn.Overloads, body[i].ExpAct[0].ExpAct[0].Value.(TuskFunc).Overloads[0])
									*static[name] = currentfn
								} else {

									if _, e := instance[name]; !e || (*instance[name]).Type() != "function" { //if the value does not exist yet (or it is not a function), make it
										var tmp TuskType = TuskFunc{}
										instance[name] = &tmp
									}

									var currentfn = (*instance[name]).(TuskFunc)
									currentfn.Overloads = append(currentfn.Overloads, body[i].ExpAct[0].ExpAct[0].Value.(TuskFunc).Overloads[0])
									*instance[name] = currentfn
								}

							} else {
								return []Action{}, makeCompilerErr("Prototype bodies can only have variable assignments, declarations, and overloading", v.File, right[0].Line)
							}

						}

						actions = append(actions, Action{
							Type: "prototype",
							Value: TuskProto{
								Static:   static,
								Instance: instance,
							},
							File: v.File,
							Line: v.Line,
						})

					case "build":

						if right[0].First[0].Type != "c-array" {
							return []Action{}, makeCompilerErr("Expected an array after build", v.File, right[0].Line)
						}

						var e error
						var dobuild bool

						right[0].First[0].Value.(TuskArray).Range(func(kk *TuskType, vv *TuskType) (none Returner, err *TuskError) {

							none = Returner{
								Type: "break",
							}

							if (*vv).Type() != "string" {
								e = makeCompilerErr("All build conditions must be strings", v.File, right[0].Line)
								return
							}

							condition := (*vv).(TuskString).ToGoType()

							if len(condition) != 0 && condition[0] == '!' { //if it detects not that type of os
								condition = condition[1:]
								if runtime.GOOS != condition {
									dobuild = false
									return
								}
							}

							if runtime.GOOS == condition {
								dobuild = true
								return
							}

							return
						})

						if e != nil {
							return nil, e
						}

						if dobuild {
							var block = right[0].Second
							actions = append(actions, block...)
						}

					case "ovld":

						if right[0].Type != "let" {
							return []Action{}, makeCompilerErr("Expected a assigner statement after ovld", v.File, right[0].Line)
						}

						if right[0].First[0].Type != "variable" {
							return []Action{}, makeCompilerErr("Cannot use :: operator in an overloader", v.File, right[0].Line)
						}

						if right[0].ExpAct[0].Type != "function" {
							return []Action{}, makeCompilerErr("Cannot overload a "+right[0].ExpAct[0].Type, v.File, right[0].Line)
						}

						right[0].Type = "ovld"

						actions = append(actions, Action{
							Type:   val,
							Name:   arraytogroup(right[0].First)[0].Name,
							ExpAct: arraytogroup(right[0].ExpAct),
							File:   v.File,
							Line:   v.Line,
						})

					case "access":

						if right[0].Value.Type() != "array" {
							return nil, makeCompilerErr("Must use an array for an access list", v.File, v.Line)
						}

						actions = append(actions, Action{
							Type:  "access",
							Value: right[0].Value,
							File:  v.File,
							Line:  v.Line,
						})

					default:

						actions = append(actions, Action{
							Type:   val,
							ExpAct: arraytogroup(right),
							File:   v.File,
							Line:   v.Line,
						})

					}

					hasStatement = true
				}
			}

			if !hasStatement {
				return []Action{}, makeCompilerErr("\""+(*v.Left).Item.Token.Name+"\" is not a statement", v.File, v.Line)
			}

		case ":=":

			if left[0].Type != "variable" {
				return []Action{}, makeCompilerErr("Expected a variable statement before := operator", v.File, right[0].Line)
			}

			actions = append(actions, Action{
				Type:   "var",
				Name:   left[0].Name,
				ExpAct: arraytogroup(right),
				File:   v.File,
				Line:   v.Line,
			})

		case "=":

			actions = append(actions, Action{
				Type:   "let",
				First:  left,
				ExpAct: right,
				File:   v.File,
				Line:   v.Line,
			})

		case "->":

			castType := v.Left.Item.Token.Name[1:]

			actions = append(actions, Action{
				Type:   "cast",
				Name:   castType, //type to cast into
				ExpAct: arraytogroup(right),
				File:   v.File,
				Line:   v.Line,
			})

		//all of these operations have the same way of appending
		case "::":

			//if it is ::, and the next action is a variable, then convert to a string
			//to get index of a variable's value, use ::()
			//for example,
			//  var a: [
			//    "hello" = "world",
			//  ]
			//  log a::hello ; would log "world"
			//
			//  var idx: "hello"
			//  log a::(idx) ; would log "world" as well
			//  log a::idx ; //would cause a panic error

			if len(right) == 0 { //safeguard
				return []Action{}, makeCompilerErr("No value was found right of a :: operator", v.File, v.Line)
			}

			if right[0].Type == "variable" {
				var str = TuskString{}
				str.FromGoType(right[0].Name)
				right[0] = Action{
					Type:  "string",
					Value: str,
					File:  right[0].File,
					Line:  right[0].Line,
				}
			} else if right[0].Type == "r-array" || right[0].Type == "c-array" { //detect ::()
				right = arraytogroup(right)
			}

			fallthrough
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
		case "&&":
			fallthrough
		case "||":
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
		case "CB-OB":
			fallthrough
		case ":":
			fallthrough
		case "?":

			rightn := right

			if v.Type != ":" && v.Type != "?" { //not a function call
				rightn = arraytogroup(right)
			}

			actions = append(actions, Action{
				Type:   v.Type,
				First:  arraytogroup(left),
				Second: rightn,
				File:   v.File,
				Line:   v.Line,
			})
		////////////////////////////////////////////////////////

		case "++":
			fallthrough
		case "--":

			if len(left) == 0 || (left[0].Type != "variable" && left[0].Type != "::") {
				return []Action{}, makeCompilerErr("Must have a variable before an increment or decrement", v.File, v.Line)
			}

			actions = append(actions, Action{
				Type:  v.Type,
				First: arraytogroup(left),
				File:  v.File,
				Line:  v.Line,
			})

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

			if len(left) == 0 || (left[0].Type != "variable" && left[0].Type != "::") {
				return []Action{}, makeCompilerErr("Must have a variable before an assignment operator", v.File, v.Line)
			}
			if len(right) == 0 {
				return []Action{}, makeCompilerErr("Could not find a value after "+v.Type, v.File, v.Line)
			}

			actions = append(actions, Action{
				Type:   v.Type,
				First:  arraytogroup(left),
				Second: arraytogroup(right),
				File:   v.File,
				Line:   v.Line,
			})

		case "break":
			fallthrough
		case "continue":

			actions = append(actions, Action{
				Type: v.Type,
				File: v.File,
				Line: v.Line,
			})

		case "none":
			vActs, e := valueActions(v.Item)

			if e != nil {
				return []Action{}, e
			}

			actions = append(actions, vActs)
		}
	}

	return actions, nil
}
