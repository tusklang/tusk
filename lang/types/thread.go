package types

import "fmt"

type TuskThread struct {
	thread *cthread
	ptr    uint64
}

//export NewThread
func NewThread(cb func() (*TuskType, *TuskError)) *TuskThread {
	//wrapper for other packages to create TuskThreads
	return newthread(cb)
}

func WaitAllThreads() {
	//wait all of the remaining threads
	for _, v := range allthreads {
		v.Join()
		v.Deallocate()
	}
}

func (ot TuskThread) Join() *TuskType {
	return jointhread(ot)
}

func (ot TuskThread) Format() string {
	//format it (with all the data)
	return fmt.Sprintf("{ Thread }")
}

func (ot TuskThread) Type() string {
	return "thread"
}

func (ot TuskThread) TypeOf() string {
	return ot.Type()
}

func (ot TuskThread) Deallocate() {
	//threads must be deallocated with a special way (because of the menace that is c)
	//but, i am grateful for c because without it, i would have has to use asm
	//and with asm, i would have to write 20 asm source files for each platform

	thread_dealloc(ot)
}

//Range ranges over a thread
func (ot TuskThread) Range(fn func(val1, val2 *TuskType) (Returner, *TuskError)) (*Returner, *TuskError) {
	return nil, nil
}
