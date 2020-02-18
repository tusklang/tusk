package main

import "strings"

func modulo(_num1 string, _num2 string, calc_params paramCalcOpts, line uint64, functions []Funcs) string {

  calc_params.precision = 0;

  divved_ := addDec(division(_num1, _num2, calc_params, line, functions))

  divved := divved_[:strings.Index(divved_, ".")]

  mult := multiply(divved, _num2, calc_params, line, functions)
  remainder := subtract(_num1, mult, calc_params, line, functions)

  return remainder
}
