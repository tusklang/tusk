package types

import "fmt"

type OmmThread struct {
	thread      *cthread
	ptr          uint64
}

//export NewThread
func NewThread(cb func() *OmmType) *OmmThread {
	//wrapper for other packages to create OmmThreads
	return newthread(cb)
}

func WaitAllThreads() {
	//wait all of the remaining threads
	for _, v := range allthreads {
		v.Join()
		v.Deallocate()
	}
}

func (ot OmmThread) Join() *OmmType {
	return jointhread(ot)
}

func (ot OmmThread) Format() string {
	//format it (with all the data)
	return fmt.Sprintf("{ Thread }")
}

func (ot OmmThread) Type() string {
	return "thread"
}

func (ot OmmThread) TypeOf() string {
	return ot.Type()
}

func (ot OmmThread) Deallocate() {
	//threads must be deallocated with a special way (because of the menace that is c)
	//but, i am grateful for c because without it, i would have has to use asm
	//and with asm, i would have to write 20 asm source files for each platform

	thread_dealloc(ot)
}