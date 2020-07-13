package interpreter

import . "lang/types"

func number__pow__number(val1, val2 OmmType, cli_params CliParams, stacktrace []string, line uint64, file string) *OmmType {
  num1, num2 := val1.(OmmNumber), val2.(OmmNumber)
  ensurePrec(&num1, &num2, cli_params)

  expNeg := false
  if isLess(num2, zero) { //account for negative exponents
    expNeg = true
    num2 = (*number__times__number(num2, neg_one, cli_params, stacktrace, line, file)).(OmmNumber)
  }

  if len(*num2.Decimal) == 0 { //if the exponent is an integer, use binary exponentiation for an O(log n) solution

    powwed := number__pow__integer(num1, num2, cli_params, stacktrace, line, file)

    if expNeg {
      powwed = *number__divide__number(one, powwed, cli_params, stacktrace, line, file)
    }

    return &powwed
  }

  var two = zero
  two.Integer = &[]int64{ 2 }

  neg := false
  if isLess(num1, zero) && isEqual((*number__mod__number(num1, two, cli_params, stacktrace, line, file)).(OmmNumber), zero) { //because ln (n < 0) is undefined
    neg = true
  }

  num1 = abs(num1, stacktrace, cli_params).(OmmNumber)

  powwed := exp((*number__times__number(num2, ln(num1, cli_params, stacktrace, line, file), cli_params, stacktrace, line, file)), cli_params, stacktrace, line, file)

  if expNeg {
    powwed = *number__divide__number(one, powwed, cli_params, stacktrace, line, file)
  }
  if neg {
    powwed = *number__times__number(powwed, neg_one, cli_params, stacktrace, line, file)
  }

  //algorithm is exp(n2 ln n1)
  return &powwed
}
