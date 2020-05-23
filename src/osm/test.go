package main

import (
    "fmt"
    "unsafe"
)

/*
extern void go_callback_int(void*, int);
static inline void CallMyFunction(void* pfoo) {
    go_callback_int(pfoo, 5);
}
*/
import "C"

//export go_callback_int
func go_callback_int(pfoo unsafe.Pointer, p1 C.int) {
    foo := (*Test)(pfoo)
    foo.cb(p1)
}

type Test struct {
}

func (s *Test) cb(x C.int) {
    fmt.Println("callback with", x)
}

func main() {
    data := &Test{}
    C.CallMyFunction(unsafe.Pointer(&data))
}
