package oatInitGob

import "encoding/gob"

import . "lang/types"
import . "lang/interpreter"

//export InitGob
func InitGob() {
	//register the value types
	gob.Register(OmmArray{})
	gob.Register(OmmBool{})
	gob.Register(OmmFunc{})
	gob.Register(OmmHash{})
	gob.Register(OmmNumber{})
	gob.Register(OmmObject{})
	gob.Register(OmmProto{})
	gob.Register(OmmRune{})
	gob.Register(OmmString{})
	gob.Register(OmmThread{})
	gob.Register(OmmUndef{})
	//////////////////////////
}