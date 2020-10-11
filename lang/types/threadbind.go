package types

//This is a binding file
//So the code is very messy
//(because C is very messy)
//But I love C

import (
	"unsafe"
)

//#include "thread.h"
import "C"

type cthread C.struct_Thread

var curptr uint64 = 0
var asyncfuncs = make(map[uint64]func() (*TuskType, *TuskError))
var allthreads = make(map[uint64]TuskThread)

func newthread(cb func() (*TuskType, *TuskError)) *TuskThread {

	asyncfuncs[curptr] = cb

	created := cthread(C.newThread(C.ulonglong(curptr)))
	var kathread TuskThread
	kathread.thread = &created
	kathread.ptr = curptr

	allthreads[curptr] = kathread
	curptr++

	return &kathread
}

func jointhread(thread TuskThread) *TuskType {
	var joined = C.waitfor((C.struct_Thread)(*thread.thread))
	return *(**TuskType)(joined)
}

func thread_dealloc(thread TuskThread) {
	C.freeThread(C.struct_Thread(*thread.thread))
	delete(asyncfuncs, thread.ptr) //remove from the asyncfuncs
}

//export CallGoCB
func CallGoCB(ptr C.ulonglong, output *unsafe.Pointer) {

	if _, e := asyncfuncs[uint64(ptr)]; e { //if it does not exist, it became corrupt, so ignore it
		ret, e := asyncfuncs[uint64(ptr)]() //call the func based on the unsafe pointer address
		if e != nil {
			e.Print()
		}
		*output = unsafe.Pointer(&ret)
	}
}
