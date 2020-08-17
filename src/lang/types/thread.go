package types

import "unsafe"

type OmmThread struct {
	thread *cthread
}

//export NewThread
func NewThread(cb unsafe.Pointer) OmmThread {
	//wrapper for other packages to create OmmThreads
	return *(*OmmThread)(newthread(cb))
}

func (ot OmmThread) Join(thread OmmThread) *OmmType {
	return jointhread(thread)
}

func (ot OmmThread) Format() string {
	return "{ Thread <~ }"
}

func (ot OmmThread) Type() string {
	return "thread"
}

func (ot OmmThread) TypeOf() string {
	return ot.Type()
}

func (ot OmmThread) Deallocate() {

}