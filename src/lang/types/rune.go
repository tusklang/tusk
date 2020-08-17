package types

import "fmt"

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
  return fmt.Sprintf("%c", *r.Rune)
}

func (r OmmRune) Type() string {
  return "rune"
}

func (r OmmRune) TypeOf() string {
  return r.Type()
}

func (_ OmmRune) Deallocate() {}
