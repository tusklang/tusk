package interpreter

func cast(val Action, nType string) Action {
  valType := val.Type

  val.Type = nType

  switch valType + "->" + nType {
    case "string->number":
      val.Integer, val.Decimal = BigNumConverter(val.ExpStr)
    case "number->string":
      val.ExpStr = num_normalize(val)
  }

  return val
}
