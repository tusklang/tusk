package oat

import "oat/helper"

import . "lang/types"
import . "lang/interpreter"

//export Run
func Run(params CliParams) {

  decoded := oatHelper.FromOat(params.Name)

  //run the oat
  RunInterpreter(decoded.Variables, params)
}
