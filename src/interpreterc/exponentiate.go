package main

import "math/big"
import "encoding/json"
import "fmt"
import "strconv"
import "github.com/ALTree/bigfloat"

// #cgo CFLAGS: -std=c99
import "C"

func exponentiate(num1 string, num2 string, calc_params paramCalcOpts, line int) string {

  num1big, _ := new(big.Float).SetPrec(PREC).SetString(num1)
  num2big, _ := new(big.Float).SetString(num2)

  var doNeg = 1
  if num1big.Cmp(big.NewFloat(0)) == -1 {

    inted, _ := new(big.Int).SetString(num2big.String(), 10)
    doNeg, _ = strconv.Atoi(new(big.Int).Quo(inted, big.NewInt(2)).String())

    num1big = new(big.Float).Mul(big.NewFloat(-1), num1big)
  }

  powwed := bigfloat.Pow(num1big, num2big)

  //if it is a negative power to 1/num1^num2
  if doNeg == 0 {
    powwed = new(big.Float).Mul(big.NewFloat(-1), powwed)
  }

  return returnInit(fmt.Sprintf("%f", powwed))
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
