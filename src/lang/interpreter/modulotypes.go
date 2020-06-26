package interpreter

func number__mod__number(num1, num2 Action, cli_params CliParams) Action {

  //ALGORITHM:
  //  num1 - floor(num1 / num2) * num2

  cli_params["Calc"]["PREC"] = 0 //force to round down
  divided := number__divide__number(num1, num2, cli_params)
  divided.Decimal = nil //just in case

  multiplied := number__times__number(divided, num2, cli_params)
  return number__minus__number(num1, multiplied, cli_params)
}
