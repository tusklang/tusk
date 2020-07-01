package interpreter

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

func isTruthy(val Action) bool {
  return !(val.ExpStr == "false" || val.Type == "falsey")
}

func isLess(val1, val2 Action) bool {

  //manage for leading zeros in the decimal
  var num1ZeroLeadingDec int
  var num2ZeroLeadingDec int

  num1DecGreater := false //if num1 has less leading zeros

  for ;len(val1.Decimal) != 0 && val1.Decimal[len(val1.Decimal) - 1] == 0; {
    val1.Decimal = val1.Decimal[:len(val1.Decimal) - 1]
    num1ZeroLeadingDec++
  }
  for ;len(val2.Decimal) != 0 && val2.Decimal[len(val2.Decimal) - 1] == 0; {
    val2.Decimal = val2.Decimal[:len(val2.Decimal) - 1]
    num2ZeroLeadingDec++
  }

  if num1ZeroLeadingDec > num2ZeroLeadingDec {
    num1DecGreater = true
  }
  /////////////////////////////////////////

  bigint1, bigint2, bigdec1, bigdec2 := toBig(val1, val2)

  if bigint1.Cmp(bigint2) != 0 {
    return bigint1.Cmp(bigint2) == -1
  }

  if num1DecGreater { //if num1 has more leading zeros, return true(because the decimal is less)
    return true
  }

  return bigdec1.Cmp(bigdec2) == -1
}

func isEqual(val1, val2 Action) bool {

  //manage for leading zeros in the decimal
  var num1ZeroLeadingDec int
  var num2ZeroLeadingDec int

  num1DecNotEqual := true //if num1 has less leading zeros

  for ;len(val1.Decimal) != 0 && val1.Decimal[len(val1.Decimal) - 1] == 0; {
    val1.Decimal = val1.Decimal[:len(val1.Decimal) - 1]
    num1ZeroLeadingDec++
  }
  for ;len(val2.Decimal) != 0 && val2.Decimal[len(val2.Decimal) - 1] == 0; {
    val2.Decimal = val2.Decimal[:len(val2.Decimal) - 1]
    num2ZeroLeadingDec++
  }

  if num1ZeroLeadingDec == num2ZeroLeadingDec {
    num1DecNotEqual = false
  }
  /////////////////////////////////////////

  bigint1, bigint2, bigdec1, bigdec2 := toBig(val1, val2)

  if bigint1.Cmp(bigint2) != 0 {
    return false
  }

  if num1DecNotEqual {
    return false
  }

  return bigdec1.Cmp(bigdec2) == 0
}

func isLessOrEqual(val1, val2 Action) bool {

  //manage for leading zeros in the decimal
  var num1ZeroLeadingDec int
  var num2ZeroLeadingDec int

  num1DecLess := true //if num1 has less leading zeros

  for ;len(val1.Decimal) != 0 && val1.Decimal[len(val1.Decimal) - 1] == 0; {
    val1.Decimal = val1.Decimal[:len(val1.Decimal) - 1]
    num1ZeroLeadingDec++
  }
  for ;len(val2.Decimal) != 0 && val2.Decimal[len(val2.Decimal) - 1] == 0; {
    val2.Decimal = val2.Decimal[:len(val2.Decimal) - 1]
    num2ZeroLeadingDec++
  }

  if num1ZeroLeadingDec > num2ZeroLeadingDec {
    num1DecLess = false
  }
  /////////////////////////////////////////

  bigint1, bigint2, bigdec1, bigdec2 := toBig(val1, val2)

  if bigdec1.Cmp(bigdec2) != 0 && !num1DecLess {
    return bigint1.Cmp(bigint2) <= 0
  }

  if num1DecLess {
    return false
  }

  return bigdec1.Cmp(bigdec2) <= 0
}

func abs(val Action, cli_params CliParams) Action {

  if isLess(val, zero) {
    return number__times__number(val, neg_one, cli_params)
  }

  return val
}

//function to determine if two actions are equal
func equals(val1, val2 Action) bool {

  if val1.Type == "number" && val2.Type == "number" {
    return isEqual(val1, val2)
  } else if val1.Name == "hashed_value" && val1.Name == "hashed_value" {
    if len(val1.Hash_Values) < len(val2.Hash_Values) { //ensure val1 has more hash values
      val1, val2 = val2, val1
    }

    //loop through the keys of val1, and if val1[k] == val2[k] they are equal
    for k := range val1.Hash_Values {

      //check if val2 has k
      if _, exists := val2.Hash_Values[k]; !exists {
        return false
      }

      if !equals(val1.Hash_Values[k][0], val2.Hash_Values[k][0]) {
        return false
      }
    }

    return true //if they are all equal, return true
  } else {
    return val1.ExpStr == val2.ExpStr
  }
}
