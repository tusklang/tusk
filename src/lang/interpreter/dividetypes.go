package interpreter

func number__divide__number(num1, num2 Action, cli_params CliParams) Action {
  ensurePrec(&num1, &num2, cli_params)

  //maybe in a future version switch to the algorithm python uses
  //https://github.com/python/cpython/blob/8bd216dfede9cb2d5bedb67f20a30c99844dbfb8/Objects/longobject.c#L2610
  //because it is faster

  //num2 is the divisor
  //num1 is the dividend

  if isEqual(num2, zero) { //if it is n/0 return undef
    return undef
  }
  if isEqual(num1, zero) { //if it is 0/n return 0
    return zero
  }

  decPlaces := len(num1.Integer) + len(num2.Decimal)
  num1n := append(num1.Decimal, num1.Integer...)

  for i := len(num1n); i < cli_params["Calc"]["PREC"].(int); i++ {
    num1n = append([]int64{ 0 }, num1n...)
  }

  curVal := zero //current value under the "house" of the division
  var final []int64 //final value

  num2Abs := abs(num2, cli_params)

  for i := len(num1n) - 1; i >= 0; i-- {
    v := num1n[i]

    curVal.Integer = append([]int64{ v }, curVal.Integer...)

    curValAbs := abs(curVal, cli_params)

    if isLess(curVal, num2Abs) {
      final = append([]int64{ 0 }, final...)
      continue
    }

    var curQuotient Action = zero
    var added Action = zero

    for addedTemp := added; func() bool {
      addedTemp = number__plus__number(addedTemp, num2Abs, cli_params)
      return isLessOrEqual(addedTemp, curValAbs)
    }(); added = addedTemp {
      curQuotient = number__plus__number(curQuotient, one, cli_params) //increment the current quotient
    }

    if isLess(curVal, zero) {
      curQuotient = number__times__number(curQuotient, neg_one, cli_params)
    }

    curVal = number__minus__number(curValAbs, added, cli_params)
    final = append(curQuotient.Integer, final...)
  }

  ret := zero
  ret.Integer, ret.Decimal = final[len(final) - decPlaces:], final[:len(final) - decPlaces]

  return ret
}
