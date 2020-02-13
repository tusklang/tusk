package main

import "strings"
import "fmt"

func division(_num1 string, _num2 string, calc_params paramCalcOpts, line uint64, functions []Funcs) string {
  if returnInit(_num2) == "0" || returnInit(_num2) == "-0"{
    return "undefined"
  }
  if returnInit(_num1) == "0" || returnInit(_num1) == "-0" {
    return "0"
  }

  _num1 = addDec(returnInit(_num1))
  _num2 = addDec(returnInit(_num2))

  decO1Index := strings.Index(_num1, ".")
  decO2Index := strings.Index(_num2, ".")
  combinedIndex := decO1Index + decO2Index

  _num1 = strings.Replace(_num1, ".", "", 1)
  _num2 = strings.Replace(_num2, ".", "", 1)

  var neg = false

  if strings.HasPrefix(_num1, "-") {
    _num1 = _num1[1:]
    neg = !neg
  }

  if strings.HasPrefix(_num2, "-") {
    _num2 = _num2[1:]
    neg = !neg
  }

  for ;len(_num2) > len(_num1); {
    _num1+="000000000000000000000000000000000000000000000"
  }

  zeroesRep := strings.Repeat("0", calc_params.precision)

  _num1+=zeroesRep

  curVal := ""
  final := ""

  for i := 0; i < len(_num1); i++ {
    curVal+=string([]rune(_num1)[i])

    if calc_params.logger {
      fmt.Println("Omm Logger ~ Division: " + final)
    }

    if isLess(curVal, _num2) {
      final+="0"
      continue
    }

    curDivisor := _num2
    curQ := "1"

    for ;isLess(add(curDivisor, _num2, calc_params, line, functions), curVal) || returnInit(add(curDivisor, _num2, calc_params, line, functions)) == returnInit(curVal); {
      curDivisor = add(curDivisor, _num2, calc_params, line, functions)
      curQ = add(curQ, "1", calc_params, line, functions)
    }

    curVal = subtract(curVal, curDivisor, calc_params, line, functions)
    final+=curQ
  }

  final = final[:combinedIndex] + "." + final[combinedIndex:]

  return returnInit(final)
}
