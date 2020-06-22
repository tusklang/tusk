package interpreter

import "strconv"

//string + (* - array - hash) = string
func string__plus__all_not_array_not_hash(num1, num2 Action, cli_params CliParams) Action {

  num1C := cast(num1, "string")
  num2C := cast(num2, "string")

  num1C.ExpStr+=num2C.ExpStr

  return num1C
}

//array + * = array
func array__plus__array(num1, num2 Action, cli_params CliParams) Action {

  if num1.Type == "array" {
    num1.Hash_Values[strconv.Itoa(len(num1.Hash_Values))] = []Action{ num2 }

    return num1
  } else {

    var newHash = make(map[string][]Action)

    for k := range num2.Hash_Values {
      _ = k
    }

    _ = newHash

    return undef //for now return undef, later change to the real code
  }
}
