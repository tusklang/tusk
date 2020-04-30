package main

import "encoding/json"
import "strings"

// #cgo CFLAGS: -std=c99
import "C"

func modulo(_num1 string, _num2 string, calc_params paramCalcOpts, line int) string {

  if _num2 == "0"  {
    return "undef"
  }

  if returnInit(_num1) == "0" {
    return "0"
  }

  divved_ := addDec(divide(_num1, _num2, calc_params, line))

  divved := divved_[:strings.Index(divved_, ".")]

  mult := multiply(divved, _num2, calc_params, line)
  remainder := subtract(_num1, mult, calc_params, line)

  return remainder
}

//export Modulo
func Modulo(_num1P *C.char, _num2P *C.char, calc_paramsP *C.char, line_ C.int) *C.char {

  _num1 := C.GoString(_num1P)
  _num2 := C.GoString(_num2P)
  calc_params_str := C.GoString(calc_paramsP)

  line := int(line_)

  _ = line

  var calc_params paramCalcOpts

  _ = json.Unmarshal([]byte(calc_params_str), &calc_params)

  var _num1P_ Action
  var _num2P_ Action

  _ = json.Unmarshal([]byte(_num1), &_num1P_)
  _ = json.Unmarshal([]byte(_num2), &_num2P_)

  /* TABLE OF TYPES:

    num % num = num
    default = falsey
  */

  nums := TypeOperations{ _num1P_.Type, _num2P_.Type }

  var finalRet Action

  switch nums {
    case TypeOperations{ "number", "number" }: //detect case "num" % "num"
      val := modulo(_num1P_.ExpStr[0], _num2P_.ExpStr[0], calc_params, line)

      finalRet = Action{ "number", "", []string{ val }, []Action{}, []string{}, [][]Action{}, []Condition{}, 39, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), false }
    default: finalRet = Action{ "falsey", "", []string{ "undef" }, []Action{}, []string{}, [][]Action{}, []Condition{}, 41, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), false }
  }

  reCalc(&finalRet)

  jsonNum, _ := json.Marshal(finalRet)

  return C.CString(string(jsonNum))
}
