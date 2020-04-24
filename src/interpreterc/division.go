package main

import "math/big"
import "encoding/json"
import "strings"
import "fmt"

// #cgo CFLAGS: -std=c99
import "C"

func divide(num1 string, num2 string, calc_params paramCalcOpts, line int) string {
  calc := new(big.Float)

  num1big, _ := new(big.Float).SetPrec(PREC).SetString(num1)
  num2big, _ := new(big.Float).SetPrec(PREC).SetString(num2)

  return returnInit(fmt.Sprintf("%f", calc.Quo(num1big, num2big)))
}

//export Division
func Division(_num1P *C.char, _num2P *C.char, calc_paramsP *C.char, line_ C.int) *C.char {

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

    num / num = num
    string / num = string | falsey (if the number is greater than string length, return falsey)
    default = falsey
  */

  nums := TypeOperations{ _num1P_.Type, _num2P_.Type }

  var finalRet Action

  switch nums {
    case TypeOperations{ "number", "number" }: //detect case "num" / "num"
      numRet := returnInit(divide(_num1P_.ExpStr[0], _num2P_.ExpStr[0], calc_params, line))

      finalRet = Action{ "number", "", []string{ numRet }, []Action{}, []string{}, []Action{}, []Condition{}, 39, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), false }
    case TypeOperations{ "string", "number" }: //detect case "string" / "num"
      //get string length
      str_ := _num1P_.ExpStr[0]

      var length string

      for length = "1"; str_ != ""; length = add(length, "1", calc_params, line) {
        str_ = str_[1:]
      }
      ////////////////////

      subtracted := subtract(length, "2", calc_params, line)

      if isLess(subtracted, _num2P_.ExpStr[0]) || returnInit(subtracted) == _num2P_.ExpStr[0] {
        finalRet = Action{ "falsey", "", []string{ "undefined" }, []Action{}, []string{}, []Action{}, []Condition{}, 41, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), false }
      } else {

        var _str []string

        for i := subtract(subtract(length, _num2P_.ExpStr[0], calc_params, line), "2", calc_params, line); isLess("0", i); i = subtract(i, "1", calc_params, line) {
          _str = append([]string{ getIndex(_num1P_.ExpStr[0], i) }, _str...)
        }

        str := strings.Join(_str, "")

        finalRet = Action{ "string", "", []string{ str }, []Action{}, []string{}, []Action{}, []Condition{}, 38, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), false }
      }
    case TypeOperations{ "number", "string" }: //detect case "string" / "num"

      //get string length
      str_ := _num2P_.ExpStr[0]

      var length string

      for length = "1"; str_ != ""; length = add(length, "1", calc_params, line) {
        str_ = str_[1:]
      }
      ////////////////////

      subtracted := subtract(length, "2", calc_params, line)

      if isLess(subtracted, _num1P_.ExpStr[0]) || returnInit(subtracted) == _num1P_.ExpStr[0] {
        finalRet = Action{ "falsey", "", []string{ "undefined" }, []Action{}, []string{}, []Action{}, []Condition{}, 41, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), false }
      } else {

        var cur string
        str := _num2P_.ExpStr[0]

        for i := add(_num1P_.ExpStr[0], "2", calc_params, line); isLess(i, length); i = add(i, "1", calc_params, line) {
          cur+=getIndex(str, i)
        }

        finalRet = Action{ "string", "", []string{ cur }, []Action{}, []string{}, []Action{}, []Condition{}, 38, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), false }
      }
    default: finalRet = Action{ "falsey", "", []string{ "undefined" }, []Action{}, []string{}, []Action{}, []Condition{}, 41, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), false }
  }

  reCalc(&finalRet)

  jsonNum, _ := json.Marshal(finalRet)

  return C.CString(string(jsonNum))
}
