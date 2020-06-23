package interpreter

import "strconv"
import "strings"

var threads []chan Returner

func processParser(v Action, cli_params CliParams, vars *map[string]Variable, this_vals []Action, dir string, isProc bool, callType string) Returner {

  name := v.Name
  var parsed Returner = Returner{
    Variables: (*vars),
    Exp: undef,
    Type: "expression",
  }

  if _, exists := (*vars)[name]; !exists {
    goto end
  } else {

    //if it is a goproc
    if (*vars)[name].Type == "goproc" {
      parsed = (*vars)[name].GoProc(v, cli_params, *vars, this_vals, dir)
    } else {

      variable := (*vars)[name].Value
      send_this := this_vals

      for _, sv := range v.Indexes {

        //convert index to string
        index := cast(interpreter(sv, cli_params, *vars, true, this_vals, dir).Exp, "string").ExpStr

        if _, exists = variable.Hash_Values[index]; !exists || variable.Hash_Values[index][0].Access == "private" {
          goto end
        }

        send_this = append([]Action{ variable }, send_this...)
        variable = interpreter(variable.Hash_Values[index], cli_params, *vars, true, this_vals, dir).Exp
      }

      if !isProc {
        parsed = Returner{
          Variables: *vars,
          Exp: variable,
          Type: "expression",
        }
        goto end
      }

      if variable.Type != "process" {
        goto end
      }

      params := variable.Params
      arguments := v.Args

      //ensure that there are the same amount of params and args
      if len(params) != len(arguments) && !containsParams(params, "$pargv") {
        goto end
      }

      sendVars := *vars

      for k, param_v := range params {
        if strings.HasPrefix(param_v, "$pargv") {
          varname := "$" + strings.TrimPrefix(param_v, "$pargv")

          //convert the rest of the args into a pargv
          var pargv = make(map[string][]Action)

          for cur, o := 0, k; o < len(arguments); cur, o = cur + 1, o + 1 {
            pargv[strconv.Itoa(cur)] = []Action{ interpreter(arguments[o], cli_params, *vars, true, this_vals, dir).Exp }
          }

          arg := arr
          arg.Hash_Values = pargv

          sendVars[varname] = Variable{
            Type: "pargv",
            Name: varname,
            Value: arg,
          }

          break
        }

        sendVars[param_v] = Variable{
          Type: "argument",
          Name: param_v,
          Value: interpreter(arguments[k], cli_params, *vars, true, this_vals, dir).Exp,
        }
      }

      if callType == "#" {
        parsed = interpreter(variable.ExpAct, cli_params, sendVars, true, send_this, dir)

        for _, sv := range parsed.Variables {
          _, exists := (*vars)[sv.Name]
          if sv.Type == "global" || exists {
            (*vars)[sv.Name] = sv
          }
        }

        for _, sv := range v.SubCall {
          curVar := parsed.Exp

          curVar.Indexes = sv.Indexes
          curVar.Args = sv.Args

          parsed = processParser(curVar, cli_params, vars, this_vals, dir, sv.IsProc, "#")
          goto end
        }
      } else if callType == "@" {

        retChannel := make(chan Returner)
        threads = append(threads, retChannel)
        go func() {
          retChannel <- interpreter(variable.ExpAct, cli_params, sendVars, true, send_this, dir)
        }()

        parsed.Exp = thread
        parsed.Exp.Thread = retChannel
      }

    }

  }

  end:
  return parsed
}
