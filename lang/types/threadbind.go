package types

//this is a binding file
//so the code is very messy
//(because c is very messy)
//but i love c

import (
	"unsafe"
)

//#include "thread.h"
import "C"

type cthread C.struct_Thread

var curptr uint64 = 0
var asyncfuncs = make(map[uint64]func() *OmmType)
var allthreads = make(map[uint64]OmmThread)

func newthread(cb func() *OmmType) *OmmThread {

	asyncfuncs[curptr] = cb

	created := cthread(C.newThread(C.ulonglong(curptr)))
	var ommthread OmmThread
	ommthread.thread = &created
	ommthread.ptr = curptr

	allthreads[curptr] = ommthread
	curptr++

	return &ommthread
}

func jointhread(thread OmmThread) *OmmType {
	var joined = C.waitfor((C.struct_Thread)(*thread.thread))
	return *(**OmmType)(joined)
}

func thread_dealloc(thread OmmThread) {
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
