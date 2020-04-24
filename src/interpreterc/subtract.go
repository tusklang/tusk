package main

import "encoding/json"
import "math/big"
import "fmt"

// #cgo CFLAGS: -std=c99
import "C"

//export SubtractStrings
func SubtractStrings(num1, num2, calc_params *C.char, line C.int) *C.char {

  var cp paramCalcOpts

  _ = json.Unmarshal([]byte(C.GoString(calc_params)), &cp)

  sum := subtract(C.GoString(num1), C.GoString(num2), cp, int(line))

  return C.CString(sum)
}

func subtract(num1 string, num2 string, calc_params paramCalcOpts, line int) string {
  calc := new(big.Float)

  num1big, _ := new(big.Float).SetPrec(PREC).SetString(num1)
  num2big, _ := new(big.Float).SetPrec(PREC).SetString(num2)

  return returnInit(fmt.Sprintf("%f", calc.Sub(num1big, num2big)))
}

//export Subtract
func Subtract(_num1P *C.char, _num2P *C.char, calc_paramsP *C.char, line_ C.int) *C.char {

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

    num - num = num
    string - num = string | falsey (if the number is greater than string length, return falsey)
    boolean - boolean = num
    array - num = array
    num - boolean = num
    default = falsey
  */

  nums := TypeOperations{ _num1P_.Type, _num2P_.Type }

  var finalRet Action

  switch nums {
    case TypeOperations{ "number", "number" }: //detect case "num" - "num"
      val := subtract(_num1P_.ExpStr[0], _num2P_.ExpStr[0], calc_params, line)

      finalRet = Action{ "number", "", []string{ val }, []Action{}, []string{}, []Action{}, []Condition{}, 39, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), false }
    case TypeOperations{ "string", "number" }: //detect case "string" - "num"
      var val string
      str := _num1P_.ExpStr[0]

      //get string length
      str_ := str

      var length string

      for length = "1"; str_ != ""; length = add(length, "1", calc_params, line) {
        str_ = str_[1:]
      }
      ////////////////////

      subtracted := subtract(length, "2", calc_params, line)

      if isLess(subtracted, _num2P_.ExpStr[0]) || returnInit(subtracted) == _num2P_.ExpStr[0] {
        finalRet = Action{ "falsey", "", []string{ "undefined" }, []Action{}, []string{}, []Action{}, []Condition{}, 41, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), false }
      } else {

        for i := subtract(length, add(_num2P_.ExpStr[0], "1", calc_params, line), calc_params, line); isLess(i, length); i = add(i, "1", calc_params, line) {
          val+=getIndex(str, i)
        }

        finalRet = Action{ "string", "", []string{ val }, []Action{}, []string{}, []Action{}, []Condition{}, 38, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), false }
      }
    case TypeOperations{ "number", "string" }: //detect case "string" - "num"
      var val string
      str := _num2P_.ExpStr[0]

      //get string length
      str_ := str

      var length string

      for length = "1"; str_ != ""; length = add(length, "1", calc_params, line) {
        str_ = str_[1:]
      }
      ////////////////////

      subtracted := subtract(length, "2", calc_params, line)

      if isLess(subtracted, _num1P_.ExpStr[0]) || returnInit(subtracted) == _num1P_.ExpStr[0] {
        finalRet = Action{ "falsey", "", []string{ "undefined" }, []Action{}, []string{}, []Action{}, []Condition{}, 41, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), false }
      } else {

        for i := "0"; isLess(i, _num1P_.ExpStr[0]) || returnInit(i) == returnInit(_num1P_.ExpStr[0]); i = add(i, "1", calc_params, line) {
          val+=string(str[0])
          str = str[1:]
        }

        finalRet = Action{ "string", "", []string{ val }, []Action{}, []string{}, []Action{}, []Condition{}, 38, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), false }
      }
    case TypeOperations{ "boolean", "boolean" }: //detect case "boolean" - "boolean"

      var val1, val2 string

      if _num1P_.ExpStr[0] == "true" {
        val1 = "1"
      } else {
        val1 = "0"
      }

      if _num2P_.ExpStr[0] == "true" {
        val2 = "1"
      } else {
        val2 = "0"
      }

      val := subtract(val1, val2, calc_params, line)

      finalRet = Action{ "number", "", []string{ val }, []Action{}, []string{}, []Action{}, []Condition{}, 39, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), false }
    case TypeOperations{ "number", "boolean" }: //detect case "number" - "boolean"

      var val2 string

      if _num2P_.ExpStr[0] == "true" {
        val2 = "1"
      } else {
        val2 = "0"
      }

      val := subtract(_num1P_.ExpStr[0], val2, calc_params, line)

      finalRet = Action{ "number", "", []string{ val }, []Action{}, []string{}, []Action{}, []Condition{}, 39, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), false }
    case TypeOperations{ "boolean", "number" }: //detect case "number" - "boolean"

      var val1 string

      if _num1P_.ExpStr[0] == "true" {
        val1 = "1"
      } else {
        val1 = "0"
      }

      val := subtract(val1, _num2P_.ExpStr[0], calc_params, line)

      finalRet = Action{ "number", "", []string{ val }, []Action{}, []string{}, []Action{}, []Condition{}, 39, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), false }
    default: finalRet = Action{ "falsey", "", []string{ "undefined" }, []Action{}, []string{}, []Action{}, []Condition{}, 41, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), false }
  }

  reCalc(&finalRet)

  jsonNum, _ := json.Marshal(finalRet)

  return C.CString(string(jsonNum))
}
