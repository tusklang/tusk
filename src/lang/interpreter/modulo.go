package interpreter

import . "lang/types"

func number__mod__number(val1, val2 OmmType, cli_params CliParams) OmmType {
  num1, num2 := val1.(OmmNumber), val2.(OmmNumber)
  ensurePrec(&num1, &num2, cli_params)

  //ALGORITHM:
  //  num1 - floor(num1 / num2) * num2

  num2P := zero //create a placeholder for num2 (because it will get mutated)
  *num2P.Integer, *num2P.Decimal = append([]int64{}, *num2.Integer...), append([]int64{}, *num2.Decimal...)

  //if you set the prec to 0 here, it will mutate it
  divided := number__divide__number(num1, num2, cli_params).(OmmNumber)
  *divided.Decimal = nil //round down

  multiplied := number__times__number(divided, num2, cli_params).(OmmNumber)
  return number__minus__number(num1, multiplied, cli_params)
}
