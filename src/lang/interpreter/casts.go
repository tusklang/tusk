package interpreter

import . "lang/types"

func cast(val OmmType, nType string, stacktrace []string, line uint64, file string) *OmmType {

  if val.Type() == nType {
    return &val
  }

  switch val.Type() + "->" + nType {

    case "number->string":
      str := NumNormalize(val.(OmmNumber)) //convert to string
      var ommstr OmmString //create an ommstring
      ommstr.FromGoType(str)
      var ommtype OmmType = ommstr //create an ommtype interface
      return &ommtype

    case "string->number":
      integer, decimal := BigNumConverter(val.(OmmString).ToGoType())
      var newNum OmmType = OmmNumber{
        Integer: &integer,
        Decimal: &decimal,
      }
      return &newNum

  }

  ommPanic("Cannot cast a " + val.Type() + " into a " + nType, line, file, stacktrace)

  //here because it wont work without it
  var none OmmType = undef
  return &none
}
