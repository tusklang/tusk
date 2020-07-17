package types

type OmmUndef struct {
  None struct{}
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
