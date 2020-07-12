package interpreter

import . "lang/types"

func number__minus__number(val1, val2 OmmType, cli_params CliParams, line uint64, file string) *OmmType {
  num1, num2 := val1.(OmmNumber), val2.(OmmNumber)
  ensurePrec(&num1, &num2, cli_params)

  num2Placeholder := zero //create a placeholder for num2 (so it wont mutate)
  tmpInt, tmpDec := make([]int64, len(*num2.Integer)), make([]int64, len(*num2.Decimal)) //allocate the length
  num2Placeholder.Integer, num2Placeholder.Decimal = &tmpInt, &tmpDec

  //looks like this
  // a - b = a + -b

  //invert the decimal
  for k, v := range *num2.Decimal {
    (*num2Placeholder.Decimal)[k] = -1 * v
  }

  //invert the integer
  for k, v := range *num2.Integer {
    (*num2Placeholder.Integer)[k] = -1 * v
  }

  return number__plus__number(num1, num2Placeholder, cli_params, line, file)
}
