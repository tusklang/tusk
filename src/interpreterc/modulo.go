package main

import "strings"
import "encoding/json"

// #cgo CFLAGS: -std=c99
import "C"

//export Modulo
func Modulo(_num1P *C.char, _num2P *C.char, calc_paramsP *C.char, line_ C.int) *C.char {

  _num1 := C.GoString(_num1P)
  _num2 := C.GoString(_num2P)
  calc_params_str := C.GoString(calc_paramsP)

  line := int(line_)

  _ = line

  var calc_params paramCalcOpts

  _ = json.Unmarshal([]byte(calc_params_str), &calc_params)

  if _num2 == "0"  {
    return C.CString("undefined")
  }

  if returnInit(_num1) == "0" {
    return C.CString("0")
  }

  calc_params.precision = 0;

  divved_ := addDec(C.GoString(Division(C.CString(_num1), C.CString(_num2), calc_paramsP, line_)))

  divved := divved_[:strings.Index(divved_, ".")]

  mult := C.GoString(Multiply(C.CString(divved), C.CString(_num2), calc_paramsP, line_))
  remainder := C.GoString(Subtract(C.CString(_num1), C.CString(mult), calc_paramsP, line_))

  return C.CString(remainder)
}
