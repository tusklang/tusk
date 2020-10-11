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

type TuskNumber struct {
	Integer *[]int64
	Decimal *[]int64
}

func (n *TuskNumber) FromGoType(val float64) {
	numStr := fmt.Sprintf("%f", val)
	integer, decimal := BigNumConverter(numStr)
	n.Integer, n.Decimal = &integer, &decimal
}

func (n TuskNumber) ToGoType() float64 {
	f, _ := strconv.ParseFloat(NumNormalize(n), 64)
	return float64(f)
}

func (n *TuskNumber) FromString(val string) {
	integer, decimal := BigNumConverter(val)
	n.Integer, n.Decimal = &integer, &decimal
}

func (n *TuskNumber) SetInt(v []int64) {
	n.Integer = &v
}

func (n *TuskNumber) SetDec(v []int64) {
	n.Decimal = &v
}

func (n TuskNumber) Clone() TuskNumber {
	var newNum TuskNumber
	newNum.SetInt(append([]int64{}, *n.Integer...))
	newNum.SetDec(append([]int64{}, *n.Decimal...))
	return newNum
}

func (n TuskNumber) Format() string {
	str := NumNormalize(n)
	return str
}

func (n TuskNumber) Type() string {
	return "number"
}

func (n TuskNumber) TypeOf() string {
	return n.Type()
}

func (n TuskNumber) Deallocate() {}

//Range ranges over a number
func (n TuskNumber) Range(fn func(val1, val2 *TuskType) (Returner, *TuskError)) (*Returner, *TuskError) {
	return nil, nil
}
