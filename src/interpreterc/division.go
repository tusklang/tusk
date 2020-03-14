package main

import "strings"
import "encoding/json"

// #cgo CFLAGS: -std=c99
import "C"

//export Division
func Division(_num1P *C.char, _num2P *C.char, calc_paramsP *C.char, line_ C.int) *C.char {

  _num1 := C.GoString(_num1P)
  _num2 := C.GoString(_num2P)
  calc_params_str := C.GoString(calc_paramsP)

  line := int(line_)

  _ = line

  var calc_params paramCalcOpts

  _ = json.Unmarshal([]byte(calc_params_str), &calc_params)

  if returnInit(_num2) == "0" || returnInit(_num2) == "-0"{
    return C.CString("undefined")
  }
  if returnInit(_num1) == "0" || returnInit(_num1) == "-0" {
    return C.CString("0")
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
    _num1+=strings.Repeat("0", 20)
  }

  curVal := ""
  final := ""

  for i := 0; i < len(_num1); i++ {
    curVal+=string([]rune(_num1)[i])

    if isLess(curVal, _num2) {
      final+="0"
      continue
    }

    curDivisor := _num2
    curQ := "1"

    for ;isLess(C.GoString(Add(C.CString(curDivisor), C.CString(_num2), calc_paramsP, line_)), curVal) || returnInit(C.GoString(Add(C.CString(curDivisor), C.CString(_num2), calc_paramsP, line_))) == returnInit(curVal); {
      curDivisor = C.GoString(Add(C.CString(curDivisor), C.CString(_num2), calc_paramsP, line_))
      curQ = C.GoString(Add(C.CString(curQ), C.CString("1"), calc_paramsP, line_))
    }

    curVal = C.GoString(Subtract(C.CString(curVal), C.CString(curDivisor), calc_paramsP, line_))
    final+=curQ
  }

  for len(final) < combinedIndex {
    final = "0" + final;
  }

  final = final[:combinedIndex] + "." + final[combinedIndex:]

  return C.CString(returnInit(final))
}
