package types

import "fmt"

type TuskRune struct {
	Rune *rune
}

func (r *TuskRune) FromGoType(val rune) {
	r.Rune = &val
}

func (r TuskRune) ToGoType() rune {
	return *r.Rune
}

func (r TuskRune) Format() string {
	return fmt.Sprintf("%c", *r.Rune)
}

func (r TuskRune) Type() string {
	return "rune"
}

func (r TuskRune) TypeOf() string {
	return r.Type()
}

func (r TuskRune) Deallocate() {}

//Range ranges over a rune
func (r TuskRune) Range(fn func(val1, val2 *TuskType) (Returner, *TuskError)) (*Returner, *TuskError) {
	return nil, nil
}
