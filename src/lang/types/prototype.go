package types

type OmmProto struct {

  ProtoName        string
  Static           map[string]*OmmType
  Instance         map[string]*OmmType

}

func (p *OmmProto) Set(static, instance map[string]*OmmType) {
  p.Static, p.Instance = static, instance
}

func (p OmmProto) GetStatic(name string) *OmmType {

  v, exists := p.Static["$" + name]

  if !exists {
    return nil
  }

  return v
}

func (p OmmProto) Format() string {
  return "{" + p.ProtoName[1:] + "}"
}

func (p OmmProto) Type() string {
  return "proto"
}

func (p OmmProto) TypeOf() string {
  return p.ProtoName[1:] /* remove the leading $ */ + " prototype"
}
