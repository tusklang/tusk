package interpreter

func number__minus__number(num1, num2 Action, cli_params CliParams) Action {

  //looks like this
  // a - b = a + -b

  //invert the decimal
  for k, v := range num2.Decimal {
    num2.Decimal[k] = -1 * v
  }
  //invert the integer
  for k, v := range num2.Integer {
    num2.Integer[k] = -1 * v
  }

  return number__plus__number(num1, num2, cli_params)
}
