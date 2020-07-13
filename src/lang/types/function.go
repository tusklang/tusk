package types

import "strconv"

type OmmFunc struct {
  Params []string
  Body   []Action
}

func (f OmmFunc) Format() string {
  return "(PARAM COUNT: " + strconv.Itoa(len(f.Params)) + ") { ... }"
}

func (arr OmmFunc) Type() string {
  return "function"
}