package types

type TuskUndef struct {
	None struct{}
}

func (u TuskUndef) Format() string {
	return "undef"
}

func (u TuskUndef) Type() string {
	return "none"
}

func (u TuskUndef) TypeOf() string {
	return u.Type()
}

func (u TuskUndef) Deallocate() {}

func (u TuskUndef) Clone() *TuskType {
	return nil
}

//Range ranges over an undef type
func (u TuskUndef) Range(fn func(val1, val2 *TuskType) (Returner, *TuskError)) (*Returner, *TuskError) {
	return nil, nil
}
