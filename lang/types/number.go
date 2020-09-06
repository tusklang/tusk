package types

import (
	"fmt"
	"math"
	"strconv"
)

//number sizes

//export DigitSize
const DigitSize = 1

//export MAX_DIGIT
var MAX_DIGIT = int64(math.Pow(10, DigitSize) - 1)

//export MIN_DIGIT
var MIN_DIGIT = -1 * MAX_DIGIT

//////////////

type OmmNumber struct {
	Integer *[]int64
	Decimal *[]int64
}

func (n *OmmNumber) FromGoType(val float64) {
	numStr := fmt.Sprintf("%f", val)
	integer, decimal := BigNumConverter(numStr)
	n.Integer, n.Decimal = &integer, &decimal
}

func (n OmmNumber) ToGoType() float64 {
	f, _ := strconv.ParseFloat(NumNormalize(n), 64)
	return float64(f)
}

func (n *OmmNumber) FromString(val string) {
	integer, decimal := BigNumConverter(val)
	n.Integer, n.Decimal = &integer, &decimal
}

func (n *OmmNumber) SetInt(v []int64) {
	n.Integer = &v
}

func (n *OmmNumber) SetDec(v []int64) {
	n.Decimal = &v
}

func (n OmmNumber) Clone() OmmNumber {
	var newNum OmmNumber
	newNum.SetInt(append([]int64{}, *n.Integer...))
	newNum.SetDec(append([]int64{}, *n.Decimal...))
	return newNum
}

func (n OmmNumber) Format() string {
	str := NumNormalize(n)
	return str
}

func (n OmmNumber) Type() string {
	return "number"
}

func (n OmmNumber) TypeOf() string {
	return n.Type()
}

func (n OmmNumber) Deallocate() {}

//Range ranges over a number
func (n OmmNumber) Range(fn func(val1, val2 *OmmType) Returner) *Returner {
	return nil
}
