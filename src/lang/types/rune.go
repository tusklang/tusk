package types

type OmmRune struct {
  Rune *rune
}

func (r *OmmRune) FromGoType(val rune) {
  r.Rune = &val
}

func (r OmmRune) ToGoType() rune {
  return *r.Rune
}

func (r OmmRune) Format() string {
  return string(*r.Rune)
}

func (r OmmRune) Type() string {
  return "rune"
}

func (r OmmRune) TypeOf() string {
  return r.Type()
}
