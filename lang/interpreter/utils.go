package interpreter

import . "omm/lang/types"

var tmpFalse = false
var tmpTrue = true

//list of commonly used values
var undef = OmmUndef{}
var zero = OmmNumber{
	Integer: &[]int64{0},
	Decimal: &[]int64{0},
}
var one = OmmNumber{
	Integer: &[]int64{1},
	Decimal: &[]int64{0},
}
var neg_one = OmmNumber{
	Integer: &[]int64{-1},
	Decimal: &[]int64{0},
}
var falsev = OmmBool{
	Boolean: &tmpFalse,
}
var truev = OmmBool{
	Boolean: &tmpTrue,
}

//////////////////////////////

//ensure that the decimal doesnt grow too much
func ensurePrec(num1, num2 *OmmNumber, cli_params CliParams) {

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
