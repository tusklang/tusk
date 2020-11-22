package types

import "fmt"

type TuskInt struct {
	Int int64
}

func (i *TuskInt) FromGoType(val int64) {
	i.Int = val
}

func (i TuskInt) ToGoType() int64 {
	return i.Int
}

func (i TuskInt) Format() string {
	return fmt.Sprintf("%d", i.Int)
}

func (i TuskInt) Type() string {
	return "int"
}

func (i TuskInt) TypeOf() string {
	return i.Type()
}

func (i TuskInt) Deallocate() {}

func (i TuskInt) Clone() *TuskType {
	var tmp TuskType = i
	return &tmp //just return a new adress to the same value
}

//Range ranges over an int
func (i TuskInt) Range(fn func(val1, val2 *TuskType) (Returner, *TuskError)) (*Returner, *TuskError) {
	return nil, nil
}
