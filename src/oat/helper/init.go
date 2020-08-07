package oatHelper

import "encoding/gob"
import . "lang/types"

func init() {
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
