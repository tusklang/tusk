package interpreter

import (
	"math/big"

	. "github.com/omm-lang/omm/lang/types"
)

//convert the numbers (integers/decimals) to bigints
func sliceToBigInt(slice []int64) *big.Int {
	bigVal := big.NewInt(0)
	places := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(len(slice)-1)), nil)

	for i := len(slice) - 1; i >= 0; i-- {
		bigVal = bigVal.Add(bigVal, new(big.Int).Mul(places, big.NewInt(slice[i])))
		places.Div(places, big.NewInt(10)) //divide the place by 10
	}

	return bigVal
}

func toBig(val1, val2 OmmNumber) (*big.Int, *big.Int, *big.Int, *big.Int) {
	int1, dec1 := val1.Integer, val1.Decimal
	int2, dec2 := val2.Integer, val2.Decimal

	var (
		bigInt1 = big.NewInt(0)
		bigInt2 = big.NewInt(0)
		bigDec1 = big.NewInt(0)
		bigDec2 = big.NewInt(0)
	)

	if int1 != nil {
		bigInt1 = sliceToBigInt(*int1)
	}
	if int2 != nil {
		bigInt2 = sliceToBigInt(*int2)
	}
	if dec1 != nil {
		bigDec1 = sliceToBigInt(*dec1)
	}
	if dec2 != nil {
		bigDec2 = sliceToBigInt(*dec2)
	}

	return bigInt1, bigInt2, bigDec1, bigDec2
}

func isTruthy(val OmmType) bool {

	switch val.(type) {
	case OmmBool:

		if val.(OmmBool).Boolean == nil {
			return false
		}

		return *val.(OmmBool).Boolean
	case OmmUndef:
		return false
	}

	return true
}

func isLess(val1, val2 OmmNumber) bool {

	//manage for leading zeros in the decimal
	var num1ZeroLeadingDec int
	var num2ZeroLeadingDec int

	num1DecGreater := false //if num1 has less leading zeros

	if val1.Decimal != nil {
		for len(*val1.Decimal) != 0 && (*val1.Decimal)[len(*val1.Decimal)-1] == 0 {
			*val1.Decimal = (*val1.Decimal)[:len(*val1.Decimal)-1]
			num1ZeroLeadingDec++
		}
		if len(*val1.Decimal) == 0 { //this means the number looks like 1.00000
			num1ZeroLeadingDec = 0
		}
	}
	if val2.Decimal != nil {
		for len(*val2.Decimal) != 0 && (*val2.Decimal)[len(*val2.Decimal)-1] == 0 {
			*val2.Decimal = (*val2.Decimal)[:len(*val2.Decimal)-1]
			num2ZeroLeadingDec++
		}
		if len(*val2.Decimal) == 0 { //this means the number looks like 1.000000000000000000000000000
			num2ZeroLeadingDec = 0
		}
	}

	if num1ZeroLeadingDec > num2ZeroLeadingDec {
		num1DecGreater = true
	}
	/////////////////////////////////////////

	bigint1, bigint2, bigdec1, bigdec2 := toBig(val1, val2)

	if bigint1.Cmp(bigint2) != 0 {
		return bigint1.Cmp(bigint2) == -1
	}

	if num1DecGreater { //if num1 has more leading zeros, return true(because the decimal is less)
		return true
	}

	return bigdec1.Cmp(bigdec2) == -1
}

func isEqual(val1, val2 OmmNumber) bool {

	//manage for leading zeros in the decimal
	var num1ZeroLeadingDec int
	var num2ZeroLeadingDec int

	num1DecNotEqual := true //if num1 has less leading zeros

	if val1.Decimal != nil {
		for len(*val1.Decimal) != 0 && (*val1.Decimal)[len(*val1.Decimal)-1] == 0 {
			*val1.Decimal = (*val1.Decimal)[:len(*val1.Decimal)-1]
			num1ZeroLeadingDec++
		}
		if len(*val1.Decimal) == 0 { //this means the number looks like 1.000000000000000000000000000
			num1ZeroLeadingDec = 0
		}
	}
	if val2.Decimal != nil {
		for len(*val2.Decimal) != 0 && (*val2.Decimal)[len(*val2.Decimal)-1] == 0 {
			*val2.Decimal = (*val2.Decimal)[:len(*val2.Decimal)-1]
			num2ZeroLeadingDec++
		}
		if len(*val2.Decimal) == 0 { //this means the number looks like 1.000000000000000000000000000
			num2ZeroLeadingDec = 0
		}
	}

	if num1ZeroLeadingDec == num2ZeroLeadingDec {
		num1DecNotEqual = false
	}
	/////////////////////////////////////////

	bigint1, bigint2, bigdec1, bigdec2 := toBig(val1, val2)

	if bigint1.Cmp(bigint2) != 0 {
		return false
	}

	if num1DecNotEqual {
		return false
	}

	return bigdec1.Cmp(bigdec2) == 0
}

func isLessOrEqual(val1, val2 OmmNumber) bool {

	//manage for leading zeros in the decimal
	var num1ZeroLeadingDec int
	var num2ZeroLeadingDec int

	num1DecLess := false //if num1 has less leading zeros

	if val1.Decimal != nil {
		for len(*val1.Decimal) != 0 && (*val1.Decimal)[len(*val1.Decimal)-1] == 0 {
			*val1.Decimal = (*val1.Decimal)[:len(*val1.Decimal)-1]
			num1ZeroLeadingDec++
		}
		if len(*val1.Decimal) == 0 { //this means the number looks like 1.000000000000000000000000000
			num1ZeroLeadingDec = 0
		}
	}
	if val2.Decimal != nil {
		for len(*val2.Decimal) != 0 && (*val2.Decimal)[len(*val2.Decimal)-1] == 0 {
			*val2.Decimal = (*val2.Decimal)[:len(*val2.Decimal)-1]
			num2ZeroLeadingDec++
		}
		if len(*val2.Decimal) == 0 { //this means the number looks like 1.000000000000000000000000000
			num2ZeroLeadingDec = 0
		}
	}

	if num1ZeroLeadingDec < num2ZeroLeadingDec {
		num1DecLess = true
	}
	/////////////////////////////////////////

	bigint1, bigint2, bigdec1, bigdec2 := toBig(val1, val2)

	if bigdec1.Cmp(bigdec2) <= 0 && !num1DecLess {
		return bigint1.Cmp(bigint2) <= 0
	}

	if num1DecLess {
		return false
	}

	return bigdec1.Cmp(bigdec2) <= 0
}

func abs(val OmmNumber, stacktrace []string, cli_params CliParams) OmmType {

	if isLess(val, zero) {
		var tmpInstance Instance
		tmpInstance.Params = cli_params
		return *number__times__number(val, neg_one, &tmpInstance, stacktrace, 0, "none" /* using this because it will not throw an error */)
	}

	return val
}
