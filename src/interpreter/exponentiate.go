package main

import "strings"
import "os"
import "fmt"

func exponentiate(_num1 string, _num2 string, calc_params paramCalcOpts, line uint64, functions []Funcs) string {
  _num1 = returnInit(_num1)
  _num2 = returnInit(_num2)

  if strings.Contains(_num2, ".") {
    fmt.Println("There Was An Error: Currently You Cannot Exponentiate By Numbers With Decimals\n\n" + _num1 + "^" + _num2 + "\n" +"^^^ <- Error On Line " + string(line))
    os.Exit(1)
  }

  var final = "1"

  if strings.HasPrefix(_num2, "-") {
    _num2 = _num2[1:]

    for ;isLess("0", _num2); {
      final = multiply(final, _num1, calc_params, line, functions)
      _num2 = subtract(_num2, "1", calc_params, line, functions)

      if calc_params.logger {
        fmt.Println("Omm Logger ~ Exponentiation: " + final)
      }
    }

    final = division("1", final, calc_params, line, functions)
  } else {

    for ;isLess("0", _num2); {

      final = multiply(final, _num1, calc_params, line, functions)
      _num2 = subtract(_num2, "1", calc_params, line, functions)

      if calc_params.logger {
        fmt.Println("Omm Logger ~ Exponentiation: " + final)
      }
    }
  }

  return final
}
