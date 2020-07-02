package interpreter

//all of the goprocs
//functions written in go that are used by omm

import "strconv"
import "time"

var goprocs = map[string]func(args [][]Action, cli_params CliParams, vars map[string]Variable, this_vals []Action, dir string) Action {

  "wait": func(args [][]Action, cli_params CliParams, vars map[string]Variable, this_vals []Action, dir string) Action {

    if len(args) != 1 {
      return undef
    }

    interpreted := interpreter(args[0], cli_params, vars, true, this_vals, dir).Exp
    amt, _ := strconv.ParseFloat(num_normalize(interpreted), 64)
    time.Sleep(time.Duration(int64(amt)) * time.Nanosecond)

    return interpreted
  },
  "self": func(args [][]Action, cli_params CliParams, vars map[string]Variable, this_vals []Action, dir string) Action {

    if len(args) == 1 { //get a value
      index, _ := strconv.Atoi(num_normalize(interpreter(args[0], cli_params, vars, true, this_vals, dir).Exp))

      //if the this_vals has the index
      if len(this_vals) > index && index >= 0 {
        return this_vals[index]
      }
    } else if len(args) == 3 { //set a value

      //wont work right now

      index, _ := strconv.Atoi(num_normalize(interpreter(args[0], cli_params, vars, true, this_vals, dir).Exp))
      val := interpreter(args[2], cli_params, vars, true, this_vals, dir).Exp

      //if the this_vals has the index
      if len(this_vals) <= index || index < 0 {
        this_vals[index].Hash_Values[cast(interpreter(args[1], cli_params, vars, true, this_vals, dir).Exp, "string").ExpStr /*because it must be a string*/] = []Action{ val }
        return val
      }
    }

    return undef
  },

}
