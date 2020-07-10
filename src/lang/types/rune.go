package types

type OmmRune struct {
  rune *rune
}

func (r *OmmRune) FromGoType(val rune) {
  r.rune = &val
}

func (r OmmRune) ToGoType() rune {
  return *r.rune
}

func (r OmmRune) Format() string {
  return string(*r.rune)
}
