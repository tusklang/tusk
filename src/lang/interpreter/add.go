package interpreter

func add(num1, num2 Action, cli_params CliParams) Action {

  /* TABLE OF TYPES:

    string + (* - array - hash) = string
    array + * = array
    num + num = num
    hash + hash = hash
    boolean + boolean = boolean
    num + boolean = num
    default = falsey
  */

  type1 := num1.Type
  type2 := num2.Type

  var final Action

  if (type1 == "string" && (type2 != "array" && type2 != "hash")) || (type2 == "string" && (type1 != "array" && type2 != "hash")) {
    //detect case `string + (* - array - hash) = string`
    final = string__plus__all_not_array_not_hash(num1, num2, cli_params)
  } else if type1 == "array" || type2 == "array" {
    final = array__plus__array(num1, num2, cli_params)
  }

  return final
}
