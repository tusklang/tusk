package types

type OmmUndef struct {
  value struct{}
}

func (u *OmmUndef) FromGoType(val struct{}) {
  u.value = struct{}{}
}

func (u OmmUndef) ToGoType() struct{} {
  var nilv struct{}
  return nilv
}

func (u OmmUndef) Format() string {
  return "undef"
}

func (u OmmUndef) Type() string {
  return "undef"
}

func (u OmmUndef) TypeOf() string {
  return u.Type()
}
