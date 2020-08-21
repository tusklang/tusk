package types

import "fmt"

//float is just a way to make a number without sacrificing so much speed
type OmmFloat struct {
	Gofloat float64
}

func (f OmmFloat) Format() string {
	return fmt.Sprint(f.Gofloat)
}

func (f OmmFloat) Type() string {
	return "float"
}

func (f OmmFloat) TypeOf() string {
	return f.Type()
}

func (_ OmmFloat) Deallocate() {}
