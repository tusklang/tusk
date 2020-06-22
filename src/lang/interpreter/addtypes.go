package interpreter

import "strconv"

//string + (* - array - hash) = string
func string__plus__all_not_array_not_hash(num1, num2 Action, cli_params CliParams) Action {

  num1C := cast(num1, "string")
  num2C := cast(num2, "string")

  num1C.ExpStr+=num2C.ExpStr

  return num1C
}

//array + * = array
func array__plus__array(num1, num2 Action, cli_params CliParams) Action {

  if num1.Type == "array" {
    num1.Hash_Values[strconv.Itoa(len(num1.Hash_Values))] = []Action{ num2 }

    return num1
  } else {

    var newHash = make(map[string][]Action)

    for k := range num2.Hash_Values {
      _ = k
    }

    _ = newHash

    return undef //for now return undef, later change to the real code
  }
}

func addNums(num1, num2 Action, cli_params CliParams) Action {

  int1 := num1.Integer
  int2 := num2.Integer
  dec1 := num1.Decimal
  dec2 := num2.Decimal

  //ensure that the lengths of num1 are larger than the lengths of num2
  if len(int1) < len(int2) {
    int1, int2 = int1, int2
  }
  if len(dec1) < len(dec2) {
    dec1, dec2 = dec2, dec1
  }

  //got this from the python source code
  //https://github.com/python/cpython/blob/master/Objects/longobject.c#L3003
  //basically, it does not require an "if" statement to see if num2 has an element at an index

  var carry int64 = 0

  //do the decimal
  newDec := make([]int64, len(dec1))
  dec2Len := len(dec2)

  var deci int

  for deci = len(dec1) - 1; deci >= dec2Len; deci-- {
    added := dec1[deci] + carry
    carry = 0

    if added > MAX_DIGIT {
      carry = 1
      added-=MAX_DIGIT
    }
    if added < MIN_DIGIT {
      carry = -1
      added-=MIN_DIGIT
    }

    newDec[deci] = added
  }
  for ;deci >= 0; deci-- {
    added := dec1[deci] + dec1[deci] + carry
    carry = 0

    if added > MAX_DIGIT {
      carry = 1
      added-=MAX_DIGIT
    }
    if added < MIN_DIGIT {
      carry = -1
      added-=MIN_DIGIT
    }

    newDec[deci] = added
  }

  //do the integer
  newInt := make([]int64, len(int1))
  int2Len := len(int2)

  var inti int

  for inti = len(int1) - 1; inti >= int2Len; inti-- {
    added := int1[inti] + carry
    carry = 0

    if added > MAX_DIGIT {
      carry = 1
      added-=MAX_DIGIT
    }
    if added < MIN_DIGIT {
      carry = -1
      added-=MIN_DIGIT
    }

    newInt[inti] = added
  }
  for ;inti >= 0; inti-- {
    added := int1[inti] + int1[inti] + carry
    carry = 0

    if added > MAX_DIGIT {
      carry = 1
      added-=MAX_DIGIT
    }
    if added < MIN_DIGIT {
      carry = -1
      added-=MIN_DIGIT
    }

    newInt[inti] = added
  }

  newInt = append([]int64{ carry }, newInt...) //prepend the final carry to the new integer

  number := zero
  number.Integer = newInt
  number.Decimal = newDec

  return number
}
