package interpreter

import "strings"
import "math/big"

//convert the numbers (integers/decimals) to bigints
func sliceToBigInt(slice []int64) *big.Int {
  bigVal := big.NewInt(0)
  places := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(len(slice) - 1)), nil)

  for i := len(slice) - 1; i >= 0; i-- {
    bigVal = bigVal.Add(bigVal, new(big.Int).Mul(places, big.NewInt(slice[i])))
    places.Div(places, big.NewInt(10)) //divide the place by 10
  }

  return bigVal
}
func sliceToBigDec(slice []int64) *big.Int {
  bigVal := big.NewInt(0)
  places := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(len(slice) - 1)), nil)

  for i := 0; i < len(slice); i++ {
    bigVal = bigVal.Add(bigVal, new(big.Int).Mul(places, big.NewInt(slice[i])))
    places.Div(places, big.NewInt(10)) //divide the place by 10
  }

  return bigVal
}

func toBig(val1, val2 Action) (*big.Int, *big.Int, *big.Int, *big.Int) {
  int1, dec1 := val1.Integer, val1.Decimal
  int2, dec2 := val2.Integer, val2.Decimal

  var (
    bigInt1 = sliceToBigInt(int1)
    bigInt2 = sliceToBigInt(int2)
    bigDec1 = sliceToBigInt(dec1)
    bigDec2 = sliceToBigInt(dec2)
  )

  return bigInt1, bigInt2, bigDec1, bigDec2
}

func num_normalize(num Action) string {

  //manage for leading zeros in the decimal
  var numZeroLeadingDec int

  for ;len(num.Decimal) != 0 && num.Decimal[len(num.Decimal) - 1] == 0; {
    num.Decimal = num.Decimal[:len(num.Decimal) - 1]
    numZeroLeadingDec++
  }
  /////////////////////////////////////////

  integer, decimal := sliceToBigInt(num.Integer), sliceToBigInt(num.Decimal)

  return integer.String() + "." + strings.Repeat("0", numZeroLeadingDec) + decimal.String()
}
