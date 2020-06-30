package interpreter

import "strconv"
import "math"

func num_normalize(num Action) string {

  /*
  ALGORITHM TO NORMALIZE:
      Starting with this number:
        [3412, -9912, 0001]
      STEP 1:
        loop through each number (from from decimal to integer)
        in each iteration of the loop, if the number is the opposite of `isNeg` (meaning if `isNeg` is false, then the current value should be positive and vice versa)
        use the following expression to get the complement
          `MAX_DIGIT` + 1 - |`current num`|
        replace the `current num` with this new value.
        next, if `isNeg`, the digit to the right should be added by one, otherwise, subtract it by 1.
        go to the next value and repeat
      STEP 2:
        join the vector of the integer and decimal with '.', then join each digit with ''
        if `isNeg` then precede the string with a '-'
        finally, return the result
  */

  integer := num.Integer
  decimal := num.Decimal

  //the first digit is actually the last index
  //because omm numbers are stored as so [1234, 5678, 9101] = 910, 156, 781, 234

  if isEqual(num, zero) {
    return "0"
  }

  //alloc amounts into the copies
  var decimalCopy = make([]int64, len(decimal))
  var integerCopy = make([]int64, len(integer))

  var isNeg bool = isLess(num, zero)

  var carry int64 = 0

  for k := range decimal {
    curIsNeg := decimal[k] < 0

    decimalCopy[k] = decimal[k]
    decimalCopy[k]+=carry
    carry = 0

    if decimalCopy[k] == 0 {
      continue
    }

    if curIsNeg != isNeg {
      complement := MAX_DIGIT + 1 - decimalCopy[k]
      decimalCopy[k] = complement

      if isNeg {
        carry = 1
      } else {
        carry = -1
      }
    }

    decimalCopy[k] = int64(math.Abs(float64(decimalCopy[k]))) //set the current digit to |current digit|
  }

  for k := range integer {
    curIsNeg := integer[k] < 0

    integerCopy[k] = integer[k]
    integerCopy[k]+=carry
    carry = 0

    if integerCopy[k] == 0 {
      continue
    }

    if curIsNeg != isNeg {
      complement := MAX_DIGIT + 1 - integerCopy[k]
      integerCopy[k] = complement

      if isNeg {
        carry = 1
      } else {
        carry = -1
      }
    }

    integerCopy[k] = int64(math.Abs(float64(integerCopy[k]))) //set the current digit to |current digit|
  }

  for ;len(integerCopy) != 0 && integerCopy[len(integerCopy) - 1] == 0; {
    integerCopy = integerCopy[:len(integerCopy) - 1]
  }
  for ;len(decimalCopy) != 0 && decimalCopy[0] == 0; {
    decimalCopy = decimalCopy[1:]
  }

  var joined = ""

  for _, v := range decimalCopy {
    joined = strconv.FormatInt(v, 10) + joined
  }
  if len(decimalCopy) != 0 {
    joined = "." + joined
  }
  for _, v := range integerCopy {
    joined = strconv.FormatInt(v, 10) + joined
  }

  if len(joined) == 0 { //just in case it is all zero
    return "0"
  }

  if joined[0] == '.' {
    /*
      if joined looks like: .123
      convert it to 0.123
    */
    joined = "0" + joined
  }
  if isNeg {
    joined = "-" + joined
  }

  return joined
}
