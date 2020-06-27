package interpreter

func exponentiate(num1, num2 Action, cli_params CliParams) Action {

  /* TABLE OF TYPES:

    num ^ num = num
    default = falsey
  */

  type1 := num1.Type
  type2 := num2.Type

  var final Action

  if type1 == "number" && type2 == "number" {
    final = number__pow__number(num1, num2, cli_params)
  } else {
    final = undef
  }

  return final
}
