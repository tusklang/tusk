package types

import (
	"fmt"
	"math"
	"math/big"
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

//ToBigInt converts a Tusk number to a *big.Int
func (n TuskNumber) ToBigInt() *big.Int {
	i, _ := new(big.Int).SetString(NumNormalize(n), 10)
	return i
}

//FromBigInt converts a *big.Int to a Tusk number
func (n *TuskNumber) FromBigInt(val *big.Int) {
	integer, decimal := BigNumConverter(val.String())
	n.Integer, n.Decimal = &integer, &decimal
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

func (n TuskNumber) Format() string {
	str := NumNormalize(n)
	return str
}

func (n TuskNumber) Type() string {
	return "big"
}

func (n TuskNumber) TypeOf() string {
	return n.Type()
}

func (n TuskNumber) Deallocate() {}

func (n TuskNumber) Clone() *TuskType {
	var number = n

	//copy the integer and decimal
	var integer = append([]int64{}, *number.Integer...)

	var decimal []int64

	if number.Decimal != nil {
		decimal = append([]int64{}, *number.Decimal...)
	}
	//////////////////////////////

	var newnum TuskType = TuskNumber{
		Integer: &integer,
		Decimal: &decimal,
	}
	return &newnum
}

//Range ranges over a number
func (n TuskNumber) Range(fn func(val1, val2 *TuskType) (Returner, *TuskError)) (*Returner, *TuskError) {
	return nil, nil
}
