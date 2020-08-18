package types

type OmmObject struct {
  Name       string
  Instance   Instance
}

func (o OmmObject) GetInstance(name string) *OmmType {

  v, exists := o.Instance.vars["$" + name]

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

func (_ OmmObject) Deallocate() {}
