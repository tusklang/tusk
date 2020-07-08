package interpreter

import . "lang/types"

func number__pow__number(num1, num2 Action, cli_params CliParams) Action {
  ensurePrec(&num1, &num2, cli_params)

  expNeg := false
  if isLess(num2, zero) { //account for negative exponents
    expNeg = true
    num2 = number__times__number(num2, neg_one, cli_params)
  }

  if len(num2.Decimal) == 0 { //if the exponent is an integer, use binary exponentiation for an O(log n) solution

    powwed := number__pow__integer(num1, num2, cli_params)

    if expNeg {
      powwed = number__divide__number(one, powwed, cli_params)
    }

    return powwed
  }

  var two = zero
  two.Integer = []int64{ 2 }

  neg := false
  if isEqual(number__mod__number(num1, two, cli_params), zero) { //because ln (n < 0) is undefined
    neg = true
  }

  num1 = abs(num1, cli_params)

  powwed := exp(number__times__number(num2, ln(num1, cli_params), cli_params), cli_params)

  if expNeg {
    powwed = number__divide__number(one, powwed, cli_params)
  }
  if neg {
    powwed = number__times__number(powwed, neg_one, cli_params)
  }

  //algorithm is exp(n2 ln n1)
  return powwed
}
