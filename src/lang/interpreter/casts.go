package interpreter

import . "lang/types"

func cast(val OmmType, nType string, stacktrace []string, line uint64, file string) *OmmType {

  if val.Type() == nType {
    return &val
  }

  switch nType + "->" + val.TypeOf() {

    case "string->number":
      str := NumNormalize(val.(OmmNumber)) //convert to string
      var ommstr OmmString //create an ommstring
      ommstr.FromGoType(str)
      var ommtype OmmType = ommstr //create an ommtype interface
      return &ommtype

    case "number->string":
      integer, decimal := BigNumConverter(val.(OmmString).ToGoType())
      var newNum OmmType = OmmNumber{
        Integer: &integer,
        Decimal: &decimal,
      }
      return &newNum

    case "number->rune":
      var gonum = float64(val.(OmmRune).ToGoType())
      var number OmmNumber
      number.FromGoType(gonum)
      var ommtype OmmType = number
      return &ommtype

    case "rune->number":
      var gorune = rune(val.(OmmNumber).ToGoType())
      var ommrune OmmRune
      ommrune.FromGoType(gorune)
      var ommtype OmmType = ommrune
      return &ommtype

    case "number->bool":
      var gobool = val.(OmmBool).ToGoType()
      if gobool {
        var ommtype OmmType = one
        return &ommtype
      } else {
        var ommtype OmmType = zero
        return &ommtype
      }

    case "string->rune":
      var runelist = val.(OmmRune).ToGoType()
      var ommstr OmmString
      ommstr.FromRuneList([]rune{ runelist })
      var ommtype OmmType = ommstr
      return &ommtype

    case "int->float":
      var goint = int64(val.(OmmFloat).Gofloat)
      var ommint OmmInteger
      ommint.Goint = goint
      var ommtype OmmType = ommint
      return &ommtype

    case "float->int":
      var gofloat = float64(val.(OmmInteger).Goint)
      var ommfloat OmmFloat
      ommfloat.Gofloat = gofloat
      var ommtype OmmType = ommfloat
      return &ommtype

    case "int->number":
      var goint = int64(val.(OmmNumber).ToGoType())
      var ommint OmmInteger
      ommint.Goint = goint
      var ommtype OmmType = ommint
      return &ommtype

    case "number->int":
      var gofloat = float64(val.(OmmInteger).Goint)
      var ommnum OmmNumber
      ommnum.FromGoType(gofloat)
      var ommtype OmmType = ommnum
      return &ommtype

    case "float->number":
      var gofloat = float64(val.(OmmNumber).ToGoType())
      var ommint OmmFloat
      ommint.Gofloat = gofloat
      var ommtype OmmType = ommint
      return &ommtype

    case "number->float":
      var gofloat = float64(val.(OmmFloat).Gofloat)
      var ommnum OmmNumber
      ommnum.FromGoType(gofloat)
      var ommtype OmmType = ommnum
      return &ommtype

  }

  OmmPanic("Cannot cast a " + val.Type() + " into a " + nType, line, file, stacktrace)

  //here because it wont work without it
  var none OmmType = undef
  return &none
}
