package types

type KaUndef struct {
	None struct{}
}

func (u KaUndef) Format() string {
	return "undef"
}

func (u KaUndef) Type() string {
	return "none"
}

func (u KaUndef) TypeOf() string {
	return u.Type()
}

func (u KaUndef) Deallocate() {}

//Range ranges over an undef type
func (u KaUndef) Range(fn func(val1, val2 *KaType) Returner) *Returner {
	return nil
}
