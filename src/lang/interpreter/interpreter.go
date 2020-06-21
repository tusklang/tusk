package interpreter

import "reflect"

import . "lang"

func interpreter(actions []Action, cli_params CliParams, map[string]Variable vars, expReturn bool, this_vals []Action, dir string) Returner {

  for _, v := range actions {

    switch v.Type {

      case "local":
        vars[v.Name] = Variable{
          Type: "local",
          Name: v.Name,
          Value: interpreter(v.ExpAct, cli_params, vars, expReturn, this_vals, dir).Exp,
        }

        if expReturn {
          return Returner{
            Variables: vars,
            Exp: vars[v.Name].Value,
            Type: "expression",
          }
        }
      case "dynamic":
        vars[v.Name] = Variable{
          Type: "dynamic",
          Name: v.Name,
          Value: v.ExpAct[0],
        }

        if expReturn {
          return Returner{
            Variables: vars,
            Exp: vars[v.Name].Value,
            Type: "expression",
          }
        }
      case "global":
        vars[v.Name] = Variable{
          Type: "global",
          Name: v.Name,
          Value: interpreter(v.ExpAct, cli_params, vars, expReturn, this_vals, dir).Exp,
        }

        if expReturn {
          return Returner{
            Variables: vars,
            Exp: vars[v.Name].Value,
            Type: "expression",
          }
        }
      case "alt":

        //if it looks like
        //alt(true);
        //instead of alt(true) {}
        //there is no alternator
        if len(v.Condition) == 0 {
          if expReturn {
            return undef
          }
        }

        curCond := 0

        //while its truthy
        for ;isTruthy(interpreter(v.Condition[curCond].Condition, cli_params, vars, true, this_vals, dir)); o++ {
          interpreted := interpreter(v.Condition[curCond].Actions, cli_params, vars, true, this_vals, dir)

          //pass the globals and already existing variables
          for _, sv := range interpreted.Variables {
            _, exists = vars[sv.Name]
            if sv.Type == "global" || exists {
              vars[sv.Name] = sv
            }
          }

          //check if they want to return/skip (continue)/break
          if (interpreted.Type == "return") return Returner{
            Variables: vars,
            Exp: interpreted.Exp,
            Type: "return"
          }
          if (interpreted.Type == "skip") continue
          if (interpreted.Type == "break") break
        }
      case "log":
        log_format(interpreter(v.ExpAct, cli_params, vars, true, this_vals, dir).exp, 2, true)
      case "print":
        log_format(interpreter(v.ExpAct, cli_params, vars, true, this_vals, dir).exp, 2, false)
      case "expressionIndex":
        val := interpreter(v.ExpAct, cli_params, vars, true, this_vals, dir)
        index := indexesCalc(val.Hash_Values, v.Indexes, cli_params, vars, this_vals, dir)

        if expReturn {
          return Returner{
            Variables: vars,
            Exp: index,
            Type: "expression",
          }
        }
      case "group":
        interpreted := interpreter(v.ExpAct, cli_params, vars, false, this_vals, dir)

        for k, sv := range interpreted.Variables {
          _, exists = vars[sv.Name]
          if sv.Type == "global" || exists {
            vars[sv.Name] = sv
          }
        }

        if expReturn {
          return Returner{
            Variables: vars,
            Exp: interpreted.Exp,
            Type: interpreted.Type,
          }
        }
      case "process":

        name := v.Name

        if name != "" {
          vars[name] = Variable{
            Type: "process",
            Name: name,
            Value: v
          }
        }

        if expReturn {
          return Returner{
            Variables: vars,
            Exp: v,
            Type: "expression",
          }
        }
      case "pargc_number":
        var pargc uint64 = 0

        //determine how many variables were passed as params
        for k, v := range vars {
          if v.Type == "argument" {
            pargc++
          } else if v.Type == "pargv" {
            pargc+=len(interpreter(v.Value, cli_params, vars, true, this_vals, dir).Exp.Hash_Values)
          }
        }

        //do later

      case "pargc_paramlist":

        var types []string

        for _, sv := range vars {

          interpreted := interpreter(sv, cli_params, vars, true, dir).Exp

          if sv.Type == "argument" {
            types = append(types, interpreted.Type)
          } else if sv.Type == "pargv" {
            for _, ssv := range interpreted.Hash_Values {
              types = append(types, interpreter(ssv, cli_params, vars, true, this_vals, dir).Exp.Type)
            }
          }
        }

        //if the given types are equal to the pargc list
        if reflect.DeepEqual(types, v.Params) {
          interpreted := interpreter(v.ExpAct, cli_params, vars, true, this_vals, dir)

          for _, sv := range interpreted.variables {
            _, exists = vars[sv.Name]
            if sv.Type == "global" || exists {
              vars[sv.Name] = sv
            }
          }

          if interpreted.Type == "return" {
            return Returner{
              Variables: vars,
              Exp: interpreted.Exp,
              Type: interpreted.Type,
            }
          }

        }

      case "#":
      case "@":
      case "return":
      case "if":

    }

  }

}
