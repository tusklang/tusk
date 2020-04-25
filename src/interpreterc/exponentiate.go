package main

import "encoding/json"
import "strings"

// #cgo CFLAGS: -std=c99
import "C"

func exponentiate(_num1 string, _num2 string, calc_params paramCalcOpts, line int) string {

  _num1 = returnInit(_num1)
  _num2 = returnInit(_num2)

  if strings.Contains(_num2, ".") {
    return "NaN"
  }

  var final = "1"

  if strings.HasPrefix(_num2, "-") {
    _num2 = _num2[1:]

    for ;isLess("0", _num2); {
      final = multiply(final, _num1, calc_params, line)
      _num2 = subtract(_num2, "1", calc_params, line)
    }

    final = divide("1", final, calc_params, line)
  } else {

    for ;isLess("0", _num2); {

      final = multiply(final, _num1, calc_params, line)
      _num2 = subtract(_num2, "1", calc_params, line)
    }
  }

  return final
}

//export Exponentiate
func Exponentiate(_num1P *C.char, _num2P *C.char, calc_paramsP *C.char, line_ C.int) *C.char {

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

    num ^ num = num
    default = num
  */

  nums := TypeOperations{ _num1P_.Type, _num2P_.Type }

  var finalRet Action

  switch nums {
    case TypeOperations{ "number", "number" }: //detect case "num" ^ "num"
      val := exponentiate(_num1P_.ExpStr[0], _num2P_.ExpStr[0], calc_params, line)

      finalRet = Action{ "number", "", []string{ val }, []Action{}, []string{}, []Action{}, []Condition{}, 39, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), false }
    default: finalRet = Action{ "number", "", []string{ "0" }, []Action{}, []string{}, []Action{}, []Condition{}, 39, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), false }
  }

  reCalc(&finalRet)

  jsonNum, _ := json.Marshal(finalRet)

  return C.CString(string(jsonNum))
}
