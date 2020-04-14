package main

import "strings"
import "strconv"
import "math"
import "encoding/json"

// #cgo CFLAGS: -std=c99
import "C"

func add(__num1 string, __num2 string, calc_params paramCalcOpts, line int) string {

  _num1 := __num1
  _num2 := __num2

  num1_, num2_ := initAdd(_num1, _num2)

  decPlace := getDec(num1_)

  var final = ""

  if !strings.HasPrefix(num1_, "-") && !strings.HasPrefix(num2_, "-") {

    var carry = 0

    num1 := Chunk(strings.ReplaceAll(num1_, ".", ""), 9)
    num2 := Chunk(strings.ReplaceAll(num2_, ".", ""), 9)

    for i := len(num1) - 1; i >= 0; i-- {
      n1, _ := strconv.ParseInt(num1[i], 10, 64)
      n2, _ := strconv.ParseInt(num2[i], 10, 64)

      sum := strconv.Itoa(int(n1 + n2 + int64(carry)))

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
  } else if strings.HasPrefix(num1_, "-") && !strings.HasPrefix(num2_, "-") {
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

  } else if !strings.HasPrefix(num1_, "-") && strings.HasPrefix(num2_, "-") {
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

  } else if strings.HasPrefix(num1_, "-") && strings.HasPrefix(num2_, "-") {

    num1_ = num1_[1:]
    num2_ = num2_[1:]

    num1 := Chunk(strings.ReplaceAll(num1_, ".", ""), 9)
    num2 := Chunk(strings.ReplaceAll(num2_, ".", ""), 9)

    var carry = 0

    for i := len(num1) - 1; i >= 0; i-- {
      n1, _ := strconv.ParseUint(string(num1[i]), 10, 64)
      n2, _ := strconv.ParseUint(string(num2[i]), 10, 64)

      sum := strconv.Itoa(int(n1 + n2 + uint64(carry)))

      carry = 0
      if len(sum) > len(string(num1[i])) {
        sum = sum[1:]
        carry = 1
      }

      nL := math.Max(float64(len(string(num1[i]))), float64(len(string(num2[i]))))

      for ;float64(len(sum)) < nL; {
        sum = "0" + sum
      }

      final = sum + final;
    }

    final = "-" + final[:decPlace] + "." + final[decPlace:]
  }

  return returnInit(final)
}

//export Add
func Add(_num1P *C.char, _num2P *C.char, calc_paramsP *C.char, line_ C.int) *C.char {

  __num1 := C.GoString(_num1P)
  __num2 := C.GoString(_num2P)
  calc_params_str := C.GoString(calc_paramsP)

  line := int(line_)

  _ = line

  var calc_params paramCalcOpts

  _ = json.Unmarshal([]byte(calc_params_str), &calc_params)

  var _num1P_ Action
  var _num2P_ Action

  _ = json.Unmarshal([]byte(__num1), &_num1P_)
  _ = json.Unmarshal([]byte(__num2), &_num2P_)

  nums := TypeOperations{ _num1P_.Type, _num2P_.Type }

  /* TABLE OF TYPES:

    string + (* - array - none - hash) = string
    array + (* - none) = array
    none + * = none
    hash + (* - hash) = none
    type + (* - hash - none) = type
    num + num = num
    hash + hash = hash
    boolean + boolean = boolean
    num + boolean = num
  */

  var finalRet Action

  switch nums {
    case TypeOperations{ "number", "number" }: { //detect case "num" + "num"

      numRet := add(_num1P_.ExpStr[0], _num2P_.ExpStr[0], calc_params, line)

      finalRet = Action{ "number", "", []string{ numRet }, []Action{}, []string{}, []Action{}, []Condition{}, 39, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), []Action{ Action{"number", "", []string{ numRet }, []Action{}, []string{}, []Action{}, []Condition{}, 39, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), []Action{}} } }
    }
    case TypeOperations{ "boolean", "boolean" }: { //detect case "boolean" + "boolean"

      val1 := _num1P_.ExpStr[0]
      val2 := _num2P_.ExpStr[0]

      boolSwitch := BoolSwitch{ val1, val2 }

      var final_ string

      switch (boolSwitch) {
        case BoolSwitch{ "true", "true" }:
          final_ = "1"
        case BoolSwitch{ "true", "false"}:
          final_ = "1"
        case BoolSwitch{ "false", "true" }:
          final_ = "1"
        case BoolSwitch{ "false", "false" }:
          final_ = "0"
      }

      finalRet = Action{ "boolean", "", []string{ final_ }, []Action{}, []string{}, []Action{}, []Condition{}, 40, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), []Action{ Action{"boolean", "", []string{ final_ }, []Action{}, []string{}, []Action{}, []Condition{}, 40, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), []Action{}} } }
    }
    case TypeOperations{ "number", "boolean" }: { //detect case "num" + "boolean"

      val1 := _num1P_.ExpStr[0]
      val2 := _num2P_.ExpStr[0]

      var final_ string

      if val2 == "true" {
        final_ = add(val1, "1", calc_params, line)
      } else {
        final_ = val1
      }

      final_ = returnInit(final_)

      finalRet = Action{ "number", "", []string{ final_ }, []Action{}, []string{}, []Action{}, []Condition{}, 39, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), []Action{ Action{"number", "", []string{ final_ }, []Action{}, []string{}, []Action{}, []Condition{}, 39, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), []Action{}} } }
    }
    case TypeOperations{ "boolean", "number" }: { //detect case "num" + "boolean"

      val1 := _num1P_.ExpStr[0]
      val2 := _num2P_.ExpStr[0]

      var final_ string

      if val1 == "true" {
        final_ = add(val2, "1", calc_params, line)
      } else {
        final_ = val2
      }

      final_ = returnInit(final_)

      finalRet = Action{ "number", "", []string{ final_ }, []Action{}, []string{}, []Action{}, []Condition{}, 39, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), []Action{ Action{"number", "", []string{ final_ }, []Action{}, []string{}, []Action{}, []Condition{}, 39, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), []Action{}} } }
    }
    default:

      if (nums.First == "string" && nums.Second != "array" && nums.Second != "none" && nums.Second != "hash") || (nums.First != "array" && nums.First != "none" && nums.First != "hash" && nums.Second == "string") { //detect case "string" + (* - "array" - "none" - "hash") = "string"
        val1 := _num1P_.ExpStr[0]
        val2 := _num2P_.ExpStr[0]

        if strings.HasPrefix(val1, "'") || strings.HasPrefix(val1, "\"") || strings.HasPrefix(val1, "`") {
          val1 = val1[1:len(val1) - 1]
        }
        if strings.HasPrefix(val2, "'") || strings.HasPrefix(val2, "\"") || strings.HasPrefix(val2, "`") {
          val2 = val2[1:len(val2) - 1]
        }

        final := "'" + val1 + val2 + "'"

        finalRet = Action{ "string", "", []string{ final }, []Action{}, []string{}, []Action{}, []Condition{}, 38, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), []Action{ Action{"string", "", []string{ final }, []Action{}, []string{}, []Action{}, []Condition{}, 38, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), []Action{}} } }
      } else if (nums.First == "array" && nums.Second != "none") || (nums.First != "none" && nums.Second == "array") { //detect case "array" + (* - "none") = "array"

      } else if nums.First == "none" || nums.Second == "none" { //detect case "none" + * = "none"

      } else if (nums.First == "hash" && nums.Second != "none") || (nums.First != "none" && nums.Second == "hash") { //detect case "hash" + (* - "hash") = "none"
        val1 := _num1P_.Hash_Values
        val2 := _num2P_.Hash_Values

        var final = make(map[string][]Action)

        for k, v := range val1 {
          final[k] = v
        }

        for k, v := range val2 {
          final[k] = v
        }

        finalRet = Action{ "hash", "", []string{}, []Action{}, []string{}, []Action{}, []Condition{}, 22, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, final, []Action{ Action{"hash", "", []string{}, []Action{}, []string{}, []Action{}, []Condition{}, 22, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, final, []Action{}} } }
      } else if (nums.First == "type" && nums.Second != "hash" && nums.Second != "none") || (nums.First != "hash" && nums.First != "none" && nums.Second == "type") { //detect case "type" + (* - "hash" - "none") = "type"

      } else {
        finalRet = Action{ "falsey", "", []string{ "undefined" }, []Action{}, []string{}, []Action{}, []Condition{}, 41, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), []Action{ Action{"falsey", "", []string{ "undefined" }, []Action{}, []string{}, []Action{}, []Condition{}, 41, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), []Action{}} } }
      }
  }

  jsonNum, _ := json.Marshal(finalRet)

  return C.CString(string(jsonNum))
}
