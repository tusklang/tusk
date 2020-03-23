package main

import "strings"
import "encoding/json"

// #cgo CFLAGS: -std=c99
import "C"

//variable value should be any number > 1
const MULT_THRESH_LEN = 6;

//export Multiply
func Multiply(_num1P *C.char, _num2P *C.char, calc_paramsP *C.char, line_ C.int) *C.char {

  _num1 := returnInit(C.GoString(_num1P))
  _num2 := returnInit(C.GoString(_num2P))
  calc_params_str := C.GoString(calc_paramsP)

  line := int(line_)

  _ = line

  var calc_params paramCalcOpts

  _ = json.Unmarshal([]byte(calc_params_str), &calc_params)

  if returnInit(_num1) == "0" || returnInit(_num2) == "0" {
    return C.CString("0")
  }

  decIndex := 0

  if strings.Contains(_num1, ".") {
    decIndex+=len(strings.Replace(_num1, ".", "", 1)) - strings.Index(_num1, ".")
  }
  if strings.Contains(_num2, ".") {
    decIndex+=len(strings.Replace(_num2, ".", "", 1)) - strings.Index(_num2, ".")
  }

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

  if len(_num1) < len(_num2) {
    _num1, _num2 = _num2, _num1
  }

  if len(_num1) >= MULT_THRESH_LEN && len(_num2) >= MULT_THRESH_LEN {

    final := "0"

    for i := len(_num2) - 1; i >= 0; i-- {
      mult := C.GoString(Multiply(C.CString(_num1), C.CString(_num2[i:i + 1]), calc_paramsP, line_)) + strings.Repeat("0", len(_num2) - i - 1)
      final = C.GoString(Add(C.CString(final), C.CString(mult), calc_paramsP, line_))
    }

    nNum = final
  } else {

    nNum = "0"

    for ;returnInit(_num2) != "0"; {
      nNum = C.GoString(Add(C.CString(nNum), C.CString(_num1), calc_paramsP, line_))
      _num2 = C.GoString(Subtract(C.CString(_num2), C.CString("1"), calc_paramsP, line_))
    }
  }

  if decIndex > 0 {

    nNum = Reverse(nNum)

    nNum = nNum[:decIndex] + "." + nNum[decIndex:]

    nNum = Reverse(nNum)
  }

  if neg == true {
    nNum = "-" + nNum
  }

  return C.CString(returnInit(nNum))
}
