package compile

import "regexp"

// #cgo CFLAGS: -std=c99
import "C"

//export IsAbsolute
func IsAbsolute(_dir *C.char) C.int {

  dir := C.GoString(_dir)

  match, _ := regexp.MatchString("^[a-zA-Z]:", dir)

  if match {
    return 1
  } else {
    return 0  
  }
}
