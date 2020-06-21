package interpreter

func indexesCalc(val Action, indexes [][]Action, cli_params CliParams, vars map[string]Variable, this_vals []Action, dir string) Action {

  if len(indexes) == 0 {
    return val
  }

  key := interpreter(indexes[0], cli_params, vars, true, this_vals, dir).Exp.ExpStr

  if _, exists := val.Hash_Values[key]; exists {
    if _, falseyValExists := val.Hash_Values["falsey"]; falseyValExists {
      return val.Hash_Values["falsey"][0]
    }
    return undef
  }

  //if it is a public key
  if val.Hash_Values[key][0].Access == "public" {

    expVal := interpreter(val.Hash_Values[key], cli_params, vars, true, this_vals, dir).Exp
    indexes = indexes[1:]
    return indexesCalc(expVal, indexes, cli_params, vars, this_vals, dir)

  }

  return undef
}
