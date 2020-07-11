package types

type OmmUndef struct {
  value struct{}
}

func (b *OmmUndef) FromGoType(val struct{}) {
  b.value = struct{}{}
}

func (b OmmUndef) ToGoType() struct{} {
  var nilv struct{}
  return nilv
}

func (b OmmUndef) Format() string {
  return "undef"
}

func (arr OmmUndef) Type() string {
  return "undef"
}
