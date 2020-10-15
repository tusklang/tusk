package compiler

import (
	"unicode"

	. "github.com/tusklang/tusk/lang/types"
)

func valueActions(item Item) (Action, error) {

	switch item.Type {

	case "{":

		oper, e := makeOperations(item.Group)

		if e != nil {
			return Action{}, e
		}

		acts, e := actionizer(oper)

		if e != nil {
			return Action{}, e
		}

		return Action{
			Type:   "{",
			ExpAct: acts,
			File:   item.File,
			Line:   item.Line,
		}, nil
	case "(":

		var arr [][]Action
		var carr TuskArray

		var arrtype = "c-array" //compile time array

		for _, v := range item.Group {
			_oper, e := makeOperations([][]Item{v})

			if e != nil {
				return Action{}, e
			}

			oper := _oper[0]

			value, e := actionizer([]Operation{oper})

			if e != nil {
				return Action{}, e
			}

			if len(value) == 0 {
				return Action{}, makeCompilerErr("Each entry in the array must have a value", item.File, item.Line)
			}

			if value[0].Type == "proto" {
				return Action{}, makeCompilerErr("Cannot have protos outside of the global scope", item.File, item.Line)
			}

			if value[0].Value == nil || value[0].Type == "function" {
				arr = append(arr, value)
				arrtype = "r-array" ///make it a runtime array
				continue
			}

			arr = append(arr, value)
			carr.PushBack(value[0].Value)
		}

		var definitearray = item.Name

		if arrtype == "c-array" {
			return Action{
				Type:  arrtype,
				Name:  definitearray,
				Value: carr,
				File:  item.File,
				Line:  item.Line,
			}, nil
		}

		return Action{
			Type:  arrtype,
			Name:  definitearray,
			Array: arr,
			File:  item.File,
			Line:  item.Line,
		}, nil
	case "[:":

		var hash = make([][2][]Action, 0)
		var chash TuskHash

		var hashtype = "c-hash" //compile time hash

		for _, v := range item.Group {
			_oper, e := makeOperations([][]Item{v})

			if e != nil {
				return Action{}, e
			}

			oper := _oper[0]

			//give errors
			if oper.Type != "=" {
				return Action{}, makeCompilerErr("Expected a '=' for a hash key", item.File, oper.Line)
			}
			/////////////

			key, e := actionizer([]Operation{*oper.Left})

			if e != nil {
				return Action{}, e
			}

			value, e := actionizer([]Operation{*oper.Right})

			if e != nil {
				return Action{}, e
			}

			if len(value) == 0 {
				return Action{}, makeCompilerErr("Expected some value as after ':'", item.File, oper.Line)
			}

			if value[0].Type == "proto" {
				return Action{}, makeCompilerErr("Cannot have protos outside of the global scope", item.File, item.Line)
			}

			if key[0].Value == nil || value[0].Value == nil {
				hashtype = "r-hash" ///make it a runtime hash
				goto runthash
			} else {
				chash.Set(&key[0].Value, value[0].Value)
			}

		runthash:
			hash = append(hash, [2][]Action{
				key,
				value,
			})
		}

		if hashtype == "c-hash" {
			return Action{
				Type:  hashtype,
				Value: chash,
				File:  item.File,
				Line:  item.Line,
			}, nil
		}

		return Action{
			Type: hashtype,
			Hash: hash,
			File: item.File,
			Line: item.Line,
		}, nil

	case "expression value":

		var val = item.Token.Name

		if len(val) == 0 {
			return Action{}, makeCompilerErr("Value is not valid", item.File, item.Line)
		}

		if val[0] == '"' || val[0] == '`' { //detect string
			var str = TuskString{}
			str.FromGoType(val[1 : len(val)-1])
			return Action{
				Type:  "string",
				Value: str,
				File:  item.File,
				Line:  item.Line,
			}, nil
		} else if val[0] == '\'' { //detect a rune
			var oRune = TuskRune{}

			qrem := val[1 : len(val)-1] //remove quotes

			if len(qrem) != 1 {
				return Action{}, makeCompilerErr("Runes must be one character long", item.File, item.Line)
			}

			oRune.FromGoType([]rune(qrem)[0])
			return Action{
				Type:  "rune",
				Value: oRune,
				File:  item.File,
				Line:  item.Line,
			}, nil
		} else if val == "true" || val == "false" { //detect a bool
			var boolean = TuskBool{}
			boolean.FromGoType(val == "true" /* convert to a boolean */)
			return Action{
				Type:  "bool",
				Value: boolean,
				File:  item.File,
				Line:  item.Line,
			}, nil
		} else if val == "undef" { //detect an undef value
			var undef TuskUndef
			return Action{
				Type:  "undef",
				Value: undef,
				File:  item.File,
				Line:  item.Line,
			}, nil
		} else if unicode.IsDigit(rune(val[0])) || val[0] == '.' || val[0] == '+' || val[0] == '-' { //detect a number
			var number = TuskNumber{}
			number.FromString(val)
			return Action{
				Type:  "number",
				Value: number,
				File:  item.File,
				Line:  item.Line,
			}, nil
		} else if val[0] == '$' { //detect a variable
			return Action{
				Type: "variable",
				Name: val[1:],
				File: item.File,
				Line: item.Line,
			}, nil
		} else { //detect nothing, which throws an error
			return Action{}, makeCompilerErr(val+" is not a value", item.File, item.Line)
		}

	}

	return Action{}, nil
}
