package interpreter

func indexesCalc(val Action, indexes [][]Action, cli_params CliParams, vars map[string]Variable, this_vals []Action, dir string) Action {

  if len(indexes) == 0 {
    return val
  }

  key := cast(interpreter(indexes[0], cli_params, vars, true, this_vals, dir).Exp, "string").ExpStr

  if _, exists := val.Hash_Values[key]; !exists {
    return undef
  }

  //if it is a hash and the key is a private value, return undef
  if val.Type == "hash" && val.Hash_Values[key][0].Access == "private" {
    return undef
  }

  expVal := interpreter(val.Hash_Values[key], cli_params, vars, true, this_vals, dir).Exp
  indexes = indexes[1:]
  return indexesCalc(expVal, indexes, cli_params, vars, this_vals, dir)
}
