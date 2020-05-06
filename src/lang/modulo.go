package lang

import "encoding/json"
import "strings"

// #cgo CFLAGS: -std=c99
import "C"

func modulo(_num1 string, _num2 string, cli_params map[string]map[string]interface{}) string {

  if _num2 == "0"  {
    return "undef"
  }

  if returnInit(_num1) == "0" {
    return "0"
  }

  divved_ := addDec(divide(_num1, _num2, cli_params))

  divved := divved_[:strings.Index(divved_, ".")]

  mult := multiply(divved, _num2, cli_params)
  remainder := subtract(_num1, mult, cli_params)

  return remainder
}

//export Modulo
func Modulo(_num1P *C.char, _num2P *C.char, cli_paramsP *C.char) *C.char {

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

    num % num = num
    default = falsey
  */

  nums := TypeOperations{ _num1P_.Type, _num2P_.Type }

  var finalRet Action

  switch nums {
    case TypeOperations{ "number", "number" }: //detect case "num" % "num"
      val := modulo(_num1P_.ExpStr[0], _num2P_.ExpStr[0], cli_params)

      finalRet = Action{ "number", "", []string{ val }, []Action{}, []string{}, [][]Action{}, []Condition{}, 39, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), false }
    default: finalRet = Action{ "falsey", "", []string{ "undef" }, []Action{}, []string{}, [][]Action{}, []Condition{}, 41, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), false }
  }

  reCalc(&finalRet)

  jsonNum, _ := json.Marshal(finalRet)

  return C.CString(string(jsonNum))
}
