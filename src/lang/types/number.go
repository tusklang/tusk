package types

import "strconv"
import "math"

//number sizes

//export DigitSize
const DigitSize = 1;
//export MAX_DIGIT
var MAX_DIGIT = int64(math.Pow(10, DigitSize) - 1)
//export MIN_DIGIT
var MIN_DIGIT = -1 * MAX_DIGIT

//////////////

type OmmNumber struct {
  Integer *[]int64
  Decimal *[]int64
}

func (n *OmmNumber) FromGoType(val int64) {
  numStr := strconv.FormatInt(val, 10)
  integer, decimal := BigNumConverter(numStr)
  n.Integer, n.Decimal = &integer, &decimal
}

func (n OmmNumber) ToGoType() int64 {
  return 0 //for now
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

func (_ OmmNumber) ValueFunc() {}
