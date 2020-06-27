package interpreter

//file that has all of the helper funcs for exponentiation

import "strconv"

func number__pow__integer(num1, num2 Action, cli_params CliParams) Action {
  ensurePrec(&num1, &num2, cli_params)

  //using binary exponentiation
  //https://cp-algorithms.com/algebra/binary-exp.html#toc-tgt-1

  if isEqual(num2, zero) {
    return one
  }

  var two = zero
  two.Integer = []int64{ 2 }

  divved := number__divide__number(num2, two, cli_params)
  divved.Decimal = nil //round down to nearest whole

  res := number__pow__integer(num1, divved, cli_params)

  two.Integer = []int64{ 2 }

  resSquared := number__times__number(res, res, cli_params)

  if isEqual(number__mod__number(num2, two, cli_params), one) {
    return number__times__number(resSquared, num1, cli_params)
  }

  return resSquared
}

func ln(x Action, cli_params CliParams) Action {
  ensurePrec(&x, &Action{}, cli_params)

  //using taylor series expansion to calculate
  //taylor looks like:
  //  ln(x) = (1/1 * ((x - 1) / x)) + (1/2 * ((x - 2) / x)) ...
  //the series is 0 < x <= prec - 1

  var onePlaceholder = zero //temp value for one (because one wil get mutated if it is passed directly)
  onePlaceholder.Integer = []int64{ 1 }

  //calculate (x - 1) / x
  xm1dx := number__divide__number(number__minus__number(x, onePlaceholder, cli_params), x, cli_params)

  var series Action = zero

  ommNumberPrec := zero
  ommNumberPrec.Integer, ommNumberPrec.Decimal = BigNumConverter(strconv.Itoa(cli_params["Calc"]["PREC"].(int)))

  for i := one; isLess(i, ommNumberPrec); i = number__plus__number(i, one, cli_params) { //make async later probably

    iplaceholder := zero //i will get mutated
    iplaceholder.Integer, iplaceholder.Decimal = append([]int64{}, i.Integer...), append([]int64{}, i.Decimal...)

    //calculate 1/i
    onedi := number__divide__number(one, iplaceholder, cli_params)

    //calculate xm1dx ^ i
    xm1dxpi := number__pow__integer(xm1dx, i, cli_params)

    //calculate onedi * xm1dxpi
    onedi__mul__xm1dxpi := number__times__number(onedi, xm1dxpi, cli_params)

    //add to the series
    series = number__plus__number(series, onedi__mul__xm1dxpi, cli_params)
  }

  return series
}

func fac(x Action, cli_params CliParams) Action {
  ensurePrec(&x, &Action{}, cli_params)

  //factorial function for taylor series
  //using a naive method, but there is probably a faster method

  prod := one

  for i := one; isLessOrEqual(i, x); i = number__plus__number(i, one, cli_params) {
    prod = number__times__number(prod, i, cli_params)
  }

  return prod
}

func exp(x Action, cli_params CliParams) Action {
  ensurePrec(&x, &Action{}, cli_params)

  //using taylor series expansion to calculate
  //taylor looks like:
  //  e^x = 1 + (x ^ 1 / 1!) + (x ^ 2 / 2!) ...
  //the series is 0 < x <= prec - 1

  var onePlaceholder = zero //temp value for one (because one wil get mutated if it is passed directly)
  onePlaceholder.Integer = []int64{ 1 }

  var series Action = one

  ommNumberPrec := zero
  ommNumberPrec.Integer, ommNumberPrec.Decimal = BigNumConverter(strconv.Itoa(cli_params["Calc"]["PREC"].(int)))

  for i := one; isLess(i, ommNumberPrec); i = number__plus__number(i, one, cli_params) { //make async later probably
    //calculate i!
    i_factorial := fac(i, cli_params)

    //calculate x^i
    xpi := number__pow__integer(x, i, cli_params)
    ensurePrec(&xpi, &Action{}, cli_params)

    //calculate x ^ i / (i!)
    xpidifac := number__divide__number(xpi, i_factorial, cli_params)
    ensurePrec(&xpidifac, &Action{}, cli_params)

    //add x ^ i / (i!) to the series
    series = number__plus__number(series, xpidifac, cli_params)
    ensurePrec(&series, &Action{}, cli_params)
  }

  return series
}
