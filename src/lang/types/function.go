package types

type OmmFunc struct {
  Params []Action
  Body   []Action
}

func (_ OmmFunc) ValueFunc() {}
