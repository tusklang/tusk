package interpreter

import "reflect"
import "unicode"
import "strconv"

func interpreter(actions []Action, cli_params CliParams, vars map[string]Variable, expReturn bool, this_vals []Action, dir string) Returner {

  for _, v := range actions {

    switch v.Type {

      case "local":
        vars[v.Name] = Variable{
          Type: "local",
          Name: v.Name,
          Value: interpreter(v.ExpAct, cli_params, vars, true, this_vals, dir).Exp,
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
          Value: interpreter(v.ExpAct, cli_params, vars, true, this_vals, dir).Exp,
        }

        if expReturn {
          return Returner{
            Variables: vars,
            Exp: vars[v.Name].Value,
            Type: "expression",
          }
        }
      case "let":

        name := v.Name
        acts := v.ExpAct
        interpreted := interpreter(acts, cli_params, vars, true, this_vals, dir).Exp

        var variable Variable

        if len(v.Indexes) == 0 {
          if _, exists := vars[name]; !exists {
            variable = Variable{
              Type: "local",
              Name: name,
              Value: interpreted,
            }
          } else {
            variable = Variable{
              Type: vars[name].Type,
              Name: name,
              Value: interpreted,
            }
          }
        } else {

          if _, exists := vars[name]; !exists {
            variable = Variable{
              Type: "local",
              Name: name,
              Value: hash,
            }
          } else {
            variable = vars[name]
          }

          oMap := &variable.Value //a pointer (ref) to the variable value

          for _, sv := range v.Indexes {
            index := cast(interpreter(sv, cli_params, vars, true, this_vals, dir).Exp, "string").ExpStr

            if _, exists := (*oMap).Hash_Values[index]; !exists {
              (*oMap).Hash_Values[index] = []Action{ hash }
            }

            oMap = &((*oMap).Hash_Values[index][0])
          }

          *oMap = interpreted
        }

        vars[name] = variable

        if expReturn {
          return Returner{
            Variables: vars,
            Exp: variable.Value,
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
      case "function":

        name := v.Name

        if name != "" {
          vars[name] = Variable{
            Type: "function",
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
      case "fargc_number":
        var fargc uint64 = 0

        //determine how many variables were passed as params
        for _, v := range vars {
          if v.Type == "argument" {
            fargc++
          } else if v.Type == "pargv" {
            fargc+=uint64(len(interpreter([]Action{ v.Value }, cli_params, vars, true, this_vals, dir).Exp.Hash_Values))
          }
        }

        fargcOmmStr := emptyString
        fargcOmmStr.ExpStr = strconv.FormatUint(fargc, 10)

        fargcOmmNum := cast(fargcOmmStr, "number")

        if isEqual(fargcOmmNum, v) { //if the given fargc is equal to the count
          interpreted := interpreter(v.ExpAct, cli_params, vars, false, this_vals, dir)

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

        //do later

      case "fargc_paramlist":

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

        //if the given types are equal to the fargc list
        if reflect.DeepEqual(types, v.Params) {
          interpreted := interpreter(v.ExpAct, cli_params, vars, false, this_vals, dir)

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
        procCall := functionParser(v, cli_params, &vars, this_vals, dir, true, "#")

        if expReturn {
          return procCall
        }
      case "@":
        procCall := functionParser(v, cli_params, &vars, this_vals, dir, true, "@")

        if expReturn {
          return procCall
        }
      case "await":

        exp := interpreter(v.ExpAct, cli_params, vars, true, this_vals, dir).Exp
        thread := exp.Thread.WaitFor() //wait for the thread

        if expReturn {
          return Returner{
            Variables: vars,
            Exp: thread.Exp,
            Type: "expression",
          }
        }
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

      case "variableIndex":

        exp := interpreter(v.ExpAct, cli_params, vars, true, this_vals, dir).Exp
        index := indexesCalc(exp, v.Indexes, cli_params, vars, this_vals, dir)

        if expReturn {
          return Returner{
            Variables: vars,
            Exp: index,
            Type: "expression",
          }
        }

      case "add":

        first := interpreter(v.First, cli_params, vars, true, this_vals, dir).Exp
        second := interpreter(v.Second, cli_params, vars, true, this_vals, dir).Exp

        val := add(first, second, cli_params)

        if expReturn {
          return Returner{
            Variables: vars,
            Exp: val,
            Type: "expression",
          }
        }

      case "subtract":

        first := interpreter(v.First, cli_params, vars, true, this_vals, dir).Exp
        second := interpreter(v.Second, cli_params, vars, true, this_vals, dir).Exp

        val := subtract(first, second, cli_params)

        if expReturn {
          return Returner{
            Variables: vars,
            Exp: val,
            Type: "expression",
          }
        }

      case "multiply":

        first := interpreter(v.First, cli_params, vars, true, this_vals, dir).Exp
        second := interpreter(v.Second, cli_params, vars, true, this_vals, dir).Exp

        val := multiply(first, second, cli_params)

        if expReturn {
          return Returner{
            Variables: vars,
            Exp: val,
            Type: "expression",
          }
        }

      case "divide":

        first := interpreter(v.First, cli_params, vars, true, this_vals, dir).Exp
        second := interpreter(v.Second, cli_params, vars, true, this_vals, dir).Exp

        val := divide(first, second, cli_params)

        if expReturn {
          return Returner{
            Variables: vars,
            Exp: val,
            Type: "expression",
          }
        }

      case "modulo":

        first := interpreter(v.First, cli_params, vars, true, this_vals, dir).Exp
        second := interpreter(v.Second, cli_params, vars, true, this_vals, dir).Exp

        val := modulo(first, second, cli_params)

        if expReturn {
          return Returner{
            Variables: vars,
            Exp: val,
            Type: "expression",
          }
        }

      case "exponentiate":

        first := interpreter(v.First, cli_params, vars, true, this_vals, dir).Exp
        second := interpreter(v.Second, cli_params, vars, true, this_vals, dir).Exp

        val := exponentiate(first, second, cli_params)

        if expReturn {
          return Returner{
            Variables: vars,
            Exp: val,
            Type: "expression",
          }
        }

      case "less":

        val1 := interpreter(v.First, cli_params, vars, true, this_vals, dir).Exp
        val2 := interpreter(v.Second, cli_params, vars, true, this_vals, dir).Exp

        if expReturn {

          //only numeric values will be tested
          if val1.Type == "number" && val2.Type == "number" && isLess(val1, val2) {
            return Returner{
              Variables: vars,
              Exp: trueAct,
              Type: "expression",
            }
          } else {
            return Returner{
              Variables: vars,
              Exp: falseAct,
              Type: "expression",
            }
          }
        }

      case "greater":

        val1 := interpreter(v.First, cli_params, vars, true, this_vals, dir).Exp
        val2 := interpreter(v.Second, cli_params, vars, true, this_vals, dir).Exp

        if expReturn {

          //only numeric values will be tested
          if val1.Type == "number" && val2.Type == "number" && !isLessOrEqual(val1, val2) {
            return Returner{
              Variables: vars,
              Exp: trueAct,
              Type: "expression",
            }
          } else {
            return Returner{
              Variables: vars,
              Exp: falseAct,
              Type: "expression",
            }
          }
        }

      case "equals":

        val1 := interpreter(v.First, cli_params, vars, true, this_vals, dir).Exp
        val2 := interpreter(v.Second, cli_params, vars, true, this_vals, dir).Exp

        if expReturn {
          if equals(val1, val2) {
            return Returner{
              Variables: vars,
              Exp: trueAct,
              Type: "expression",
            }
          } else {
            return Returner{
              Variables: vars,
              Exp: falseAct,
              Type: "expression",
            }
          }

        }

      case "lessOrEqual":

        val1 := interpreter(v.First, cli_params, vars, true, this_vals, dir).Exp
        val2 := interpreter(v.Second, cli_params, vars, true, this_vals, dir).Exp

        if expReturn {

          //only numeric values will be tested
          if val1.Type == "number" && val2.Type == "number" && isLessOrEqual(val1, val2) {
            return Returner{
              Variables: vars,
              Exp: trueAct,
              Type: "expression",
            }
          } else {
            return Returner{
              Variables: vars,
              Exp: falseAct,
              Type: "expression",
            }
          }
        }

      case "greaterOrEqual":

        val1 := interpreter(v.First, cli_params, vars, true, this_vals, dir).Exp
        val2 := interpreter(v.Second, cli_params, vars, true, this_vals, dir).Exp

        if expReturn {

          //only numeric values will be tested
          if val1.Type == "number" && val2.Type == "number" && !isLess(val1, val2) {
            return Returner{
              Variables: vars,
              Exp: trueAct,
              Type: "expression",
            }
          } else {
            return Returner{
              Variables: vars,
              Exp: falseAct,
              Type: "expression",
            }
          }
        }

      case "notEqual":

        val1 := interpreter(v.First, cli_params, vars, true, this_vals, dir).Exp
        val2 := interpreter(v.Second, cli_params, vars, true, this_vals, dir).Exp

        if expReturn {
          if !equals(val1, val2) {
            return Returner{
              Variables: vars,
              Exp: trueAct,
              Type: "expression",
            }
          } else {
            return Returner{
              Variables: vars,
              Exp: falseAct,
              Type: "expression",
            }
          }
        }

      case "and":

        val1 := interpreter(v.First, cli_params, vars, true, this_vals, dir).Exp

        if expReturn {

          //if num1 is not truthy, the "and" is automatically false
          if !isTruthy(val1) {
            return Returner{
              Variables: vars,
              Exp: falseAct,
              Type: "expression",
            }
          }

          val2 := interpreter(v.Second, cli_params, vars, true, this_vals, dir).Exp
          if isTruthy(val2) {
            return Returner{
              Variables: vars,
              Exp: trueAct,
              Type: "expression",
            }
          } else {
            return Returner{
              Variables: vars,
              Exp: falseAct,
              Type: "expression",
            }
          }
        }

      case "or":

        val1 := interpreter(v.First, cli_params, vars, true, this_vals, dir).Exp

        if expReturn {

          //if num1 is truthy, the "or" is automatically true
          if isTruthy(val1) {
            return Returner{
              Variables: vars,
              Exp: trueAct,
              Type: "expression",
            }
          }

          val2 := interpreter(v.Second, cli_params, vars, true, this_vals, dir).Exp
          if isTruthy(val2) {
            return Returner{
              Variables: vars,
              Exp: trueAct,
              Type: "expression",
            }
          } else {
            return Returner{
              Variables: vars,
              Exp: falseAct,
              Type: "expression",
            }
          }
        }

      case "not":

        val := interpreter(v.First, cli_params, vars, true, this_vals, dir).Exp

        if expReturn {
          if isTruthy(val) {
            return Returner{
              Variables: vars,
              Exp: falseAct,
              Type: "expression",
            }
          } else {
            return Returner{
              Variables: vars,
              Exp: trueAct,
              Type: "expression",
            }
          }
        }

      case "nand":

        val1 := interpreter(v.First, cli_params, vars, true, this_vals, dir).Exp

        if expReturn {

          //if num1 is truthy, the "nor" is automatically false
          if !isTruthy(val1) {
            return Returner{
              Variables: vars,
              Exp: trueAct,
              Type: "expression",
            }
          }

          val2 := interpreter(v.Second, cli_params, vars, true, this_vals, dir).Exp
          if !isTruthy(val2) {
            return Returner{
              Variables: vars,
              Exp: falseAct,
              Type: "expression",
            }
          } else {
            return Returner{
              Variables: vars,
              Exp: trueAct,
              Type: "expression",
            }
          }
        }

      case "nor":

        val1 := interpreter(v.First, cli_params, vars, true, this_vals, dir).Exp

        if expReturn {

          //if num1 is not truthy, the "nand" is automatically true
          if !isTruthy(val1) {
            return Returner{
              Variables: vars,
              Exp: trueAct,
              Type: "expression",
            }
          }

          val2 := interpreter(v.Second, cli_params, vars, true, this_vals, dir).Exp
          if isTruthy(val2) {
            return Returner{
              Variables: vars,
              Exp: falseAct,
              Type: "expression",
            }
          } else {
            return Returner{
              Variables: vars,
              Exp: trueAct,
              Type: "expression",
            }
          }
        }

      case "xor":

        val1 := interpreter(v.First, cli_params, vars, true, this_vals, dir).Exp
        val2 := interpreter(v.Second, cli_params, vars, true, this_vals, dir).Exp

        if expReturn {

          //xor is the same as !=
          if isTruthy(val1) != isTruthy(val2) {
            return Returner{
              Variables: vars,
              Exp: falseAct,
              Type: "expression",
            }
          } else {
            return Returner{
              Variables: vars,
              Exp: trueAct,
              Type: "expression",
            }
          }
        }

      case "xnor":

        val1 := interpreter(v.First, cli_params, vars, true, this_vals, dir).Exp
        val2 := interpreter(v.Second, cli_params, vars, true, this_vals, dir).Exp

        if expReturn {

          //xnor is the same as ==
          if isTruthy(val1) == isTruthy(val2) {
            return Returner{
              Variables: vars,
              Exp: falseAct,
              Type: "expression",
            }
          } else {
            return Returner{
              Variables: vars,
              Exp: trueAct,
              Type: "expression",
            }
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
