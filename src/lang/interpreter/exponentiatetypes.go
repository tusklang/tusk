package interpreter

func number__pow__number(num1, num2 Action, cli_params CliParams) Action {
  ensurePrec(&num1, &num2, cli_params)

  //algorithm is exp(n1 ln n2)

  if len(num2.Decimal) == 0 { //if the exponent is an integer, use binary exponentiation for an O(log n) solution
    return number__pow__integer(num1, num2, cli_params)
  }

  return exp(number__times__number(num2, ln(num1, cli_params), cli_params), cli_params)
}
