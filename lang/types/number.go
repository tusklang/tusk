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

type KaNumber struct {
	Integer *[]int64
	Decimal *[]int64
}

func (n *KaNumber) FromGoType(val float64) {
	numStr := fmt.Sprintf("%f", val)
	integer, decimal := BigNumConverter(numStr)
	n.Integer, n.Decimal = &integer, &decimal
}

func (n KaNumber) ToGoType() float64 {
	f, _ := strconv.ParseFloat(NumNormalize(n), 64)
	return float64(f)
}

func (n *KaNumber) FromString(val string) {
	integer, decimal := BigNumConverter(val)
	n.Integer, n.Decimal = &integer, &decimal
}

func (n *KaNumber) SetInt(v []int64) {
	n.Integer = &v
}

func (n *KaNumber) SetDec(v []int64) {
	n.Decimal = &v
}

func (n KaNumber) Clone() KaNumber {
	var newNum KaNumber
	newNum.SetInt(append([]int64{}, *n.Integer...))
	newNum.SetDec(append([]int64{}, *n.Decimal...))
	return newNum
}

func (n KaNumber) Format() string {
	str := NumNormalize(n)
	return str
}

func (n KaNumber) Type() string {
	return "number"
}

func (n KaNumber) TypeOf() string {
	return n.Type()
}

func (n KaNumber) Deallocate() {}

//Range ranges over a number
func (n KaNumber) Range(fn func(val1, val2 *KaType) Returner) *Returner {
	return nil
}
