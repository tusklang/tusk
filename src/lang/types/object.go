package types

type OmmObject struct {
  Name       string
  Instance   map[string]*OmmType
}

func (o OmmObject) New(proto OmmProto) OmmObject { //factory to create a new object
  o.Name = proto.ProtoName
  o.Instance = make(map[string]*OmmType)

  //copy the default instance to the object
  for k, v := range proto.Instance {
    o.Instance[k] = v
  }
  /////////////////////////////////////////

  return o
}

func (o OmmObject) GetInstance(name string) *OmmType {

  v, exists := o.Instance["$" + name]

  if !exists {
    return nil
  }

  return v
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
