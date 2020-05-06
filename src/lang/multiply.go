package lang

import "encoding/json"
import "strings"

// #cgo CFLAGS: -std=c99
import "C"

func multiply(_num1 string, _num2 string, cli_params map[string]map[string]interface{}) string {

  if returnInit(_num1) == "0" || returnInit(_num2) == "0" {
    return "0"
  }

  decIndex := 0

  if strings.Contains(_num1, ".") {
    decIndex+=len(strings.Replace(_num1, ".", "", 1)) - strings.Index(_num1, ".")
  }
  if strings.Contains(_num2, ".") {
    decIndex+=len(strings.Replace(_num2, ".", "", 1)) - strings.Index(_num2, ".")
  }

  _num1 = strings.Replace(_num1, ".", "", 1)
  _num2 = strings.Replace(_num2, ".", "", 1)

  if isLess(_num1, _num2) {
    _num1, _num2 = _num2, _num1
  }

  neg := false

  if strings.HasPrefix(_num1, "-") {
    _num1 = _num1[1:]
    neg = !neg
  }
  if strings.HasPrefix(_num2, "-") {
    _num2= _num2[1:]
    neg = !neg
  }

  var nNum string

  if len(_num1) < len(_num2) {
    _num1, _num2 = _num2, _num1
  }

  if len(_num1) >= cli_params["Calc"]["LONG_MULT_THRESH"].(int) && len(_num2) >= cli_params["Calc"]["LONG_MULT_THRESH"].(int) {

    final := "0"

    for i := len(_num2) - 1; i >= 0; i-- {
      mult := multiply(_num1, _num2[i:i + 1], cli_params) + strings.Repeat("0", len(_num2) - i - 1)
      final = add(final, mult, cli_params)
    }

    nNum = final
  } else {

    nNum = "0"

    for ;returnInit(_num2) != "0"; {
      nNum = add(nNum, _num1, cli_params)
      _num2 = subtract(_num2, "1", cli_params)
    }
  }

  if decIndex > 0 {

    nNum = Reverse(nNum)

    nNum = nNum[:decIndex] + "." + nNum[decIndex:]

    nNum = Reverse(nNum)
  }

  if neg == true {
    nNum = "-" + nNum
  }

  return returnInit(nNum)
}

//export Multiply
func Multiply(_num1P *C.char, _num2P *C.char, cli_paramsP *C.char) *C.char {

  _num1 := returnInit(C.GoString(_num1P))
  _num2 := returnInit(C.GoString(_num2P))
  cli_params_str := C.GoString(cli_paramsP)

  var cli_params map[string]map[string]interface{}

  _ = json.Unmarshal([]byte(cli_params_str), &cli_params)

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
      val := multiply(_num1P_.ExpStr[0], _num2P_.ExpStr[0], cli_params)

      finalRet = Action{ "number", "", []string{ val }, []Action{}, []string{}, [][]Action{}, []Condition{}, 39, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), false }
    case TypeOperations{ "string", "number" }: //detect case "string" * "num"
      str := _num1P_.ExpStr[0]
      str = str[1:len(str) - 1]

      var val string

      for i := "0"; isLess(i, _num2P_.ExpStr[0]); i = add(i, "1", cli_params) {
        val+=str
      }

      val = val
      finalRet = Action{ "string", "", []string{ val }, []Action{}, []string{}, [][]Action{}, []Condition{}, 38, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), false }
    case TypeOperations{ "number", "string" }: //detect case "string" * "num"
      str := _num2P_.ExpStr[0]
      str = str[1:len(str) - 1]

      var val string

      for i := "0"; isLess(i, _num1P_.ExpStr[0]); i = add(i, "1", cli_params) {
        val+=str
      }

      val = val
      finalRet = Action{ "string", "", []string{ val }, []Action{}, []string{}, [][]Action{}, []Condition{}, 38, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), false }
    case TypeOperations{ "array", "array" }: //detect case "array" * "array"

      //get length of first array

      length := "0"

      for _, _ = range _num1P_.Hash_Values {
        length = add(length, "1", cli_params)
      }
      ///////////////////////////

      nMap := make(map[string][]Action)

      for k, v := range _num2P_.Hash_Values {
        nMap[add(length, k, cli_params)] = v
      }

      //merge the two maps

      finalMap := make(map[string][]Action)

      for k, v := range _num1P_.Hash_Values {
        finalMap[k] = v
      }

      for k, v := range nMap {
        finalMap[k] = v
      }

      finalRet = Action{ "array", "", []string{}, []Action{}, []string{}, [][]Action{}, []Condition{}, 24, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, finalMap, false }
    default: finalRet = Action{ "falsey", "", []string{ "undef" }, []Action{}, []string{}, [][]Action{}, []Condition{}, 41, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), false }
  }

  reCalc(&finalRet)

  jsonNum, _ := json.Marshal(finalRet)

  return C.CString(string(jsonNum))
}
