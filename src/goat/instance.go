package goat

import . "lang/interpreter"
import . "lang/types"

//export NewInstance
func NewInstance(oat Oat) *Instance {
  var ins Instance
  ins.FromOatStruct(oat)
  return &ins
}
