package main

import "strings"
import "os"
import "fmt"
import "encoding/json"

// #cgo CFLAGS: -std=c99
import "C"

//export Exponentiate
func Exponentiate(_num1P *C.char, _num2P *C.char, calc_paramsP *C.char, line_ C.int) *C.char {

  _num1 := C.GoString(_num1P)
  _num2 := C.GoString(_num2P)
  calc_params_str := C.GoString(calc_paramsP)

  line := int(line_)

  _ = line

  var calc_params paramCalcOpts

  _ = json.Unmarshal([]byte(calc_params_str), &calc_params)

  _num1 = returnInit(_num1)
  _num2 = returnInit(_num2)

  if strings.Contains(_num2, ".") {
    fmt.Println("There Was An Error: Currently You Cannot Exponentiate By Numbers With Decimals\n\n" + _num1 + "^" + _num2 + "\n" +"^^^ <- Error On Line " + string(line))
    os.Exit(1)
  }

  var final = "1"

  if strings.HasPrefix(_num2, "-") {
    _num2 = _num2[1:]

    for ;isLess("0", _num2); {
      final = C.GoString(Multiply(C.CString(final), C.CString(_num1), calc_paramsP, line_))
      _num2 = C.GoString(Subtract(C.CString(_num2), C.CString("1"), calc_paramsP, line_))
    }

    final = C.GoString(Division(C.CString("1"), C.CString(final), calc_paramsP, line_))
  } else {

    for ;isLess("0", _num2); {

      final = C.GoString(Multiply(C.CString(final), C.CString(_num1), calc_paramsP, line_))
      _num2 = C.GoString(Subtract(C.CString(_num2), C.CString("1"), calc_paramsP, line_))
    }
  }

  return C.CString(final)
}
