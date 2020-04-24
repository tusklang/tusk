package main

import "encoding/json"
import "math/big"
import "fmt"

// #cgo CFLAGS: -std=c99
import "C"

func multiply(num1 string, num2 string, calc_params paramCalcOpts, line int) string {
  calc := new(big.Float)

  num1big, _ := new(big.Float).SetPrec(PREC).SetString(num1)
  num2big, _ := new(big.Float).SetPrec(PREC).SetString(num2)

  return returnInit(fmt.Sprintf("%f", calc.Mul(num1big, num2big)))
}

//export Multiply
func Multiply(_num1P *C.char, _num2P *C.char, calc_paramsP *C.char, line_ C.int) *C.char {

  _num1 := returnInit(C.GoString(_num1P))
  _num2 := returnInit(C.GoString(_num2P))
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

    num * num = num
    string * num = string
    array * array = array
    default = falsey
  */

  nums := TypeOperations{ _num1P_.Type, _num2P_.Type }

  var finalRet Action

  switch nums {
    case TypeOperations{ "number", "number" }: //detect case "num" * "num"
      val := multiply(_num1P_.ExpStr[0], _num2P_.ExpStr[0], calc_params, line)

      finalRet = Action{ "number", "", []string{ val }, []Action{}, []string{}, []Action{}, []Condition{}, 39, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), false }
    case TypeOperations{ "string", "number" }: //detect case "string" * "num"
      str := _num1P_.ExpStr[0]
      str = str[1:len(str) - 1]

      var val string

      for i := "0"; isLess(i, _num2P_.ExpStr[0]); i = add(i, "1", calc_params, line) {
        val+=str
      }

      val = val
      finalRet = Action{ "string", "", []string{ val }, []Action{}, []string{}, []Action{}, []Condition{}, 38, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), false }
    case TypeOperations{ "number", "string" }: //detect case "string" * "num"
      str := _num2P_.ExpStr[0]
      str = str[1:len(str) - 1]

      var val string

      for i := "0"; isLess(i, _num1P_.ExpStr[0]); i = add(i, "1", calc_params, line) {
        val+=str
      }

      val = val
      finalRet = Action{ "string", "", []string{ val }, []Action{}, []string{}, []Action{}, []Condition{}, 38, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), false }
    case TypeOperations{ "array", "array" }: //detect case "array" * "array"

      //get length of first array

      length := "0"

      for _, _ = range _num1P_.Hash_Values {
        length = add(length, "1", calc_params, line)
      }
      ///////////////////////////

      nMap := make(map[string][]Action)

      for k, v := range _num2P_.Hash_Values {
        nMap[add(length, k, calc_params, line)] = v
      }

      //merge the two maps

      finalMap := make(map[string][]Action)

      for k, v := range _num1P_.Hash_Values {
        finalMap[k] = v
      }

      for k, v := range nMap {
        finalMap[k] = v
      }

      finalRet = Action{ "array", "", []string{}, []Action{}, []string{}, []Action{}, []Condition{}, 24, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, finalMap, false }
    default: finalRet = Action{ "falsey", "", []string{ "undefined" }, []Action{}, []string{}, []Action{}, []Condition{}, 41, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), false }
  }

  reCalc(&finalRet)

  jsonNum, _ := json.Marshal(finalRet)

  return C.CString(string(jsonNum))
}
