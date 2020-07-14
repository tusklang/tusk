package types

type OmmObject struct {
  prototype  OmmProto
  instance   map[string]*OmmType
}

func (o OmmObject) New(proto OmmProto) OmmObject { //factory to create a new object
  o.prototype = proto
  o.instance = make(map[string]*OmmType)

  //copy the default instance to the object
  for k, v := range proto.Instance {
    o.instance[k] = v
  }
  /////////////////////////////////////////

  return o
}

func (o OmmObject) GetInstance(name string) *OmmType {

  v, exists := o.instance["$" + name]

  if !exists {
    return nil
  }

  return v
}

func (o OmmObject) Format() string {
  return "{object}"
}

func (o OmmObject) Type() string {
  return "object"
}
