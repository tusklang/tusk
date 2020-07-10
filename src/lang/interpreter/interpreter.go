package interpreter

import "fmt"
import . "lang/types"

func interpreter(actions []Action, cli_params CliParams, passedVars map[string]Variable, this_vals []Action) Returner {

  var expReturn = false //if it is inside an expression

  if len(actions) == 1 {
    expReturn = true
  }

  var vars = make(map[string]Variable)

  for k, v := range passedVars { //copy passedVars into vars
    vars[k] = v
  }

  for _, v := range actions {
    switch v.Type {

      case "local":
        vars[v.Name] = Variable{
          Type: "local",
          Value: interpreter(v.ExpAct, cli_params, vars, this_vals).Exp,
        }

        if expReturn {
          return Returner{
            Type: "expression",
            Exp: vars[v.Name].Value,
            Variables: vars,
          }
        }

      case "global":
        vars[v.Name] = Variable{
          Type: "global",
          Value: interpreter(v.ExpAct, cli_params, vars, this_vals).Exp,
        }

        if expReturn {
          return Returner{
            Type: "expression",
            Exp: vars[v.Name].Value,
            Variables: vars,
          }
        }

      case "log":
        fmt.Println(interpreter(v.ExpAct, cli_params, vars, this_vals).Exp.Format())
      case "print":
        fmt.Print(interpreter(v.ExpAct, cli_params, vars, this_vals).Exp.Format())

      //all of the types
      case "string": fallthrough
      case "rune": fallthrough
      case "number": fallthrough
      case "bool": fallthrough
      case "function": fallthrough
      case "array": fallthrough
      case "hash": fallthrough
      case "undef": fallthrough
      case "thread":

        if expReturn {
          return Returner{
            Type: "expression",
            Exp: v.Value,
            Variables: vars,
          }
        }
      //////////////////

      case "{":

        groupRet := interpreter(v.ExpAct, cli_params, vars, this_vals)

        for k, v := range groupRet.Variables {
          if _, exists := vars[k]; exists || v.Type == "global" {
            vars[k] = v
          }
        }

        if expReturn {
          return Returner{
            Type: "expression",
            Exp: groupRet.Exp,
            Variables: groupRet.Variables,
          }
        }

      case "(":

        groupRet := interpreter(v.ExpAct, cli_params, vars, this_vals)

        for k, v := range groupRet.Variables {
          vars[k] = v
        }

        if expReturn {
          return Returner{
            Type: "expression",
            Exp: groupRet.Exp,
            Variables: vars,
          }
        }

    }
  }

  return Returner{
    Type: "none",
    Variables: vars,
  }
}
