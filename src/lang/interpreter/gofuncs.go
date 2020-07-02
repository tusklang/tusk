package interpreter

//all of the gofuncs
//functions written in go that are used by omm

import "strconv"
import "strings"
import "time"
import "os/exec"

var gofuncs = map[string]func(args [][]Action, cli_params CliParams, vars map[string]Variable, this_vals []Action, dir string) Action {

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
    }

    return undef
  },
  "clone": func(args [][]Action, cli_params CliParams, vars map[string]Variable, this_vals []Action, dir string) Action {

    //clone will take a value and create a "non refrence"
    //this is because go is weird, and omm is based on go
    //example:
    //  value1: [::]
    //  value2: value1
    //
    //  value2::SubValue: 'Hi There'
    //  log value1
    //
    //will log
    //  [:
    //    SubValue: 'Hi There'
    //  :]
    //instead of
    //  [::]
    //
    //clone() will just allow you to get the latter value

    if len(args) == 1 {
      val := interpreter(args[0], cli_params, vars, true, this_vals, dir).Exp

      //copy the elements of the hash
      var hash = make(map[string][]Action)
      for k, v := range val.Hash_Values {
        hash[k] = v
      }

      rval := Action{
        Type: val.Type,
        Name: val.Name,
        ExpStr: val.ExpStr,
        ExpAct: append([]Action{}, val.ExpAct...),
        Params: append([]string{}, val.Params...),
        Args: append([][]Action{}, val.Args...),
        Condition: append([]Condition{}, val.Condition...),
        First: append([]Action{}, val.First...),
        Second: append([]Action{}, val.Second...),
        Degree: append([]Action{}, val.Degree...),
        Value: append([][]Action{}, val.Value...),
        Indexes: append([][]Action{}, val.Indexes...),
        Hash_Values: hash,
        Access: val.Access,
        Integer: append([]int64{}, val.Integer...),
        Decimal: append([]int64{}, val.Decimal...),
        Thread: val.Thread,
      }

      return rval
    }

    return undef
  },
  "exec": func(args [][]Action, cli_params CliParams, vars map[string]Variable, this_vals []Action, dir string) Action {

    if len(args) == 1 {
      cmd := interpreter(args[0], cli_params, vars, true, this_vals, dir).Exp.ExpStr

      command := exec.Command(cmd)
      command.Dir = dir
      out, _ := command.CombinedOutput()

      stringValue := emptyString
      stringValue.ExpStr = string(out)
      return stringValue
    } else if len(args) == 2 {
      cmd := interpreter(args[0], cli_params, vars, true, this_vals, dir).Exp.ExpStr
      stdin := interpreter(args[1], cli_params, vars, true, this_vals, dir).Exp.ExpStr

      command := exec.Command(cmd)

      command.Dir = dir
      command.Stdin = strings.NewReader(stdin)

      out, _ := command.CombinedOutput()

      stringValue := emptyString
      stringValue.ExpStr = string(out)
      return stringValue
    }

    return undef
  },

}
