package interpreter

import "strconv"
import "math"

func num_normalize(num Action) string {

  /*

  ALGORITHM TO NORMALIZE:

      Starting with this number:
        [3412, -9912, 0001]

      STEP 0: (initializer step)
        remove the leading zeros in the integer

        if the first digit non zero digit is negative set isNeg = true, otherwise isNeg = false
        (if there is no first digit, go to the decimals and see)

      STEP 1:
        loop through each number (from from decimal to integer)
        in each iteration of the loop, if the number is the opposite of `isNeg` (meaning if `isNeg` is false, then the current value should be positive and vice versa)
        use the following expression to get the complement

          `OMM_MAX_DIGIT` - |`current num`|

        replace the `current num` with this new value.
        next, if `isNeg`, the digit to the left should be added by one, otherwise, subtract it by 1.
        go to the next value and repeat

      STEP 2:
        join the vector of the integer and decimal with '.', then join each digit with ''
        if `isNeg` then precede the string with a '-'
        finally, return the result

  */

  //determine if the number is + or -
  merged := append(num.Integer, num.Decimal...)
  mergedLen := len(merged)

  if mergedLen == 0 {
    return "0"
  }

  firstNonZeroDigitIndex := 0
  firstNonZeroDigit := merged[firstNonZeroDigitIndex]

  for ;firstNonZeroDigit == 0; {
    firstNonZeroDigitIndex++

    //if there are no "non-zero" digits, the final value is "0"
    if firstNonZeroDigitIndex >= mergedLen {
      return "0"
    }

    firstNonZeroDigit = merged[firstNonZeroDigitIndex]
  }

  var isNegative bool = false

  if firstNonZeroDigit < 0 {
    isNegative = true
  }

  var isNegativeNum int64 = 1

  if isNegative {
    isNegativeNum = -1
  }

  var numStr string = ""
  var carry int64

  //do the decimal
  deci := num.Decimal

  for i := len(deci) - 1; i >= 0; i-- {
    deci[i]+=carry
    curNeg := deci[i] < 0
    carry = 0

    if curNeg != isNegative && deci[i] != 0 {
      carry+=isNegativeNum
    }

    numStr+=strconv.Itoa(int(math.Abs(float64(deci[i]))))
  }

  //do the integer
  inte := num.Integer

  for i := len(inte) - 1; i >= 0; i-- {
    inte[i]+=carry
    curNeg := inte[i] < 0
    carry = 0

    if curNeg != isNegative && inte[i] != 0 {
      inte[i] = int64(MAX_DIGIT) + 1 - int64(math.Abs(float64(inte[i])))
      carry+=isNegativeNum * isNegativeNum
    }

    numStr = strconv.Itoa(int(math.Abs(float64(inte[i])))) + numStr
  }

  numStr = "-" + numStr

  return numStr
}
