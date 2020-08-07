package interpreter

//type-function and type-object files are here because they use the `Interpreter` structure

import . "lang/types"

type OmmObject struct {
  Name       string
  Instance   Instance
}

func (o OmmObject) GetInstance(name string) *OmmType {

  v, exists := o.Instance.globals["$" + name]

  if !exists {
    return nil
  }

  return v.Value
}

func (o OmmObject) Format() string {
  return "{" + o.Name[1:] + "}"
}

func (o OmmObject) Type() string {
  return "object"
}

func (o OmmObject) TypeOf() string {
  return o.Name[1:]
}
