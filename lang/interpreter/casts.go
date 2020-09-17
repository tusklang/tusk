package interpreter

import (
	"github.com/tusklang/tusk/lang/types"
)

//Cast casts one tusk value to another type
func Cast(val types.TuskType, nType string, stacktrace []string, line uint64, file string) *types.TuskType {

	if val.Type() == nType {
		return &val
	}

	switch nType + "->" + val.TypeOf() {

	case "string->number":
		str := types.NumNormalize(val.(types.TuskNumber)) //convert to string
		var kastr types.TuskString                       //create an kastring
		kastr.FromGoType(str)
		var tusktype types.TuskType = kastr //create an tusktype interface
		return &tusktype

	case "number->string":
		integer, decimal := types.BigNumConverter(val.(types.TuskString).ToGoType())
		var newNum types.TuskType = types.TuskNumber{
			Integer: &integer,
			Decimal: &decimal,
		}
		return &newNum

	case "number->rune":
		var gonum = float64(val.(types.TuskRune).ToGoType())
		var number types.TuskNumber
		number.FromGoType(gonum)
		var tusktype types.TuskType = number
		return &tusktype

	case "rune->number":
		var gorune = rune(val.(types.TuskNumber).ToGoType())
		var karune types.TuskRune
		karune.FromGoType(gorune)
		var tusktype types.TuskType = karune
		return &tusktype

	case "number->bool":
		var gobool = val.(types.TuskBool).ToGoType()
		if gobool {
			var tusktype types.TuskType = one
			return &tusktype
		}

		var tusktype types.TuskType = zero
		return &tusktype

	case "string->rune":
		var runelist = val.(types.TuskRune).ToGoType()
		var kastr types.TuskString
		kastr.FromRuneList([]rune{runelist})
		var tusktype types.TuskType = kastr
		return &tusktype

	}

	TuskPanic("Cannot cast a "+val.Type()+" into a "+nType, line, file, stacktrace)

	//here because it wont work without it
	var none types.TuskType = undef
	return &none
}
