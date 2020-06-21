package interpreter

import . "lang"

type Variable struct {
  Type      string
  Name      string
  Value     Action
  GoProc    func(actions []Action, cli_params CliParams, map[string]Variable vars, expReturn bool, this_vals []Action, dir string) Returner
}

type Returner struct {
  Variables map[string]Variable
  Exp       Action
  Type      string
}
