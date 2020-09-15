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
var asyncfuncs = make(map[uint64]func() *KaType)
var allthreads = make(map[uint64]KaThread)

func newthread(cb func() *KaType) *KaThread {

	asyncfuncs[curptr] = cb

	created := cthread(C.newThread(C.ulonglong(curptr)))
	var kathread KaThread
	kathread.thread = &created
	kathread.ptr = curptr

	allthreads[curptr] = kathread
	curptr++

	return &kathread
}

func jointhread(thread KaThread) *KaType {
	var joined = C.waitfor((C.struct_Thread)(*thread.thread))
	return *(**KaType)(joined)
}

func thread_dealloc(thread KaThread) {
	C.freeThread(C.struct_Thread(*thread.thread))
	delete(asyncfuncs, thread.ptr) //remove from the asyncfuncs
}

//export CallGoCB
func CallGoCB(ptr C.ulonglong, output *unsafe.Pointer) {

	if _, e := asyncfuncs[uint64(ptr)]; e { //if it does not exist, it became corrupt, so ignore it
		ret := asyncfuncs[uint64(ptr)]() //call the func based on the unsafe pointer address
		*output = unsafe.Pointer(&ret)
	}
}
