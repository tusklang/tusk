package lang

// #cgo CFLAGS: -std=c99
import "C"

//export GetActNum
func GetActNum(val string) int {

  switch val {
    case "string":
      return 38
    case "boolean":
      return 40
    case "variable":
      return 43
    case "array":
      return 24
    case "hash":
      return 22
    case "statement":
      return 0
    case "falsey":
      return 41
    case "none":
      return 42
    case "process":
      return 10
    case "type":
      return 44
  }

  return 42
}

//export GetActNumC
func GetActNumC(val *C.char) C.int {

  return C.int(GetActNum(C.GoString(val)))
}
