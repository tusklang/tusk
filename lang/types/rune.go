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

func (r OmmRune) Deallocate() {}

//Range ranges over a rune
func (r OmmRune) Range(fn func(val1, val2 *OmmType) Returner) *Returner {
	return nil
}
