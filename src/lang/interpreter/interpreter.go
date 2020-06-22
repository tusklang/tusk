package interpreter

import "reflect"
import "unicode"

func interpreter(actions []Action, cli_params CliParams, vars map[string]Variable, expReturn bool, this_vals []Action, dir string) Returner {

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
      case "let":
      case "alt":

        //if it looks like
        //alt(true);
        //instead of alt(true) {}
        //there is no alternator
        if len(v.Condition) == 0 {
          if expReturn {
            return Returner{
              Variables: vars,
              Exp: undef,
              Type: "expression",
            }
          }
        }

        curCond := 0

        //while its truthy
        for o := 0; isTruthy(interpreter(v.Condition[curCond].Condition, cli_params, vars, true, this_vals, dir).Exp); o++ {
          interpreted := interpreter(v.Condition[curCond].Actions, cli_params, vars, true, this_vals, dir)

          //pass the globals and already existing variables
          for _, sv := range interpreted.Variables {
            _, exists := vars[sv.Name]
            if sv.Type == "global" || exists {
              vars[sv.Name] = sv
            }
          }

          //check if they want to return/skip (continue)/break
          if interpreted.Type == "return" {
            return Returner{
              Variables: vars,
              Exp: interpreted.Exp,
              Type: "return",
            }
          }
          if interpreted.Type == "skip" {
            continue
          }
          if interpreted.Type == "break" {
            break
          }
        }
      case "log":
        log_format(interpreter(v.ExpAct, cli_params, vars, true, this_vals, dir).Exp, 2, true)
      case "print":
        log_format(interpreter(v.ExpAct, cli_params, vars, true, this_vals, dir).Exp, 2, false)
      case "expressionIndex":
        val := interpreter(v.ExpAct, cli_params, vars, true, this_vals, dir).Exp
        index := indexesCalc(val, v.Indexes, cli_params, vars, this_vals, dir)

        if expReturn {
          return Returner{
            Variables: vars,
            Exp: index,
            Type: "expression",
          }
        }
      case "group":
        interpreted := interpreter(v.ExpAct, cli_params, vars, false, this_vals, dir)

        for _, sv := range interpreted.Variables {
          _, exists := vars[sv.Name]
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
            Value: v,
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
        for _, v := range vars {
          if v.Type == "argument" {
            pargc++
          } else if v.Type == "pargv" {
            pargc+=uint64(len(interpreter([]Action{ v.Value }, cli_params, vars, true, this_vals, dir).Exp.Hash_Values))
          }
        }

        //do later

      case "pargc_paramlist":

        var types []string

        for _, sv := range vars {

          interpreted := interpreter([]Action{ sv.Value }, cli_params, vars, true, this_vals, dir).Exp

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

          for _, sv := range interpreted.Variables {
            _, exists := vars[sv.Name]
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
          if interpreted.Type == "break" {
            return Returner{
              Variables: vars,
              Exp: interpreted.Exp,
              Type: interpreted.Type,
            }
          }
          if interpreted.Type == "skip" {
            return Returner{
              Variables: vars,
              Exp: interpreted.Exp,
              Type: interpreted.Type,
            }
          }

        }

      case "#":
      case "@":
      case "conditional":

        for _, sv := range v.Condition {

          ////////////////////////////////////
          /*
          for example, if the dev writes:
          if (
            a: 1
            b: 2
            return a = b
          );
          it would not do the condition (because a != b)
          */
          expRetCond := len(sv.Condition) == 1
          ////////////////////////////////////

          val := interpreter(sv.Condition, cli_params, vars, expRetCond, this_vals, dir).Exp

          if isTruthy(val) || sv.Type == "else" {
            interpreted := interpreter(sv.Actions, cli_params, vars, true, this_vals, dir)

            for _, ssv := range interpreted.Variables {
              _, exists := vars[ssv.Name]
              if ssv.Type == "global" || exists {
                vars[ssv.Name] = ssv
              }
            }

            //if the dev wants to return/break/skip, the outer proc/loop
            if interpreted.Type == "return" {
              return Returner{
                Variables: vars,
                Exp: interpreted.Exp,
                Type: interpreted.Type,
              }
            }
            if interpreted.Type == "break" {
              break
            }
            if interpreted.Type == "skip" {
              continue
            }

            //dont test any more conditions if this condition was true
            break
          }
        }

      case "import":

        files := v.Value

        for _, sv := range files {
          interpreted := interpreter(sv, cli_params, vars, false, this_vals, dir)

          for _, ssv := range interpreted.Variables {
            if ssv.Type == "global" { //dont pass if it already exists because each file should keep its own local variables (security)
              vars[ssv.Name] = ssv
            }
          }
        }
      case "break":
        return Returner{
          Variables: vars,
          Exp: Action{},
          Type: "break",
        }
      case "skip":
        return Returner{
          Variables: vars,
          Exp: Action{},
          Type: "skip",
        }
      case "return":
        return Returner{
          Variables: vars,
          Exp: interpreter(v.ExpAct, cli_params, vars, true, this_vals, dir).Exp,
          Type: "break",
        }
      case "loop":

        cond := v.Condition[0].Condition
        expRetCond := len(v.Condition) == 1

        for ;isTruthy(interpreter(cond, cli_params, vars, expRetCond, this_vals, dir).Exp); {
          interpreted := interpreter(v.Condition[0].Actions, cli_params, vars, true, this_vals, dir)

          for _, sv := range interpreted.Variables {
            _, exists := vars[sv.Name]
            if sv.Type == "global" || exists {
              vars[sv.Name] = sv
            }
          }

          //if the dev wants to return/break/skip, the outer proc/loop
          if interpreted.Type == "return" {
            return Returner{
              Variables: vars,
              Exp: interpreted.Exp,
              Type: interpreted.Type,
            }
          }
          if interpreted.Type == "break" {
            return Returner{
              Variables: vars,
              Exp: interpreted.Exp,
              Type: interpreted.Type,
            }
          }
          if interpreted.Type == "skip" {
            return Returner{
              Variables: vars,
              Exp: interpreted.Exp,
              Type: interpreted.Type,
            }
          }
        }

      case "hash":

        val := v

        for k, sv := range v.Hash_Values {
          exp := interpreter(sv, cli_params, vars, true, this_vals, dir).Exp

          if !unicode.IsLower([]rune(k)[0]) { //if it starts with an uppercase letter
            exp.Access = "public"
          }

          val.Hash_Values[k] = []Action{ exp }
        }

        if expReturn {
          return Returner{
            Variables: vars,
            Exp: val,
            Type: "expression",
          }
        }

      case "array":

        val := v

        for k, sv := range v.Hash_Values {

          exp := interpreter(sv, cli_params, vars, true, this_vals, dir).Exp

          exp.Access = "public"

          val.Hash_Values[k] = []Action{ exp }
        }

        if expReturn {
          return Returner{
            Variables: vars,
            Exp: val,
            Type: "expression",
          }
        }

      case "hashIndex":

        var toInterpreter = hash

        toInterpreter.Hash_Values = v.Hash_Values

        index := indexesCalc(
          interpreter([]Action{ toInterpreter }, cli_params, vars, true, this_vals, dir).Exp,
          v.Indexes, cli_params, vars, this_vals, dir)

        if expReturn {
          return Returner{
            Variables: vars,
            Exp: interpreter([]Action{ index }, cli_params, vars, true, this_vals, dir).Exp,
            Type: "expression",
          }
        }

      case "arrayIndex":

        var toInterpreter = arr

        toInterpreter.Hash_Values = v.Hash_Values

        index := indexesCalc(
          interpreter([]Action{ toInterpreter }, cli_params, vars, true, this_vals, dir).Exp,
          v.Indexes, cli_params, vars, this_vals, dir)

        if expReturn {
          return Returner{
            Variables: vars,
            Exp: interpreter([]Action{ index }, cli_params, vars, true, this_vals, dir).Exp,
            Type: "expression",
          }
        }

      //basic value types
      case "string": fallthrough
      case "number": fallthrough
      case "boolean": fallthrough
      case "falsey": fallthrough
      case "thread": fallthrough
      case "none":

        if expReturn {
          return Returner{
            Variables: vars,
            Exp: v,
            Type: "expression",
          }
        }

      case "variable":

        var val Action

        if _, exists := vars[v.Name]; !exists {
          val = undef
        } else {
          val = interpreter([]Action{ vars[v.Name].Value }, cli_params, vars, true, this_vals, dir).Exp
        }

        if expReturn {
          return Returner{
            Variables: vars,
            Exp: val,
            Type: "expression",
          }
        }

    }

  }

  //if nothing was returned, return undef
  return Returner{
    Variables: vars,
    Exp: undef,
    Type: "none",
  }
}
