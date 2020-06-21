package lang

// #cgo CFLAGS: -std=c99
// #include "bind.h"
import "C"

func colorprint(msg string, color int) {
  C.colorprint(C.CString(msg), C.int(color))
}
