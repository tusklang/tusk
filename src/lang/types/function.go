package types

import "strconv"

type OmmFunc struct {
  Params []string
  Body   []Action
}

func (f OmmFunc) Format() string {
  return "(PARAM COUNT: " + strconv.Itoa(len(f.Params)) + ") { ... }"
}

func (f OmmFunc) Type() string {
  return "function"
}

func (f OmmFunc) TypeOf() string {
  return f.Type()
}
