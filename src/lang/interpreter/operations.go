package interpreter

import . "lang/types"

//list of operations
var operations = map[string]func(val1, val2 OmmType, cli_params CliParams, line uint64, file string) *OmmType {
	"number + number": number__plus__number,
	"number - number": number__minus__number,
	"number * number": number__times__number,
	"number / number": number__divide__number,
	"number % number": number__mod__number,
	"number ^ number": number__pow__number,
	"number = number": func(val1, val2 OmmType, cli_params CliParams, line uint64, file string) *OmmType {

		var final = falsev

		if isEqual(val1.(OmmNumber), val2.(OmmNumber)) {
			final = truev
		}

		var finalType OmmType = final

		return &finalType
	},
	"number != number": func(val1, val2 OmmType, cli_params CliParams, line uint64, file string) *OmmType {

		var final = truev

		if isEqual(val1.(OmmNumber), val2.(OmmNumber)) {
			final = falsev
		}

		var finalType OmmType = final

		return &finalType
	},
	"string = string": func(val1, val2 OmmType, cli_params CliParams, line uint64, file string) *OmmType {

		var isEqual OmmType = falsev

		if val1.(OmmString).ToGoType() == val2.(OmmString).ToGoType() {
			isEqual = truev
		}

		return &isEqual
	},
	"string != string": func(val1, val2 OmmType, cli_params CliParams, line uint64, file string) *OmmType {

		var isEqual OmmType = truev

		if val1.(OmmString).ToGoType() == val2.(OmmString).ToGoType() {
			isEqual = falsev
		}

		return &isEqual
	},
	"bool = bool": func(val1, val2 OmmType, cli_params CliParams, line uint64, file string) *OmmType {

		var isEqual OmmType = falsev

		if val1.(OmmBool).ToGoType() == val2.(OmmBool).ToGoType() {
			isEqual = truev
		}

		return &isEqual
	},
	"bool != bool": func(val1, val2 OmmType, cli_params CliParams, line uint64, file string) *OmmType {

		var isEqual OmmType = truev

		if val1.(OmmBool).ToGoType() == val2.(OmmBool).ToGoType() {
			isEqual = falsev
		}

		return &isEqual
	},
	"rune = rune": func(val1, val2 OmmType, cli_params CliParams, line uint64, file string) *OmmType {

		var isEqual OmmType = falsev

		if val1.(OmmBool).ToGoType() == val2.(OmmBool).ToGoType() {
			isEqual = truev
		}

		return &isEqual
	},
	"rune != rune": func(val1, val2 OmmType, cli_params CliParams, line uint64, file string) *OmmType {

		var isEqual OmmType = truev

		if val1.(OmmBool).ToGoType() == val2.(OmmBool).ToGoType() {
			isEqual = falsev
		}

		return &isEqual
	},
	"undef ! bool": func(val1, val2 OmmType, cli_params CliParams, line uint64, file string) *OmmType {

		boolean := !val2.(OmmBool).ToGoType()
		var converted OmmType = OmmBool{
			Boolean: &boolean,
		}

		return &converted
	},
}
