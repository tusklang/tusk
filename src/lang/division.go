package lang

import "encoding/json"
import "strings"

// #cgo CFLAGS: -std=c99
import "C"

//export DivisionC
func DivisionC(_num1C *C.char, _num2C *C.char, cli_paramsP *C.char) *C.char {

  cli_params_str := C.GoString(cli_paramsP)

  var cli_params map[string]map[string]interface{}

  _ = json.Unmarshal([]byte(cli_params_str), &cli_params)

  _num1, _num2 := C.GoString(_num1C), C.GoString(_num2C)

  if _num1 == "undef" || _num2 == "undef" {
    return C.CString("undef")
  }

  if returnInit(_num2) == "0" || returnInit(_num2) == "-0" {
    return C.CString("undef")
  }
  if returnInit(_num1) == "0" || returnInit(_num1) == "-0" {
    return C.CString("0")
  }

  _num1 = addDec(returnInit(_num1))
  _num2 = addDec(returnInit(_num2))

  decO1Index := strings.Index(_num1, ".")
  decO2Index := len(_num2) - strings.Index(_num2, ".") - 1
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

  _num1+=strings.Repeat("0", cli_params["Calc"]["PREC"].(int))

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

    for ;isLess(C.GoString(AddC(C.CString(curDivisor), C.CString(_num2), cli_paramsP)), curVal) || returnInit(C.GoString(AddC(C.CString(curDivisor), C.CString(_num2), cli_paramsP))) == returnInit(curVal); {
      curDivisor = C.GoString(AddC(C.CString(curDivisor), C.CString(_num2), cli_paramsP))
      curQ = C.GoString(AddC(C.CString(curQ), C.CString("1"), cli_paramsP))
    }

    curVal = C.GoString(SubtractC(C.CString(curVal), C.CString(curDivisor), cli_paramsP))

    final+=curQ
  }

  for ;len(final) < combinedIndex; final = "0" + final {}

  final = final[:combinedIndex] + "." + final[combinedIndex:]

  if neg {
    final = "-" + final
  }

  return C.CString(final)
}
