package main

import "strings"
import "fmt"

func multiply(_num1 string, _num2 string, calc_params paramCalcOpts, line uint64, functions []Funcs) string {

  decIndex := 0

  if strings.Contains(_num1, ".") {
    decIndex+=len(_num1) - strings.Index(_num1, ".")
  }
  if strings.Contains(_num2, ".") {
    decIndex+=len(_num1) - strings.Index(_num2, ".")
  }

  decIndex--

  _num1 = strings.Replace(_num1, ".", "", 1)
  _num2 = strings.Replace(_num2, ".", "", 1)

  if isLess(_num1, _num2) {
    _num1, _num2 = _num2, _num1
  }

  neg := false

  if strings.HasPrefix(_num1, "-") {
    _num1 = _num1[1:]
    neg = !neg
  }
  if strings.HasPrefix(_num2, "-") {
    _num2= _num2[1:]
    neg = !neg
  }

  var nNum string

  if isLess(_num2, calc_params.mult_thresh) {
    nNum = "0"

    for ;returnInit(_num2) != "0"; {
      nNum = add(nNum, _num1, calc_params, line, functions)
      _num2 = subtract(_num2, "1", calc_params, line, functions)

      if calc_params.logger {
        fmt.Println("Omm Logger ~ Multiplication: " + nNum)
      }
    }

    if decIndex != -1 {
      nNum = Reverse(nNum)

      nNum = nNum[:decIndex] + "." + nNum[decIndex:]

      nNum = Reverse(nNum)
    }

    if neg == true {
      nNum = "-" + nNum
    }
  } else {
    nNum = "0"

    for i, o := len(_num2) - 1, 0; i >= 0; i, o = i - 1, o + 1 {
      nNum = add(nNum, multiply(string([]rune(_num2)[i]), _num1, calc_params, line, functions) + RepeatAdd("0", o), calc_params, line, functions)

      if calc_params.logger {
        fmt.Println("Omm Logger ~ Multiplication: " + nNum)
      }
    }
  }

  return returnInit(nNum)
}
