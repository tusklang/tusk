package types

type OmmProto struct {

  SelfName         string
  static           map[string]*OmmType
  Instance         map[string]*OmmType

}

func (p *OmmProto) Set(static, instance map[string]*OmmType) {
  p.static, p.Instance = static, instance
}

func (p OmmProto) GetStatic(name string) *OmmType {

  v, exists := p.static["$" + name]

  if !exists {
    return nil
  }

  return v
}

func (p OmmProto) Format() string {
  return "{prototype}"
}

func (p OmmProto) Type() string {
  return "proto"
}
