package ox

/*
  ffi for omm
  (named ox because go -> omm is called goat, and I needed another animal)
*/

import . "lang/types"

type Ox interface {
  CallFunc(name string) *OmmType
}
