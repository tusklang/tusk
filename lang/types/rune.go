package types

import "fmt"

type KaRune struct {
	Rune *rune
}

func (r *KaRune) FromGoType(val rune) {
	r.Rune = &val
}

func (r KaRune) ToGoType() rune {
	return *r.Rune
}

func (r KaRune) Format() string {
	return fmt.Sprintf("%c", *r.Rune)
}

func (r KaRune) Type() string {
	return "rune"
}

func (r KaRune) TypeOf() string {
	return r.Type()
}

func (r KaRune) Deallocate() {}

//Range ranges over a rune
func (r KaRune) Range(fn func(val1, val2 *KaType) Returner) *Returner {
	return nil
}
