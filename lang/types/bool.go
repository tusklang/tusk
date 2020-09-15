package types

import "strconv"

type KaBool struct {
	Boolean *bool
}

func (b *KaBool) FromGoType(val bool) {
	b.Boolean = &val
}

func (b KaBool) ToGoType() bool {

	if b.Boolean == nil {
		return false
	}

	return *b.Boolean
}

func (b KaBool) Format() string {
	return strconv.FormatBool(*b.Boolean)
}

func (b KaBool) Type() string {
	return "bool"
}

func (b KaBool) TypeOf() string {
	return b.Type()
}

func (b KaBool) Deallocate() {}

//Range ranges over a bool
func (b KaBool) Range(fn func(val1, val2 *KaType) Returner) *Returner {
	return nil
}
