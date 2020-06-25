package interpreter

func isTruthy(val Action) bool {
  return !(val.ExpStr == "false" || val.Type == "falsey")
}

func isLess(val1, val2 Action) bool {

  //loop through the integer and decimal
  //if it is less, val1 is less, if it is greater, val1 is greater

  //if a swap happened (for the integer and the decimal)
  var swappedInt bool = false
  var swappedDec bool = false

  //ensure that the values of val1 are greater than the values of val2
  if len(val1.Integer) < len(val2.Integer) {
    swappedInt = true
    val1.Integer, val2.Integer = val2.Integer, val1.Integer
  }
  if len(val1.Decimal) < len(val2.Decimal) {
    swappedDec = true
    val1.Decimal, val2.Decimal = val2.Decimal, val1.Decimal
  }

  //do the integer
  var intI int

  for intI = len(val1.Integer) - 1; intI >= len(val2.Integer); intI-- {
    if val1.Integer[intI] < 0 {
      return !swappedInt
    } else if val1.Integer[intI] > 0 {
      return swappedInt
    }
  }
  for ;intI >= 0; intI-- {
    if val1.Integer[intI] < val2.Integer[intI] {
      return !swappedInt
    } else if val1.Integer[intI] > val2.Integer[intI] {
      return swappedInt
    }
  }

  //do the decimal
  var decI int

  for decI = len(val1.Decimal) - 1; decI >= len(val2.Decimal); decI-- {
    if val1.Decimal[decI] < 0 {
      return !swappedDec
    } else if val1.Decimal[decI] > 0 {
      return swappedDec
    }
  }
  for ;decI >= 0; decI-- {
    if val1.Decimal[decI] < val2.Decimal[decI] {
      return !swappedDec
    } else if val1.Decimal[decI] > val2.Decimal[decI] {
      return swappedDec
    }
  }

  return false //if nothing passed, return false
}

func isEqual(val1, val2 Action) bool {

  //loop through the integer and decimal, and if the current value of val1 is not equal to current value of val1, it is not equal

  //ensure that the values of val1 are greater than the values of val2
  if len(val1.Integer) < len(val2.Integer) {
    val1.Integer, val2.Integer = val2.Integer, val1.Integer
  }
  if len(val1.Decimal) < len(val2.Decimal) {
    val1.Decimal, val2.Decimal = val2.Decimal, val1.Decimal
  }

  //do the integer
  var intI int

  for intI = 0; intI < len(val2.Integer); intI++ {
    if val1.Integer[intI] != val2.Integer[intI] {
      return false
    }
  }
  for ;intI < len(val1.Integer); intI++ {
    if val1.Integer[intI] != 0 {
      return false
    }
  }

  //do the decimal
  var decI int

  for decI = 0; decI < len(val2.Decimal); decI++ {
    if val1.Decimal[decI] != val2.Decimal[decI] {
      return false
    }
  }
  for ;decI < len(val1.Decimal); decI++ {
    if val1.Decimal[decI] != 0 {
      return false
    }
  }

  return true //if it passed all tests, return true
}
