package interpreter

func number__divide__number(num1, num2 Action, cli_params CliParams) Action {
  ensurePrec(&num1, &num2, cli_params)

  //maybe in a future version switch to the algorithm python uses
  //https://github.com/python/cpython/blob/8bd216dfede9cb2d5bedb67f20a30c99844dbfb8/Objects/longobject.c#L2610
  //because it is faster
  //also look into this:
  //http://citeseerx.ist.psu.edu/viewdoc/download?doi=10.1.1.47.565&rep=rep1&type=pdf

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
  num2n := zero
  num2n.Integer = append(num2.Decimal, num2.Integer...)

  a := zero
  a.Integer = num1n

  for i := len(num1n); i < cli_params["Calc"]["PREC"].(int); i++ {
    num1n = append([]int64{ 0 }, num1n...)
  }

  curVal := zero //current value under the "house" of the division
  var final []int64 //final value

  num2Abs := abs(num2n, cli_params)

  a = zero
  a.Integer = num1n

  for i := len(num1n) - 1; i >= 0; i-- {
    v := num1n[i]

    curVal.Integer = append([]int64{ v }, curVal.Integer...)
    curValAbs := abs(curVal, cli_params)

    if isLess(curValAbs, num2Abs) {
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

    apn2 := number__plus__number(added, num2Abs, cli_params)

    if isEqual(apn2, curValAbs) {
      added = apn2
      curQuotient = number__plus__number(curQuotient, one, cli_params)
    }

    if isLess(num1, zero) {
      curQuotient = number__times__number(curQuotient, neg_one, cli_params)
    }

    //remove leading zeros from the curQuotient
    for ;len(curQuotient.Integer) != 1 && curQuotient.Integer[len(curQuotient.Integer) - 1] == 0; {
      curQuotient.Integer = curQuotient.Integer[:len(curQuotient.Integer) - 1]
    }

    curVal = number__minus__number(curValAbs, added, cli_params)
    final = append(curQuotient.Integer, final...)
  }

  if isLess(num2, zero) { //if num2 is negative, multiply the final by -1
    finalAct := zero
    finalAct.Integer = final
    finalAct = number__times__number(finalAct, neg_one, cli_params)
    final = finalAct.Integer
  }

  ret := zero
  ret.Integer, ret.Decimal = final[len(final) - decPlaces:], final[:len(final) - decPlaces]

  return ret
}
