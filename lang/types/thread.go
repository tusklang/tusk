package types

import "fmt"

type KaThread struct {
	thread *cthread
	ptr    uint64
}

//export NewThread
func NewThread(cb func() *KaType) *KaThread {
	//wrapper for other packages to create KaThreads
	return newthread(cb)
}

func WaitAllThreads() {
	//wait all of the remaining threads
	for _, v := range allthreads {
		v.Join()
		v.Deallocate()
	}
}

func (ot KaThread) Join() *KaType {
	return jointhread(ot)
}

func (ot KaThread) Format() string {
	//format it (with all the data)
	return fmt.Sprintf("{ Thread }")
}

func (ot KaThread) Type() string {
	return "thread"
}

func (ot KaThread) TypeOf() string {
	return ot.Type()
}

func (ot KaThread) Deallocate() {
	//threads must be deallocated with a special way (because of the menace that is c)
	//but, i am grateful for c because without it, i would have has to use asm
	//and with asm, i would have to write 20 asm source files for each platform

	thread_dealloc(ot)
}

//Range ranges over a thread
func (ot KaThread) Range(fn func(val1, val2 *KaType) Returner) *Returner {
	return nil
}
