package main

// #cgo CFLAGS: -std=c99
import "C"

//export GetOp
func GetOp(_val *C.char) *C.char {
  val := C.GoString(_val)

  switch (val) {
    case "add":
      return C.CString("+")
    case "subtract":
      return C.CString("-")
    case "multiply":
      return C.CString("*")
    case "divide":
      return C.CString("/")
    case "exponentiate":
      return C.CString("^")
    case "modulo":
      return C.CString("%")
  }

  //if no operation was detected, return ?OP?
  return C.CString("?OP?")
}
