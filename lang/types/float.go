package types

import "fmt"

//float is just a way to make a number without sacrificing so much speed
type OmmFloat struct {
	Gofloat float64
}

func (b OmmFloat) Format() string {
	return fmt.Sprint(b.Gofloat)
}

func (b OmmFloat) Type() string {
	return "float"
}

func (b OmmFloat) TypeOf() string {
	return b.Type()
}

func (_ OmmFloat) Deallocate() {}
