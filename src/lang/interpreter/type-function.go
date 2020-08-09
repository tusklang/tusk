package interpreter

//type-function and type-object files are here because they use the `Instance` structure

import . "lang/types"

type Overload struct {
  Params    []string
  Types     []string
  Body      []Action
}

type OmmFunc struct {
  Overloads []Overload
  Instance   *Instance
}

func (f OmmFunc) Format() string {
  return "(function) { ... }"
}

func (f OmmFunc) Type() string {
  return "function"
}

func (f OmmFunc) TypeOf() string {
  return f.Type()
}
