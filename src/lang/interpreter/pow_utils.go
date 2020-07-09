package interpreter

//file that has all of the helper funcs for exponentiation

import "strconv"

import . "lang/types"

func number__pow__integer(val1, val2 OmmType, cli_params CliParams) OmmType {
  num1, num2 := val1.(OmmNumber), val2.(OmmNumber)
  ensurePrec(&num1, &num2, cli_params)

  //using binary exponentiation
  //https://cp-algorithms.com/algebra/binary-exp.html#toc-tgt-1

  if isEqual(num2, zero) {
    return one
  }

  var two = zero
  *two.Integer = []int64{ 2 }

  divved := number__divide__number(num2, two, cli_params).(OmmNumber)
  *divved.Decimal = nil //round down to nearest whole

  res := number__pow__integer(num1, divved, cli_params).(OmmNumber)

  *two.Integer = []int64{ 2 } //because it gets mutated

  resSquared := number__times__number(res, res, cli_params).(OmmNumber)

  if isEqual(number__mod__number(num2, two, cli_params).(OmmNumber), one) {
    return number__times__number(resSquared, num1, cli_params).(OmmNumber)
  }

  return resSquared
}

func ln(val OmmType, cli_params CliParams) OmmType {
  x := val.(OmmNumber)
  ensurePrec(&x, &OmmNumber{}, cli_params)

  //using taylor series expansion to calculate
  //found here https://www.efunda.com/math/taylor_series/logarithmic.cfm
  //algorithm 2

  var series = zero

  var two = zero
  *two.Integer = []int64{ 2, 0 }

  //calculate (x - 1) / (x + 1)
  xm1dxp1 := number__divide__number(number__minus__number(x, one, cli_params), number__plus__number(x, one, cli_params), cli_params).(OmmNumber)

  //convert precision to omm number
  ommNumberPrec := zero
  *ommNumberPrec.Integer, *ommNumberPrec.Decimal = BigNumConverter(strconv.Itoa(cli_params["Calc"]["PREC"].(int)))

  //calculate taylor series to prec
  for i := one; isLess(i, ommNumberPrec); i = number__plus__number(i, two, cli_params).(OmmNumber) {

    //calculate 1/i
    onedi := number__divide__number(one, i, cli_params)

    //calculate xm1dxp1 ^ i
    xm1dxp1pi := number__pow__integer(xm1dxp1, i, cli_params)

    //calculate onedi * xm1dxp1pi
    oneditxm1dxp1pi := number__times__number(onedi, xm1dxp1pi, cli_params)

    //add to the series
    series = number__plus__number(series, oneditxm1dxp1pi, cli_params).(OmmNumber)
  }

  series = number__times__number(series, two, cli_params).(OmmNumber)
  return series
}

func fac(val OmmType, cli_params CliParams) OmmType {
  x := val.(OmmNumber)
  ensurePrec(&x, &OmmNumber{}, cli_params)

  //factorial function for taylor series
  //using a naive method, but there is probably a faster method

  prod := one

  for i := one; isLessOrEqual(i, x); i = number__plus__number(i, one, cli_params).(OmmNumber) {
    prod = number__times__number(prod, i, cli_params).(OmmNumber)
  }

  return prod
}

func exp(val OmmType, cli_params CliParams) OmmType {
  x := val.(OmmNumber)
  ensurePrec(&x, &OmmNumber{}, cli_params)

  //using taylor series expansion to calculate
  //found here https://www.efunda.com/math/taylor_series/exponential.cfm
  //algorithm 1

  var onePlaceholder = zero //temp value for one (because one wil get mutated if it is passed directly)
  *onePlaceholder.Integer = []int64{ 1 }

  var series OmmNumber = one

  ommNumberPrec := zero
  *ommNumberPrec.Integer, *ommNumberPrec.Decimal = BigNumConverter(strconv.Itoa(cli_params["Calc"]["PREC"].(int)))

  for i := one; isLess(i, ommNumberPrec); i = number__plus__number(i, one, cli_params).(OmmNumber) {
    //calculate i!
    i_factorial := fac(i, cli_params)

    //calculate x^i
    xpi := number__pow__integer(x, i, cli_params)

    //calculate x ^ i / (i!)
    xpidifac := number__divide__number(xpi, i_factorial, cli_params)

    //add x ^ i / (i!) to the series
    series = number__plus__number(series, xpidifac, cli_params).(OmmNumber)
  }

  return series
}
