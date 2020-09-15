package compiler

import (
	"runtime"

	. "ka/lang/types"
)

func arraytogroup(arractions []Action) []Action {
	//convert a parenthesis array to a {} group
	//	if (true || false && !true) { ;would be registered as an array
	//		;blah blah
	//	}

	var converted []Action

	switch arractions[0].Type {
	case "c-array":

		//range through it and append to the converted
		arractions[0].Value.(KaArray).Range(func(_, v *KaType) Returner {
			converted = append(converted, Action{
				Type:  (*v).Type(),
				Value: *v,
			})

			return Returner{}
		})

	case "r-array":

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
		case "~":

			var statements = []string{"var", "if", "elif", "else", "while", "each", "include", "function", "return", "proto", "static", "instance", "ifwin", "ifnwin", "ovld", "defer", "access"} //list of statements

			var hasStatement bool = false

			for _, val := range statements {
				if val == (*v.Left).Item.Token.Name {

					switch val {
					case "function":
						if right[0].Type != "=>" {
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
							Value: KaFunc{
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

						if right[0].Type != "=>" {
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

						if right[0].Type != "=>" {
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

					case "while":
						if right[0].Type != "=>" {
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
						if right[0].Type != "=>" {
							return []Action{}, makeCompilerErr("Each loops need a iterator and a body", v.File, right[0].Line)
						}

						iter := arraytogroup(right[0].First)

						if len(iter) != 3 {
							return []Action{}, makeCompilerErr("Each loops must look like this: each(iterator, key, value)", v.File, right[0].Line)
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
								ExpAct: right[0].ExpAct,
								File:   v.File,
								Line:   v.Line,
							})
						}

					case "proto":

						if len(right) == 0 {
							return []Action{}, makeCompilerErr("Prototypes require a body", v.File, right[0].Line)
						}

						if right[0].Type != "{" {
							return []Action{}, makeCompilerErr("Prototype bodies can only be curly brace enclosed", v.File, right[0].Line)
						}

						var (
							static   = make(map[string]*KaType)
							instance = make(map[string]*KaType)
							access   = make(map[string][]string)
						)
						var body = right[0].ExpAct //get the struct body
						var currentaccess []string

						for i := range body {

							if body[i].Type == "access" { //protected field
								if currentaccess != nil {
									return nil, makeCompilerErr("Access found twice in a row", body[i].File, body[i].Line)
								}

								for _, v := range body[i].Value.(KaArray).Array {
									if (*v).Type() != "string" {
										return nil, makeCompilerErr("Expect a string list to access", body[i].File, body[i].Line)
									}

									var cur = (*v).(KaString).ToGoType()

									if cur == "thisf" { //"thisf" means this file
										cur = body[i].File
									}

									currentaccess = append(currentaccess, cur)
								}

								continue
							}

							if body[i].Type != "static" && body[i].Type != "instance" { //if it does not name static or instance, automatically make it instance
								body[i] = Action{
									Type:   "instance",
									ExpAct: []Action{body[i]},
									File:   body[i].File,
									Line:   body[i].Line,
								}
							}

							name := arraytogroup(body[i].ExpAct)[0].Name[1:]

							if currentaccess != nil {
								access[name] = currentaccess
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

								var tmp KaType = KaUndef{}

								if body[i].Type == "static" {
									static[name] = &tmp
								} else {
									instance[name] = &tmp
								}
							} else if body[i].ExpAct[0].Type == "ovld" {

								if body[i].Type == "static" {

									if _, e := static[name]; !e || (*static[name]).Type() != "function" { //if the value does not exist yet (or it is not a function), make it
										var tmp KaType = KaFunc{}
										static[name] = &tmp
									}

									var currentfn = (*static[name]).(KaFunc)
									currentfn.Overloads = append(currentfn.Overloads, body[i].ExpAct[0].ExpAct[0].Value.(KaFunc).Overloads[0])
									*static[name] = currentfn
								} else {

									if _, e := instance[name]; !e || (*instance[name]).Type() != "function" { //if the value does not exist yet (or it is not a function), make it
										var tmp KaType = KaFunc{}
										instance[name] = &tmp
									}

									var currentfn = (*instance[name]).(KaFunc)
									currentfn.Overloads = append(currentfn.Overloads, body[i].ExpAct[0].ExpAct[0].Value.(KaFunc).Overloads[0])
									*instance[name] = currentfn
								}

							} else {
								return []Action{}, makeCompilerErr("Prototype bodies can only have variable assignments, declarations, and overloading", v.File, right[0].Line)
							}

							currentaccess = nil
						}

						actions = append(actions, Action{
							Type: "proto",
							Value: KaProto{
								Static:     static,
								Instance:   instance,
								AccessList: access,
							},
							File: v.File,
							Line: v.Line,
						})

					case "ifwin":

						if runtime.GOOS == "windows" {
							actions = append(actions, right...)
						}

					case "ifnwin":

						if runtime.GOOS != "windows" {
							actions = append(actions, right...)
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
							ExpAct: right,
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
				ExpAct: right,
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
				ExpAct: right,
				File:   v.File,
				Line:   v.Line,
			})

		//all of these operations have the same way of appending
		case "::":

			//if it is ::, and the next action is a variable, then convert to a string
			//to get index of a variable's value, use ::()
			//for example,
			//  var a: [:
			//    "hello": "world",
			//  :]
			//  log a::hello ; would log "world"
			//
			//  var idx: "hello"
			//  log a::(idx) ; would log "world" as well
			//  log a::idx ; //would cause a panic error

			if len(right) == 0 { //safeguard
				return []Action{}, makeCompilerErr("No value was found right of a :: operator", v.File, v.Line)
			}

			if right[0].Type == "variable" {
				var str = KaString{}
				str.FromGoType(right[0].Name[1:]) //remove the $ from the varname
				right[0] = Action{
					Type:  "string",
					Value: str,
					File:  right[0].File,
					Line:  right[0].Line,
				}
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
		case "&":
			fallthrough
		case "|":
			fallthrough
		case "=>":
			fallthrough
		case ":":
			fallthrough
		case "?":

			actions = append(actions, Action{
				Type:   v.Type,
				First:  left,
				Second: right,
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
				First: left,
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
		case "%=":
			fallthrough
		case "^=":

			if len(left) == 0 || (left[0].Type != "variable" && left[0].Type != "::") {
				return []Action{}, makeCompilerErr("Must have a variable before an assignment operator", v.File, v.Line)
			}
			if len(right) == 0 {
				return []Action{}, makeCompilerErr("Could not find a value after "+v.Type, v.File, v.Line)
			}

			actions = append(actions, Action{
				Type:   v.Type,
				First:  left,
				Second: right,
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
