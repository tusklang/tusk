package interpreter

import . "lang/types"

func add(num1, num2 Action, cli_params CliParams) Action {

  /* TABLE OF TYPES:

    string + (* - array - hash) = string
    array + * = array
    num + num = num
    hash + hash = hash
    boolean + boolean = boolean
    default = falsey
  */

  type1 := num1.Type
  type2 := num2.Type

  var final Action

  if (type1 == "string" && (type2 != "array" && type2 != "hash")) || (type2 == "string" && (type1 != "array" && type2 != "hash")) {
    //detect case `string + (* - array - hash) = string`
    final = string__plus__all_not_array_not_hash(num1, num2, cli_params)
  } else if type1 == "array" || type2 == "array" {
    //detect case `array + * = array`
    final = array__plus__all(num1, num2, cli_params)
  } else if type1 == "number" && type2 == "number" {
    //detect case `num + num = num`
    final = number__plus__number(num1, num2, cli_params)
  } else if type1 == "hash" && type2 == "hash" {
    //detect case `hash + hash = hash`
    final = hash__plus__hash(num1, num2, cli_params)
  } else if type1 == "boolean" && type2 == "boolean" {
    //detect case `boolean + boolean = boolean`
    final = bool__plus__bool(num1, num2, cli_params)
  } else {
    //detect default case
    final = undef
  }

  return final
}
