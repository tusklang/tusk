package types

import "strconv"

//integer is just a way to make a number without sacrificing so much speed
type OmmInteger struct {
	Goint int64
}

func (b OmmInteger) Format() string {
	return strconv.FormatInt(b.Goint, 10)
}
  
func (b OmmInteger) Type() string {
	return "int"
}
  
func (b OmmInteger) TypeOf() string {
	return b.Type()
}

func (_ OmmInteger) Deallocate() {}
