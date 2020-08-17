package types

import "unsafe"

//#include "thread.h"
import "C"

type cthread C.struct_Thread

func newthread(cb unsafe.Pointer) unsafe.Pointer {
	created := C.newThread(cb)
	return unsafe.Pointer(&created)
}

func jointhread(thread OmmThread) *OmmType {
	var joined = C.waitfor((C.struct_Thread)(*thread.thread))
	return (*OmmType)(joined)
}

//export CallGoCB
func CallGoCB(_fn unsafe.Pointer) unsafe.Pointer {
	fn := *(*(func() Returner))(_fn)
	return unsafe.Pointer(&fn)
}
