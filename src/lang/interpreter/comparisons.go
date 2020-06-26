package interpreter

import "math/big"

func isTruthy(val Action) bool {
  return !(val.ExpStr == "false" || val.Type == "falsey")
}

func toBig(val1, val2 Action) (*big.Int, *big.Int, *big.Int, *big.Int) {
  int1, dec1 := val1.Integer, val1.Decimal
  int2, dec2 := val2.Integer, val2.Decimal

  //convert the numbers (integers/decimals) to bigints
  numConv := func(slice []int64) *big.Int {
    bigVal := big.NewInt(0)
    places := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(len(slice) - 1)), nil)

    for i := len(slice) - 1; i >= 0; i-- {
      bigVal = bigVal.Add(bigVal, new(big.Int).Mul(places, big.NewInt(slice[i])))
      places.Div(places, big.NewInt(10)) //divide the place by 10
    }

    return bigVal
  }

  var (
    bigInt1 = numConv(int1)
    bigInt2 = numConv(int2)
    bigDec1 = numConv(dec1)
    bigDec2 = numConv(dec2)
  )

  return bigInt1, bigInt2, bigDec1, bigDec2
}

func isLess(val1, val2 Action) bool {

  bigint1, bigint2, bigdec1, bigdec2 := toBig(val1, val2)

  if bigint1.Cmp(bigint2) != 0 {
    return bigint1.Cmp(bigint2) == -1
  }

  return bigdec1.Cmp(bigdec2) == -1
}

func isEqual(val1, val2 Action) bool {

  bigint1, bigint2, bigdec1, bigdec2 := toBig(val1, val2)

  if bigint1.Cmp(bigint2) != 0 {
    return false
  }

  return bigdec1.Cmp(bigdec2) == 0
}

func isLessOrEqual(val1, val2 Action) bool {

  bigint1, bigint2, bigdec1, bigdec2 := toBig(val1, val2)

  if bigint1.Cmp(bigint2) != 0 {
    return bigint1.Cmp(bigint2) == -1
  }

  return bigdec1.Cmp(bigdec2) <= 0
}

func abs(val Action, cli_params CliParams) Action {

  if isLess(val, zero) {
    return number__times__number(val, neg_one, cli_params)
  }

  return val
}
