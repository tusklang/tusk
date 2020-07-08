package interpreter

import . "lang/types"

//list of commonly used values
var undef = Action{}
var arr = Action{}
var hash = Action{}
var zero = Action{}
var one = Action{}
var neg_one = Action{}
var falseAct = Action{}
var trueAct = Action{}
var emptyString = Action{}
var emptyRune = Action{}
var thread = Action{}
//////////////////////////////

//ensure that the decimal doesnt grow too much
func ensurePrec(num1, num2 *Action, cli_params CliParams) {

  if len((*num1).Decimal) > cli_params["Calc"]["PREC"].(int) {
    (*num1).Decimal = (*num1).Decimal[len((*num1).Decimal) - cli_params["Calc"]["PREC"].(int):]
  }
  if len((*num2).Decimal) > cli_params["Calc"]["PREC"].(int) {
    (*num2).Decimal = (*num2).Decimal[len((*num2).Decimal) - cli_params["Calc"]["PREC"].(int):]
  }
}
