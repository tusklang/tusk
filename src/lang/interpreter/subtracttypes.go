package interpreter

func number__minus__number(num1, num2 Action, cli_params CliParams) Action {
  ensurePrec(&num1, &num2, cli_params)

  num2Placeholder := zero //create a placeholder for num2 (so it wont mutate)
  num2Placeholder.Integer, num2Placeholder.Decimal = make([]int64, len(num2.Integer)), make([]int64, len(num2.Decimal)) //allocate the length

  //looks like this
  // a - b = a + -b

  //invert the decimal
  for k, v := range num2.Decimal {
    num2Placeholder.Decimal[k] = -1 * v
  }
  //invert the integer
  for k, v := range num2.Integer {
    num2Placeholder.Integer[k] = -1 * v
  }

  return number__plus__number(num1, num2Placeholder, cli_params)
}
