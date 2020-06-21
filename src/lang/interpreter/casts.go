package interpreter

import . "lang"

func cast(val Action, type string) Action {
  valType := val.Type

  val.Type = type

  switch valType + "->" + type {
    case "string->number":
      val.Integer, val.Decimal = BigNumConverter(val.ExpStr)
    default:
      return val
  }
}
