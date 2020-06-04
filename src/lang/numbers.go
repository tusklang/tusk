package lang

import "strings"
import "strconv"

// #cgo CFLAGS: -std=c99
import "C"

type TypeOperations struct {
  First    string
  Second   string
}

type BoolSwitch struct {
  First    string
  Second   string
}

const PREC = 1000

//export ReturnInitC
func ReturnInitC(str *C.char) *C.char {
  return C.CString( returnInit( C.GoString(str) ) )
}

//export IsLessC
func IsLessC(n1, n2 *C.char) C.int {

  if isLess( C.GoString(n1), C.GoString(n2) ) {
    return 1
  }

  return 0
}

func returnInit(str string) string {

  if strings.HasPrefix(str, "+") {
    str = str[1:]
  }

  if getType(str) == "string" {
    str = strconv.Itoa(len(str) - 2)
  }

  if str == "true" {
    str = "1"
  }
  if str == "false" {
    str = "0"
  }

  if getType(str) == "falsey" {
    str = "-1"
  }

  for ;strings.HasPrefix(str, "0"); {
    str = str[1:]
  }

  for ;strings.HasSuffix(str, "0") && strings.Contains(str, "."); {
    str = strings.TrimSuffix(str, "0")
  }

  if strings.HasSuffix(str, ".") {
    str = strings.TrimSuffix(str, ".")
  }

  if strings.HasPrefix(str, "-") {
    str = str[1:]

    for ;strings.HasPrefix(str, "0"); {
      str = str[1:]
    }

    str = "-" + str
  }

  if len(str) == 0 {
    str = "0"
  }

  if str == "-0" {
    str = "0"
  }

  if strings.HasPrefix(str, ".") {
    str = "0" + str
  }

  return str
}

func splitEvery(str string, by uint64) []string {
  nString := []string{}

  for i := 0; i < len(str); i++ {
    if uint64(i) % by == 0 {
      nString = append(nString, "")
    }

    nString[len(nString) - 1]+=string(i)
  }

  return nString
}

func initAdd(num1 string, num2 string) (string, string) {
  if !strings.Contains(num1, ".") {
    num1+=".0"
  }
  if !strings.Contains(num2, ".") {
    num2+=".0"
  }

  num1Neg := false
  num2Neg := false

  if strings.HasPrefix(num1, "-") {
    num1 = num1[1:]
    num1Neg = true
  }
  if strings.HasPrefix(num2, "-") {
    num2 = num2[1:]
    num2Neg = true
  }

  num1_ := strings.Split(num1, ".")
  num2_ := strings.Split(num2, ".")

  for ;len(num1_[0]) != len(num2_[0]); {
    if len(num1_[0]) < len(num2_[0]) {
      num1_[0] = "0" + num1_[0]
    } else {
      num2_[0] = "0" + num2_[0]
    }
  }

  for ;len(num1_[1]) != len(num2_[1]); {
    if len(num1_[1]) < len(num2_[1]) {
      num1_[1] = num1_[1] + "0"
    } else {
      num2_[1] = num2_[1] + "0"
    }
  }

  num1 = num1_[0] + "." + num1_[1]
  num2 = num2_[0] + "." + num2_[1]

  num1 = "00" + num1
  num2 = "00" + num2

  if num1Neg {
    num1 = "-" + num1
  }

  if num2Neg {
    num2 = "-" + num2
  }

  return num1, num2
}

func addDec(num string) string {

  if num == "undef" {
    return "undef";
  }

  if !strings.Contains(num, ".") {
    return num + ".0"
  }

  return num
}

func abs(num []string) []string {

  //ensure that the given number will not change
  num = append(num, "-")

  for k, v := range num {
    if v[0] == '-' {
      num[k] = num[k][1:]
    }
  }

  return num
}

func isLess(num1 string, num2 string) bool {

  if num1 == "undef" || num2 == "undef" || num1 == "NaN" || num2 == "NaN" {
    return false
  }

  if getType(num1) == "string" {
    num1 = strconv.Itoa(len(num1) - 2)
  }

  if getType(num2) == "string" {
    num2 = strconv.Itoa(len(num2) - 2)
  }

  if num1 == "true" {
    num1 = "1"
  }
  if num1 == "false" {
    num1 = "0"
  }

  if num2 == "true" {
    num2 = "1"
  }
  if num2 == "false" {
    num2 = "0"
  }

  if getType(num1) == "falsey" {
    num1 = "-1"
  }

  if getType(num2) == "falsey" {
    num2 = "-1"
  }

  num1 = returnInit(num1)
  num2 = returnInit(num2)

  if num1 == num2 {

    return false
  } else if strings.HasPrefix(num1, "-") && !strings.HasPrefix(num2, "-") {

    return true
  } else if !strings.HasPrefix(num1, "-") && strings.HasPrefix(num2, "-") {

    return false
  } else {

    num1_, num2_ := initAdd(num1, num2)

    for i := 0; i < len(num1_); i++ {
      n1, _ := strconv.ParseUint(string(num1_[i]), 10, 64)
      n2, _ := strconv.ParseUint(string(num2_[i]), 10, 64)

      if n1 < n2 {
        return true
      } else if n1 > n2 {
        break
      }
    }
  }

  return false
}
func getDec(num string) int {

  if strings.HasPrefix(num, "-") {
    return strings.Index(num[1:], ".")
  }

  return strings.Index(num, ".")
}

//function to re calculate the hash values (because operations can manipulate expStr but not hash values)
func reCalc(val *Action) {

  switch ((*val).Type) {
    case "string":
      expstr := val.ExpStr[0]

      for i := 0; i < len(expstr); i++ {
        val.Hash_Values[strconv.Itoa(i)] = []Action{ Action{ "string", "", []string{ string(expstr[i]) }, []Action{}, []string{}, [][]Action{}, []Condition{}, 38, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), false, "private", []SubCaller{}, []int{}, []int{} } }
      }
  }
}

func getIndex(val, index string) string {

  for i := "1"; isLess(i, index); i = C.GoString(AddC(C.CString(i), C.CString("1"), C.CString("{}"))) {
    val = val[1:]
  }

  return string(val[0])
}

//function to chunk a string into pieces
func Chunk(val string, by int) []string {

  //val to return
  nVal := []string{}
  last := 0

  for i := 0; i < len(val); i++ {

    //if it is on the loop then append the last to the current index to nVar
    if i % by == 0 {
      nVal = append(nVal, val[last:i])
      last = i
    }
  }

  //if val's length does not equal the last value append the remaining bits of the string to nVal
  if (last != len(val)) {
    nVal = append(nVal, val[last:len(val)])
  }

  return nVal;
}
