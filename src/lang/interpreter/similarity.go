package interpreter

import . "lang/types"

func similar(val1, val2, degree Action, cli_params CliParams, vars map[string]Variable, this_vals []Action, dir string) bool {

  if degree.Type == "falsey" { //if the degree is undef, the degree is zero
    degree = zero
  }

  //if the degree is not a number return false
  if degree.Type != "number" {
    return false
  }

  if val1.Name == "hashed_value" && val2.Name == "hashed_value" { //if both values are hashes

    //ensure val1 has more hash values than val2
    if len(val1.Hash_Values) < len(val2.Hash_Values) {
      val1, val2 = val2, val1
    }

    var difAmt = zero

    for k, v := range val1.Hash_Values {

      //if val2 does not have the value, add one to the difAmt
      if _, exists := val2.Hash_Values[k]; !exists {
        goto addone
      }

      if isEqual(interpreter([]Action{ v }, cli_params, vars, true, this_vals, dir).Exp, interpreter([]Action{ val2.Hash_Values[k] }, cli_params, vars, true, this_vals, dir).Exp) {
        continue
      }

      addone:
      difAmt = number__plus__number(difAmt, one, cli_params)

      //if the difAmt > degree, they are not similar
      if !isLessOrEqual(difAmt, degree) {
        return false
      }
    }

  } else {

    if isEqual(degree, zero) {
      //if it is 0, no need to add (also it serves as lazy equality)

      //cast to a string, then compare
      return cast(val1, "string").ExpStr == cast(val2, "string").ExpStr
    } else {

      upper := add(val1, degree, cli_params)
      lower := subtract(val1, degree, cli_params)

      if !isLess(upper, val2) /* if upper >= val2 */ && isLessOrEqual(lower, val2) /* lower <= val2 */ {
        return true
      } else {
        return false
      }

    }

  }

  return true
}

func strictSimilar(val1, val2, degree Action, cli_params CliParams, vars map[string]Variable, this_vals []Action, dir string) bool {

  if degree.Type == "falsey" { //if the degree is undef, the degree is zero
    degree = zero
  }

  //if the degree is not a number return false
  if degree.Type != "number" {
    return false
  }

  if val1.Name == "hashed_value" && val2.Name == "hashed_value" { //if both values are hashes

    //ensure val1 has more hash values than val2
    if len(val1.Hash_Values) < len(val2.Hash_Values) {
      val1, val2 = val2, val1
    }

    var difAmt = zero

    for k, v := range val1.Hash_Values {

      //if val2 does not have the value, add one to the difAmt
      if _, exists := val2.Hash_Values[k]; !exists {
        goto addone
      }

      if isEqual(interpreter([]Action{ v }, cli_params, vars, true, this_vals, dir).Exp, interpreter([]Action{ val2.Hash_Values[k] }, cli_params, vars, true, this_vals, dir).Exp) {
        continue
      }

      addone:
      difAmt = number__plus__number(difAmt, one, cli_params)
    }

    if !isEqual(difAmt, degree) { //if degree != difAmt, they are not strictly similar
      return false
    }

  } else {

    if isEqual(degree, zero) {
      //if it is 0, no need to add
      return equals(val1, val2)
    } else {

      upper := add(val1, degree, cli_params)
      lower := subtract(val1, degree, cli_params)

      if !equals(upper, val2) /* if upper != val2 */ && !equals(lower, val2) /* lower != val2 */ {
        return false
      }

    }

  }

  return true
}
