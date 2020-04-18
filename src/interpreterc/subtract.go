package main

import "strings"
import "strconv"
import "math"
import "encoding/json"

// #cgo CFLAGS: -std=c99
import "C"

//export SubtractStrings
func SubtractStrings(num1, num2, calc_params *C.char, line C.int) *C.char {

  var cp paramCalcOpts

  _ = json.Unmarshal([]byte(C.GoString(calc_params)), &cp)

  sum := subtract(C.GoString(num1), C.GoString(num2), cp, int(line))

  return C.CString(sum)
}

func subtract(_num1 string, _num2 string, calc_params paramCalcOpts, line int) string {
  num1_, num2_ := initAdd(_num1, _num2)

  decPlace := getDec(num1_)

  var final = ""

  if !strings.HasPrefix(num1_, "-") && !strings.HasPrefix(num2_, "-") {

    _num1_ := addDec(strings.Replace(num1_, "-", "0", 1))
    _num2_ := addDec(strings.Replace(num2_, "-", "0", 1))

    num1, num2 := initAdd(_num1_, _num2_)

    decPlace = getDec(num1)

    num1 = strings.Replace(num1, ".", "", 1)
    num2 = strings.Replace(num2, ".", "", 1)

    if !isLess(num1, num2) {

      for i := len(num1) - 1; i >= 0; i-- {
        n1, _ := strconv.ParseUint(string(num1[i]), 10, 64)
        n2, _ := strconv.ParseUint(string(num2[i]), 10, 64)

        if n1 < n2 {
          n1+=10

          o := i - 1
          o1 := []rune(num1)

          if o1[o] == '.' {
            o--
          }

          for ;o1[o] == '0'; {
            o1[o] = '9'
            o--
          }
          cur, _ := strconv.ParseInt(string(o1[o]), 10, 64)
          o1[o] = []rune(strconv.Itoa(int(cur - 1)))[0]

          num1 = string(o1)
        }

        sum := strconv.Itoa(int(n1 - n2))

        final = sum + final
      }

      final = final[:decPlace] + "." + final[decPlace:]
    } else {
      switchOpts := []string{ num1, num2 }

      num1 = switchOpts[1]
      num2 = switchOpts[0]

      for i := len(num1) - 1; i >= 0; i-- {
        n1, _ := strconv.ParseUint(string(num1[i]), 10, 64)
        n2, _ := strconv.ParseUint(string(num2[i]), 10, 64)

        if n1 < n2 {
          n1+=10

          o := i - 1
          o1 := []rune(num1)

          if o1[o] == '.' {
            o--
          }

          for ;o1[o] == '0'; {
            o1[o] = '9'
            o--
          }
          cur, _ := strconv.ParseInt(string(o1[o]), 10, 64)
          o1[o] = []rune(strconv.Itoa(int(cur - 1)))[0]

          num1 = string(o1)
        }

        sum := strconv.Itoa(int(n1 - n2))

        final = sum + final
      }

      final = "-" + final[:decPlace] + "." + final[decPlace:]
    }

  } else if strings.HasPrefix(num1_, "-") && !strings.HasPrefix(num2_, "-") {
    var carry = 0

    _num1_ := addDec(strings.Replace(num1_, "-", "0", 1))
    _num2_ := addDec(strings.Replace(num2_, "-", "0", 1))

    num1__, num2__ := initAdd(_num1_, _num2_)

    decPlace = getDec(num1__)

    num1 := Chunk(strings.ReplaceAll(num1__, ".", ""), 9)
    num2 := Chunk(strings.ReplaceAll(num2__, ".", ""), 9)

    for i := len(num1) - 1; i >= 0; i-- {
      n1, _ := strconv.ParseUint(num1[i], 10, 64)
      n2, _ := strconv.ParseUint(num2[i], 10, 64)

      sum := strconv.Itoa(int(n1 + n2 + uint64(carry)))

      carry = 0
      if len(sum) > len(num1[i]) {
        sum = sum[1:]
        carry = 1
      }

      nL := math.Max(float64(len(num1[i])), float64(len(num2[i])))

      for ;float64(len(sum)) < nL; {
        sum = "0" + sum
      }

      final = sum + final;
    }

    final = "-" + final[:decPlace] + "." + final[decPlace:]

  } else if !strings.HasPrefix(num1_, "-") && strings.HasPrefix(num2_, "-") {
    var carry = 0

    num1_ := addDec(strings.Replace(num1_, "-", "0", 1))
    num2_ := addDec(strings.Replace(num2_, "-", "0", 1))

    num1__, num2__ := initAdd(num1_, num2_)

    decPlace = getDec(num1__)

    num1 := Chunk(strings.ReplaceAll(num1__, ".", ""), 9)
    num2 := Chunk(strings.ReplaceAll(num2__, ".", ""), 9)

    for i := len(num1) - 1; i >= 0; i-- {
      n1, _ := strconv.ParseUint(num1[i], 10, 64)
      n2, _ := strconv.ParseUint(num2[i], 10, 64)

      sum := strconv.Itoa(int(n1 + n2 + uint64(carry)))

      carry = 0
      if len(sum) > len(num1[i]) {
        sum = sum[1:]
        carry = 1
      }

      nL := math.Max(float64(len(num1[i])), float64(len(num2[i])))

      for ;float64(len(sum)) < nL; {
        sum = "0" + sum
      }

      final = sum + final;
    }

    final = final[:decPlace] + "." + final[decPlace:]


  } else if strings.HasPrefix(num1_, "-") && strings.HasPrefix(num2_, "-") {

    num1_ = num1_[1:]
    num2_ = num2_[1:]

    num1 := addDec(strings.Replace(num1_, "-", "0", 1))
    num2 := addDec(strings.Replace(num2_, "-", "0", 1))

    num1, num2 = initAdd(num1, num2)

    decPlace = getDec(num1)

    num1 = strings.Replace(num1, ".", "", 1)
    num2 = strings.Replace(num2, ".", "", 1)

    if !isLess(num1, num2) {

      for i := len(num1) - 1; i >= 0; i-- {
        n1, _ := strconv.ParseUint(string(num1[i]), 10, 64)
        n2, _ := strconv.ParseUint(string(num2[i]), 10, 64)

        if n1 < n2 {
          n1+=10

          o := i - 1
          o1 := []rune(num1)

          if o1[o] == '.' {
            o--
          }

          for ;o1[o] == '0'; {
            o1[o] = '9'
            o--
          }
          cur, _ := strconv.ParseInt(string(o1[o]), 10, 64)
          o1[o] = []rune(strconv.Itoa(int(cur - 1)))[0]

          num1 = string(o1)
        }

        sum := strconv.Itoa(int(n1 - n2))

        final = sum + final
      }

      final = "-" + final[:decPlace] + "." + final[decPlace:]
    } else {
      switchOpts := []string{ num1, num2 }

      num1 = switchOpts[1]
      num2 = switchOpts[0]

      for i := len(num1) - 1; i >= 0; i-- {
        n1, _ := strconv.ParseUint(string(num1[i]), 10, 64)
        n2, _ := strconv.ParseUint(string(num2[i]), 10, 64)

        if n1 < n2 {
          n1+=10

          o := i - 1
          o1 := []rune(num1)

          if o1[o] == '.' {
            o--
          }

          for ;o1[o] == '0'; {
            o1[o] = '9'
            o--
          }
          cur, _ := strconv.ParseInt(string(o1[o]), 10, 64)
          o1[o] = []rune(strconv.Itoa(int(cur - 1)))[0]

          num1 = string(o1)
        }

        sum := strconv.Itoa(int(n1 - n2))

        final = sum + final
      }

      final = final[:decPlace] + "." + final[decPlace:]
    }
  }

  return returnInit(final)
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

      finalRet = Action{ "number", "", []string{ val }, []Action{}, []string{}, []Action{}, []Condition{}, 39, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), []Action{ Action{"number", "", []string{ val }, []Action{}, []string{}, []Action{}, []Condition{}, 39, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), []Action{}, false} }, false }
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
        finalRet = Action{ "falsey", "", []string{ "undefined" }, []Action{}, []string{}, []Action{}, []Condition{}, 41, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), []Action{ Action{"falsey", "", []string{ "undefined" }, []Action{}, []string{}, []Action{}, []Condition{}, 41, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), []Action{}, false} }, false }
      } else {

        for i := subtract(length, add(_num2P_.ExpStr[0], "1", calc_params, line), calc_params, line); isLess(i, length); i = add(i, "1", calc_params, line) {
          val+=getIndex(str, i)
        }

        finalRet = Action{ "string", "", []string{ val }, []Action{}, []string{}, []Action{}, []Condition{}, 38, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), []Action{ Action{"string", "", []string{ val }, []Action{}, []string{}, []Action{}, []Condition{}, 38, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), []Action{}, false} }, false }
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
        finalRet = Action{ "falsey", "", []string{ "undefined" }, []Action{}, []string{}, []Action{}, []Condition{}, 41, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), []Action{ Action{"falsey", "", []string{ "undefined" }, []Action{}, []string{}, []Action{}, []Condition{}, 41, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), []Action{}, false} }, false }
      } else {

        for i := "0"; isLess(i, _num1P_.ExpStr[0]) || returnInit(i) == returnInit(_num1P_.ExpStr[0]); i = add(i, "1", calc_params, line) {
          val+=string(str[0])
          str = str[1:]
        }

        finalRet = Action{ "string", "", []string{ val }, []Action{}, []string{}, []Action{}, []Condition{}, 38, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), []Action{ Action{"string", "", []string{ val }, []Action{}, []string{}, []Action{}, []Condition{}, 38, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), []Action{}, false} }, false }
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

      finalRet = Action{ "number", "", []string{ val }, []Action{}, []string{}, []Action{}, []Condition{}, 39, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), []Action{ Action{"number", "", []string{ val }, []Action{}, []string{}, []Action{}, []Condition{}, 39, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), []Action{}, false} }, false }
    case TypeOperations{ "number", "boolean" }: //detect case "number" - "boolean"

      var val2 string

      if _num2P_.ExpStr[0] == "true" {
        val2 = "1"
      } else {
        val2 = "0"
      }

      val := subtract(_num1P_.ExpStr[0], val2, calc_params, line)

      finalRet = Action{ "number", "", []string{ val }, []Action{}, []string{}, []Action{}, []Condition{}, 39, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), []Action{ Action{"number", "", []string{ val }, []Action{}, []string{}, []Action{}, []Condition{}, 39, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), []Action{}, false} }, false }
    case TypeOperations{ "boolean", "number" }: //detect case "number" - "boolean"

      var val1 string

      if _num1P_.ExpStr[0] == "true" {
        val1 = "1"
      } else {
        val1 = "0"
      }

      val := subtract(val1, _num2P_.ExpStr[0], calc_params, line)

      finalRet = Action{ "number", "", []string{ val }, []Action{}, []string{}, []Action{}, []Condition{}, 39, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), []Action{ Action{"number", "", []string{ val }, []Action{}, []string{}, []Action{}, []Condition{}, 39, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), []Action{}, false} }, false }
    default: finalRet = Action{ "falsey", "", []string{ "undefined" }, []Action{}, []string{}, []Action{}, []Condition{}, 41, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), []Action{ Action{"falsey", "", []string{ "undefined" }, []Action{}, []string{}, []Action{}, []Condition{}, 41, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), []Action{}, false} }, false }
  }

  reCalc(&finalRet)

  jsonNum, _ := json.Marshal(finalRet)

  return C.CString(string(jsonNum))
}
