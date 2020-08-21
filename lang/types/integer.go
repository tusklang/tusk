package types

import "strconv"

//integer is just a way to make a number without sacrificing so much speed
type OmmInteger struct {
	Goint int64
}

func (i OmmInteger) Format() string {
	return strconv.FormatInt(i.Goint, 10)
}

func (i OmmInteger) Type() string {
	return "int"
}

func (i OmmInteger) TypeOf() string {
	return i.Type()
}

func (_ OmmInteger) Deallocate() {}
