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
      var gostring = string(val.(OmmRune).ToGoType())
      var ommstr OmmString
      ommstr.FromGoType(gostring)
      var ommtype OmmType = ommstr
      return &ommtype

  }

  OmmPanic("Cannot cast a " + val.Type() + " into a " + nType, line, file, stacktrace)

  //here because it wont work without it
  var none OmmType = undef
  return &none
}
