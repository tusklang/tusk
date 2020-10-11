package types

import "strconv"

type TuskBool struct {
	Boolean *bool
}

func (b *TuskBool) FromGoType(val bool) {
	b.Boolean = &val
}

func (b TuskBool) ToGoType() bool {

	if b.Boolean == nil {
		return false
	}

	return *b.Boolean
}

func (b TuskBool) Format() string {
	return strconv.FormatBool(*b.Boolean)
}

func (b TuskBool) Type() string {
	return "bool"
}

func (b TuskBool) TypeOf() string {
	return b.Type()
}

func (b TuskBool) Deallocate() {}

//Range ranges over a bool
func (b TuskBool) Range(fn func(val1, val2 *TuskType) (Returner, *TuskError)) (*Returner, *TuskError) {
	return nil, nil
}
