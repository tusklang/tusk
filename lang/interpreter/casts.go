package interpreter

import (
	"fmt"
	"strconv"

	"github.com/tusklang/tusk/lang/types"
	"github.com/tusklang/tusk/native"
)

//Cast casts one tusk value to another type
func Cast(val types.TuskType, nType string, stacktrace []string, line uint64, file string) (*types.TuskType, *types.TuskError) {

	if val.Type() == nType {
		return &val, nil
	}

	switch nType + "->" + val.TypeOf() {

	case "string->int":
		s := strconv.FormatInt(val.(types.TuskInt).Int, 10)
		var tuskstr types.TuskString
		tuskstr.FromGoType(s)
		var tusktype types.TuskType = tuskstr
		return &tusktype, nil

	case "int->string":
		str := val.(types.TuskString).ToGoType()
		i, e := strconv.ParseInt(str, 10, 64)

		if e != nil {
			return nil, native.TuskPanic("Invalid integer literal: \""+str+"\"", line, file, stacktrace, native.ErrCodes["INVALIDLIT"])
		}
		var integertype types.TuskType = types.TuskInt{
			Int: i,
		}
		return &integertype, nil

	case "string->float":
		s := fmt.Sprint(val.(types.TuskFloat).Float)
		var tuskstr types.TuskString
		tuskstr.FromGoType(s)
		var tusktype types.TuskType = tuskstr
		return &tusktype, nil

	case "float->string":
		str := val.(types.TuskString).ToGoType()
		f, e := strconv.ParseFloat(str, 64)

		if e != nil {
			return nil, native.TuskPanic("Invalid float literal: \""+str+"\"", line, file, stacktrace, native.ErrCodes["INVALIDLIT"])
		}
		var floattype types.TuskType = types.TuskFloat{
			Float: f,
		}
		return &floattype, nil

	case "float->int":
		gofloat := float64(val.(types.TuskInt).Int)
		var tuskfloat types.TuskFloat
		tuskfloat.FromGoType(gofloat)
		var tusktype types.TuskType = tuskfloat
		return &tusktype, nil

	case "int->float":
		goint := int64(val.(types.TuskFloat).Float)
		var tuskint types.TuskInt
		tuskint.FromGoType(goint)
		var tusktype types.TuskType = tuskint
		return &tusktype, nil

	case "string->big":
		str := types.NumNormalize(val.(types.TuskNumber)) //convert to string
		var tuskstr types.TuskString                      //create an tuskstring
		tuskstr.FromGoType(str)
		var tusktype types.TuskType = tuskstr //create an tusktype interface
		return &tusktype, nil

	case "big->string":
		integer, decimal := types.BigNumConverter(val.(types.TuskString).ToGoType())
		var newNum types.TuskType = types.TuskNumber{
			Integer: &integer,
			Decimal: &decimal,
		}
		return &newNum, nil

	case "int->rune":
		var gonum = val.(types.TuskRune).ToGoType()
		var number types.TuskInt
		number.FromGoType(int64(gonum))
		var tusktype types.TuskType = number
		return &tusktype, nil

	case "rune->int":
		var gorune = rune(val.(types.TuskInt).ToGoType())
		var tuskrune types.TuskRune
		tuskrune.FromGoType(gorune)
		var tusktype types.TuskType = tuskrune
		return &tusktype, nil

	case "int->bool":
		var gobool = val.(types.TuskBool).ToGoType()
		if gobool {
			var tusktype types.TuskType = types.TuskInt{
				Int: 1,
			}
			return &tusktype, nil
		}

		var tusktype types.TuskType = types.TuskInt{
			Int: 0,
		}
		return &tusktype, nil

	case "string->rune":
		var runelist = val.(types.TuskRune).ToGoType()
		var tuskstr types.TuskString
		tuskstr.FromRuneList([]rune{runelist})
		var tusktype types.TuskType = tuskstr
		return &tusktype, nil

	case "array->string":
		var runelist = val.(types.TuskString).ToRuneList()
		var tuskrunelist types.TuskArray

		//convert each go rune in the array to a tusk rune
		for _, v := range runelist {
			var currune types.TuskRune
			currune.FromGoType(v)
			tuskrunelist.PushBack(currune)
		}

		var tusktype types.TuskType = tuskrunelist
		return &tusktype, nil

	}

	return nil, native.TuskPanic("Cannot cast a "+val.TypeOf()+" into a "+nType, line, file, stacktrace, native.ErrCodes["INVALIDCAST"])
}
