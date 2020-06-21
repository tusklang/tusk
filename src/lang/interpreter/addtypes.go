package interpreter

//string + (* - array - hash) = string
func string__plus__all_not_array_not_hash(num1, num2 Action, cli_params CliParams) Action {

  num2C := cast(num2, "string")

  var val = num1
  val.ExpStr+=num2C.ExpStr

  return val
}
