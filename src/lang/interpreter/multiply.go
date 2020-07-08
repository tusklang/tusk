package interpreter

import . "lang/types"

func multiply(num1, num2 Action, cli_params CliParams) Action {

  /* TABLE OF TYPES:

    num * num = num
    string * num = string
    array * array = array
    default = falsey
  */

  type1 := num1.Type
  type2 := num2.Type

  var final Action

  if type1 == "number" && type2 == "number" {
    //detect case `num * num = num`
    final = number__times__number(num1, num2, cli_params)
  } else if (type1 == "string" && type2 == "number") || (type1 == "number" && type2 == "string") {
    //detect case `string * num = string`
    final = string__times__number(num1, num2, cli_params)
  } else if type1 == "array" && type2 == "array" {
    //detect case `array * array = array`
    final = array__times__array(num1, num2, cli_params)
  } else {
    //detect default case
    final = undef
  }

  return final
}
