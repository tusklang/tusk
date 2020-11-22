package types

import "fmt"

type TuskFloat struct {
	Float float64
}

func (f *TuskFloat) FromGoType(val float64) {
	f.Float = val
}

func (f TuskFloat) ToGoType() float64 {
	return f.Float
}

func (f TuskFloat) Format() string {
	return fmt.Sprintf("%f", f.Float)
}

func (f TuskFloat) Type() string {
	return "float"
}

func (f TuskFloat) TypeOf() string {
	return f.Type()
}

func (f TuskFloat) Deallocate() {}

func (f TuskFloat) Clone() *TuskType {
	var tmp TuskType = f
	return &tmp //just return a new adress to the same value
}

//Range ranges over a float
func (f TuskFloat) Range(fn func(val1, val2 *TuskType) (Returner, *TuskError)) (*Returner, *TuskError) {
	return nil, nil
}
