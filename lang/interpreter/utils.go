package interpreter

import . "github.com/tusklang/tusk/lang/types"

var tmpFalse = false
var tmpTrue = true

//list of ckaonly used values
var undef = TuskUndef{}
var zero = TuskNumber{
	Integer: &[]int64{0},
	Decimal: &[]int64{0},
}
var one = TuskNumber{
	Integer: &[]int64{1},
	Decimal: &[]int64{0},
}
var neg_one = TuskNumber{
	Integer: &[]int64{-1},
	Decimal: &[]int64{0},
}
var falsev = TuskBool{
	Boolean: &tmpFalse,
}
var truev = TuskBool{
	Boolean: &tmpTrue,
}

//////////////////////////////

//ensure that the decimal doesnt grow too much
func ensurePrec(num1, num2 *TuskNumber, cli_params CliParams) {

	//ensure a nil pointer error doesnt happen
	if (*num1).Decimal == nil {
		tmp := []int64{0}
		(*num1).Decimal = &tmp
	}
	if (*num2).Decimal == nil {
		tmp := []int64{0}
		(*num2).Decimal = &tmp
	}

	//using cli_params precision + 1 because everything must be a float (and decimal must be >= 1)
	if uint64(len(*(*num1).Decimal)) > cli_params.Prec+1 {
		*(*num1).Decimal = (*(*num1).Decimal)[uint64(len(*(*num1).Decimal))-cli_params.Prec:]
	}
	if uint64(len(*(*num2).Decimal)) > cli_params.Prec+1 {
		(*(*num2).Decimal) = (*(*num2).Decimal)[uint64(len(*(*num2).Decimal))-cli_params.Prec:]
	}

}
