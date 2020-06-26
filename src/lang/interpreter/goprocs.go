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

}
