package interpreter

import (
	"github.com/omm-lang/omm/lang/types"
	"github.com/omm-lang/omm/ommstd/native"
)

//Cast casts one omm value to another type
func Cast(val types.OmmType, nType string, stacktrace []string, line uint64, file string) *types.OmmType {

	if val.Type() == nType {
		return &val
	}

	switch nType + "->" + val.TypeOf() {

	case "string->number":
		str := types.NumNormalize(val.(types.OmmNumber)) //convert to string
		var ommstr types.OmmString                       //create an ommstring
		ommstr.FromGoType(str)
		var ommtype types.OmmType = ommstr //create an ommtype interface
		return &ommtype

	case "number->string":
		integer, decimal := types.BigNumConverter(val.(types.OmmString).ToGoType())
		var newNum types.OmmType = types.OmmNumber{
			Integer: &integer,
			Decimal: &decimal,
		}
		return &newNum

	case "number->rune":
		var gonum = float64(val.(types.OmmRune).ToGoType())
		var number types.OmmNumber
		number.FromGoType(gonum)
		var ommtype types.OmmType = number
		return &ommtype

	case "rune->number":
		var gorune = rune(val.(types.OmmNumber).ToGoType())
		var ommrune types.OmmRune
		ommrune.FromGoType(gorune)
		var ommtype types.OmmType = ommrune
		return &ommtype

	case "number->bool":
		var gobool = val.(types.OmmBool).ToGoType()
		if gobool {
			var ommtype types.OmmType = one
			return &ommtype
		}

		var ommtype types.OmmType = zero
		return &ommtype

	case "string->rune":
		var runelist = val.(types.OmmRune).ToGoType()
		var ommstr types.OmmString
		ommstr.FromRuneList([]rune{runelist})
		var ommtype types.OmmType = ommstr
		return &ommtype

	case "int->float":
		var goint = int64(val.(types.OmmFloat).Gofloat)
		var ommint types.OmmInteger
		ommint.Goint = goint
		var ommtype types.OmmType = ommint
		return &ommtype

	case "float->int":
		var gofloat = float64(val.(types.OmmInteger).Goint)
		var ommfloat types.OmmFloat
		ommfloat.Gofloat = gofloat
		var ommtype types.OmmType = ommfloat
		return &ommtype

	case "int->number":
		var goint = int64(val.(types.OmmNumber).ToGoType())
		var ommint types.OmmInteger
		ommint.Goint = goint
		var ommtype types.OmmType = ommint
		return &ommtype

	case "number->int":
		var gofloat = float64(val.(types.OmmInteger).Goint)
		var ommnum types.OmmNumber
		ommnum.FromGoType(gofloat)
		var ommtype types.OmmType = ommnum
		return &ommtype

	case "float->number":
		var gofloat = float64(val.(types.OmmNumber).ToGoType())
		var ommint types.OmmFloat
		ommint.Gofloat = gofloat
		var ommtype types.OmmType = ommint
		return &ommtype

	case "number->float":
		var gofloat = float64(val.(types.OmmFloat).Gofloat)
		var ommnum types.OmmNumber
		ommnum.FromGoType(gofloat)
		var ommtype types.OmmType = ommnum
		return &ommtype

	}

	native.OmmPanic("Cannot cast a "+val.Type()+" into a "+nType, line, file, stacktrace)

	//here because it wont work without it
	var none types.OmmType = undef
	return &none
}
