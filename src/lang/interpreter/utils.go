package interpreter

import . "lang/types"

var tmpFalse = false
var tmpTrue = true

//list of commonly used values
var undef = OmmUndef{}
var zero = OmmNumber{
  Integer: &[]int64{0},
}
var one = OmmNumber{
  Integer: &[]int64{1},
}
var neg_one = OmmNumber{
  Integer: &[]int64{-1},
}
var falseAct = OmmBool{
  Boolean: &tmpFalse,
}
var trueAct = OmmBool{
  Boolean: &tmpTrue,
}
//////////////////////////////

//ensure that the decimal doesnt grow too much
func ensurePrec(num1, num2 *OmmNumber, cli_params CliParams) {

  if len(*(*num1).Decimal) > cli_params["Calc"]["PREC"].(int) {
    *(*num1).Decimal = (*(*num1).Decimal)[len(*(*num1).Decimal) - cli_params["Calc"]["PREC"].(int):]
  }
  if len(*(*num2).Decimal) > cli_params["Calc"]["PREC"].(int) {
    (*(*num2).Decimal) = (*(*num2).Decimal)[len(*(*num2).Decimal) - cli_params["Calc"]["PREC"].(int):]
  }
}
