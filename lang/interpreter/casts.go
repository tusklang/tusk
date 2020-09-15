package interpreter

import (
	"ka/lang/types"
)

//Cast casts one ka value to another type
func Cast(val types.KaType, nType string, stacktrace []string, line uint64, file string) *types.KaType {

	if val.Type() == nType {
		return &val
	}

	switch nType + "->" + val.TypeOf() {

	case "string->number":
		str := types.NumNormalize(val.(types.KaNumber)) //convert to string
		var kastr types.KaString                       //create an kastring
		kastr.FromGoType(str)
		var katype types.KaType = kastr //create an katype interface
		return &katype

	case "number->string":
		integer, decimal := types.BigNumConverter(val.(types.KaString).ToGoType())
		var newNum types.KaType = types.KaNumber{
			Integer: &integer,
			Decimal: &decimal,
		}
		return &newNum

	case "number->rune":
		var gonum = float64(val.(types.KaRune).ToGoType())
		var number types.KaNumber
		number.FromGoType(gonum)
		var katype types.KaType = number
		return &katype

	case "rune->number":
		var gorune = rune(val.(types.KaNumber).ToGoType())
		var karune types.KaRune
		karune.FromGoType(gorune)
		var katype types.KaType = karune
		return &katype

	case "number->bool":
		var gobool = val.(types.KaBool).ToGoType()
		if gobool {
			var katype types.KaType = one
			return &katype
		}

		var katype types.KaType = zero
		return &katype

	case "string->rune":
		var runelist = val.(types.KaRune).ToGoType()
		var kastr types.KaString
		kastr.FromRuneList([]rune{runelist})
		var katype types.KaType = kastr
		return &katype

	}

	KaPanic("Cannot cast a "+val.Type()+" into a "+nType, line, file, stacktrace)

	//here because it wont work without it
	var none types.KaType = undef
	return &none
}
