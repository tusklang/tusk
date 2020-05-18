package lang

import "encoding/json"
import "strings"

// #cgo CFLAGS: -std=c99
import "C"

//export MultiplyC
func MultiplyC(_num1C *C.char, _num2C *C.char, cli_paramsP *C.char) *C.char {

  cli_params_str := C.GoString(cli_paramsP)

  var cli_params map[string]map[string]interface{}

  _ = json.Unmarshal([]byte(cli_params_str), &cli_params)

  _num1, _num2 := C.GoString(_num1C), C.GoString(_num2C)

  if _num1 == "undef" || _num2 == "undef" {
    return C.CString("undef")
  }

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

  if float64(len(_num1)) >= cli_params["Calc"]["LONG_MULT_THRESH"].(float64) && float64(len(_num2)) >= cli_params["Calc"]["LONG_MULT_THRESH"].(float64) {

    final := "0"

    for i := len(_num2) - 1; i >= 0; i-- {
      mult := C.GoString(MultiplyC(C.CString(_num1), C.CString(_num2[i:i + 1]), cli_paramsP)) + strings.Repeat("0", len(_num2) - i - 1)
      final = C.GoString(AddC(C.CString(final), C.CString(mult), cli_paramsP))
    }

    nNum = final
  } else {

    nNum = "0"

    for ;returnInit(_num2) != "0"; {
      nNum = C.GoString(AddC(C.CString(nNum), C.CString(_num1), cli_paramsP))
      _num2 = C.GoString(SubtractC(C.CString(_num2), C.CString("1"), cli_paramsP))
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
