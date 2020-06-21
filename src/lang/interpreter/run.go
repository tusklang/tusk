package interpreter

type CliParams map[string]map[string]interface{}

//export RunInterpreter
func RunInterpreter(actions []Action, cli_params map[string]map[string]interface{}, dir string) {
  var vars map[string]Variable
  interpreter(actions, CliParams(cli_params), vars, false, []Action{}, dir)
}
