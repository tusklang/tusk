package lang

import "encoding/json"
import "strings"

// #cgo CFLAGS: -std=c99
import "C"

func divide(_num1 string, _num2 string, cli_params map[string]map[string]interface{}) string {

  if returnInit(_num2) == "0" || returnInit(_num2) == "-0" {
    return "undef"
  }
  if returnInit(_num1) == "0" || returnInit(_num1) == "-0" {
    return "0"
  }

  _num1 = addDec(returnInit(_num1))
  _num2 = addDec(returnInit(_num2))

  decO1Index := strings.Index(_num1, ".")
  decO2Index := len(_num2) - strings.Index(_num2, ".") - 1
  combinedIndex := decO1Index + decO2Index

  _num1 = strings.Replace(_num1, ".", "", 1)
  _num2 = strings.Replace(_num2, ".", "", 1)

  var neg = false

  if strings.HasPrefix(_num1, "-") {
    _num1 = _num1[1:]
    neg = !neg
  }

  if strings.HasPrefix(_num2, "-") {
    _num2 = _num2[1:]
    neg = !neg
  }

  for ;len(_num2) > len(_num1); {
    _num1+=strings.Repeat("0", 20)
  }

  _num1+=strings.Repeat("0", cli_params["Calc"]["PREC"].(int))

  curVal := ""
  final := ""

  for i := 0; i < len(_num1); i++ {

    curVal+=string([]rune(_num1)[i])

    if isLess(curVal, _num2) {
      final+="0"
      continue
    }

    curDivisor := _num2
    curQ := "1"

    for ;isLess(add(curDivisor, _num2, cli_params), curVal) || returnInit(add(curDivisor, _num2, cli_params)) == returnInit(curVal); {
      curDivisor = add(curDivisor, _num2, cli_params)
      curQ = add(curQ, "1", cli_params)
    }

    curVal = subtract(curVal, curDivisor, cli_params)

    final+=curQ
  }

  for ;len(final) < combinedIndex; final = "0" + final {}

  final = final[:combinedIndex] + "." + final[combinedIndex:]

  if neg {
    final = "-" + final
  }

  return final
}

//export Division
func Division(_num1P *C.char, _num2P *C.char, cli_paramsP *C.char) *C.char {

  _num1 := C.GoString(_num1P)
  _num2 := C.GoString(_num2P)
  cli_params_str := C.GoString(cli_paramsP)

  var cli_params map[string]map[string]interface{}

  _ = json.Unmarshal([]byte(cli_params_str), &cli_params)

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
      numRet := returnInit(divide(_num1P_.ExpStr[0], _num2P_.ExpStr[0], cli_params))

      finalRet = Action{ "number", "", []string{ numRet }, []Action{}, []string{}, [][]Action{}, []Condition{}, 39, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), false }
    case TypeOperations{ "string", "number" }: //detect case "string" / "num"
      //get string length
      str_ := _num1P_.ExpStr[0]

      var length string

      for length = "1"; str_ != ""; length = add(length, "1", cli_params) {
        str_ = str_[1:]
      }
      ////////////////////

      subtracted := subtract(length, "2", cli_params)

      if isLess(subtracted, _num2P_.ExpStr[0]) || returnInit(subtracted) == _num2P_.ExpStr[0] {
        finalRet = Action{ "falsey", "", []string{ "undef" }, []Action{}, []string{}, [][]Action{}, []Condition{}, 41, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), false }
      } else {

        var _str []string

        for i := subtract(subtract(length, _num2P_.ExpStr[0], cli_params), "2", cli_params); isLess("0", i); i = subtract(i, "1", cli_params) {
          _str = append([]string{ getIndex(_num1P_.ExpStr[0], i) }, _str...)
        }

        str := strings.Join(_str, "")

        finalRet = Action{ "string", "", []string{ str }, []Action{}, []string{}, [][]Action{}, []Condition{}, 38, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), false }
      }
    case TypeOperations{ "number", "string" }: //detect case "string" / "num"

      //get string length
      str_ := _num2P_.ExpStr[0]

      var length string

      for length = "1"; str_ != ""; length = add(length, "1", cli_params) {
        str_ = str_[1:]
      }
      ////////////////////

      subtracted := subtract(length, "2", cli_params)

      if isLess(subtracted, _num1P_.ExpStr[0]) || returnInit(subtracted) == _num1P_.ExpStr[0] {
        finalRet = Action{ "falsey", "", []string{ "undef" }, []Action{}, []string{}, [][]Action{}, []Condition{}, 41, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), false }
      } else {

        var cur string
        str := _num2P_.ExpStr[0]

        for i := add(_num1P_.ExpStr[0], "2", cli_params); isLess(i, length); i = add(i, "1", cli_params) {
          cur+=getIndex(str, i)
        }

        finalRet = Action{ "string", "", []string{ cur }, []Action{}, []string{}, [][]Action{}, []Condition{}, 38, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), false }
      }
    default: finalRet = Action{ "falsey", "", []string{ "undef" }, []Action{}, []string{}, [][]Action{}, []Condition{}, 41, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), false }
  }

  reCalc(&finalRet)

  jsonNum, _ := json.Marshal(finalRet)

  return C.CString(string(jsonNum))
}
