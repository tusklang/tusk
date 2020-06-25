package interpreter

import "strconv"
import "strings"

func number__times__number(num1, num2 Action, cli_params CliParams) Action {

  var multFin [][]int64 //store the final values that were multiplied
  trailingZeroCount := 0

  //amount of decimal places there are
  decPlaceCount := len(num1.Decimal) + len(num2.Decimal)

  //actual numbers for num1 and num2
  num1n := append(num1.Decimal, num1.Integer...)
  num2n := append(num2.Decimal, num2.Integer...)

  if len(num1n) < len(num2n) { //swap num1n and num2n if num2n is greater (improves performance)
    num1n, num2n = num2n, num1n
  }

  for _, v := range num1n {
    //add a new row to the multiplied values
    multFin = append(multFin, []int64{})

    var carry int64 = 0 //variable to carry over overflowed numbers

    for i := 0; i < trailingZeroCount; i++ { //insert the trailing zeros
      multFin[len(multFin) - 1] = append(multFin[len(multFin) - 1], 0)
    }

    for _, sv := range num2n {
      var product int64 = v * sv + carry
      carry = 0 //reverrt carry after it was factored in

      if product >= MAX_DIGIT {
        var rounded int64 = product / (MAX_DIGIT + 1) * (MAX_DIGIT + 1) //round down by `MAX_DIGIT`
        carry = product / MAX_DIGIT //divide by the max digit to get the carry
        product-=rounded
      }
      if product <= MIN_DIGIT {
        var rounded int64 = ((product + (MIN_DIGIT - 1) - 1) / (MIN_DIGIT - 1)) * (MIN_DIGIT - 1) //round up by `MIN_DIGIT`
        carry = product / MAX_DIGIT //divide by the max digit to get the carry
        product-=rounded
      }

      multFin[len(multFin) - 1] = append(multFin[len(multFin) - 1], product)
    }

    multFin[len(multFin) - 1] = append(multFin[len(multFin) - 1], carry)
    trailingZeroCount++
  }

  totalSum := []int64{ 0 }

  //multiply the values
  for _, v := range multFin {
    totalSumAct := zero //placeholder number to pass into add
    totalSumAct.Integer = totalSum

    multFinAct := zero
    multFinAct.Integer = v

    totalSum = number__plus__number(totalSumAct, multFinAct, cli_params).Integer
  }

  decimalRet := totalSum[:decPlaceCount]
  integerRet := totalSum[decPlaceCount:]

  returner := zero
  returner.Integer, returner.Decimal = integerRet, decimalRet

  return returner
}

func string__times__number(num1, num2 Action, cli_params CliParams) Action {

  //ensure that the string is num1 and the number is num2
  if num2.Type == "string" {
    num1, num2 = num2, num1
  }

  finalAct := emptyString

  for i := zero; isLess(i, num2); i = number__plus__number(i, one, cli_params) {
    finalAct.ExpStr+=num1.ExpStr
  }

  i := zero
  for _, v := range finalAct.ExpStr {
    var curRune = emptyRune
    curRune.ExpStr = string(v)
    finalAct.Hash_Values[strings.TrimPrefix(num_normalize(i), ".0") /* remove the ".0" from the end */ ] = []Action{ curRune }
  }

  return finalAct
}

func array__times__array(num1, num2 Action, cli_params CliParams) Action {
  var arrLen = len(num1.Hash_Values)

  for k, v := range num2.Hash_Values {
    intK, _ := strconv.Atoi(k)
    num1.Hash_Values[strconv.Itoa(arrLen + intK)] = v
  }

  return num1
}
