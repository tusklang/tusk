package main

import "math/big"
import "encoding/json"
import "fmt"

// #cgo CFLAGS: -std=c99
import "C"

func modulo(num1 string, num2 string, calc_params paramCalcOpts, line int) string {
  if num2 == "0"  {
    return "undefined"
  }

  divved := divide(num1, num2, calc_params, line)

  floored, _ := new(big.Int).SetString(divved, 10)

  mult := multiply(fmt.Sprint(floored), num2, calc_params, line)
  remainder := subtract(num1, mult, calc_params, line)

  return returnInit(fmt.Sprintf("%f", remainder))
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

      finalRet = Action{ "number", "", []string{ val }, []Action{}, []string{}, []Action{}, []Condition{}, 39, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), false }
    default: finalRet = Action{ "falsey", "", []string{ "undefined" }, []Action{}, []string{}, []Action{}, []Condition{}, 41, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), false }
  }

  reCalc(&finalRet)

  jsonNum, _ := json.Marshal(finalRet)

  return C.CString(string(jsonNum))
}
