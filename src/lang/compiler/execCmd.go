package compiler

import "os/exec"
import "strings"

// #cgo CFLAGS: -std=c99
import "C"

//export ExecCmd
func ExecCmd(cmd *C.char, stdin *C.char, dir *C.char) *C.char {

  command := exec.Command(C.GoString(cmd))

  command.Dir = C.GoString(dir)
  command.Stdin = strings.NewReader(C.GoString(stdin))

  out, _ := command.CombinedOutput()

  return C.CString(string(out))
}
