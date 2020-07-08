package interpreter

//all of the gofuncs
//functions written in go that are used by omm

import "strconv"
import "strings"
import "time"
import "os"
import "os/exec"
import "io/ioutil"

import . "lang/types"

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
      var hash = make(map[string]Action)
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
      cmdstr := interpreter(args[0], cli_params, vars, true, this_vals, dir).Exp.ExpStr

      command := exec.Command(cmdstr)
      out, _ := command.CombinedOutput()

      stringValue := emptyString
      stringValue.ExpStr = string(out)
      return stringValue
    } else if len(args) == 2 {
      cmdstr := interpreter(args[0], cli_params, vars, true, this_vals, dir).Exp.ExpStr
      stdin := interpreter(args[1], cli_params, vars, true, this_vals, dir).Exp.ExpStr

      command := exec.Command(cmdstr)
      command.Stdin = strings.NewReader(stdin)

      out, _ := command.CombinedOutput()

      stringValue := emptyString
      stringValue.ExpStr = string(out)
      return stringValue
    }

    return undef
  },
  "chdir": func(args [][]Action, cli_params CliParams, vars map[string]Variable, this_vals []Action, dir string) Action {

    if len(args) == 1 {
      //get the expstr of the first argument
      putDir := cast(interpreter(args[0], cli_params, vars, true, this_vals, dir).Exp, "string").ExpStr

      os.Chdir(putDir)
    }

    return undef
  },

  //filesystem gofuncs
  "files.isDir": func(args [][]Action, cli_params CliParams, vars map[string]Variable, this_vals []Action, dir string) Action {

    if len(args) == 1 {
      //get the expstr of the first argument
      name := cast(interpreter(args[0], cli_params, vars, true, this_vals, dir).Exp, "string").ExpStr

      stat, e := os.Stat(name)

      if e != nil { //if there is an error, then it isn't a directory
        return falseAct
      }

      switch mode := stat.Mode(); {
        case mode.IsDir():
          return trueAct
        default:
          return falseAct
      }
    }

    return undef
  },
  "files.isFile": func(args [][]Action, cli_params CliParams, vars map[string]Variable, this_vals []Action, dir string) Action {

    if len(args) == 1 {
      //get the expstr of the first argument
      name := cast(interpreter(args[0], cli_params, vars, true, this_vals, dir).Exp, "string").ExpStr

      stat, e := os.Stat(name)

      if e != nil { //if there is an error, then it isn't a file
        return falseAct
      }

      switch mode := stat.Mode(); {
        case mode.IsRegular(): //IsRegular means it is a file
          return trueAct
        default:
          return falseAct
      }
    }

    return undef
  },
  "files.exists": func(args [][]Action, cli_params CliParams, vars map[string]Variable, this_vals []Action, dir string) Action {

    if len(args) == 1 {
      //get the expstr of the first argument
      name := cast(interpreter(args[0], cli_params, vars, true, this_vals, dir).Exp, "string").ExpStr

      _, e := os.Stat(name)

      if e != nil { //if there is an error, then it does not exist
        return falseAct
      } else {
        return trueAct //otherwise it does exist
      }
    }

    return undef
  },
  "files.readFile": func(args [][]Action, cli_params CliParams, vars map[string]Variable, this_vals []Action, dir string) Action {

    if len(args) == 1 {
      //get the expstr of the first argument
      name := cast(interpreter(args[0], cli_params, vars, true, this_vals, dir).Exp, "string").ExpStr

      _, e := os.Stat(name)

      if e != nil { //if there is an error, then it does not exist
        return undef
      }

      fileByte, err := ioutil.ReadFile(name)

      if err != nil {
        return undef //if there was an error, return undef
      }

      var retStr = emptyString
      retStr.ExpStr = string(fileByte) //give the expstr

      for k, v := range fileByte {

        runeP := emptyRune
        runeP.ExpStr = string(v)

        retStr.Hash_Values[strconv.Itoa(k)] = runeP
      }

      return retStr
    }

    return undef
  },
  "files.readDir": func(args [][]Action, cli_params CliParams, vars map[string]Variable, this_vals []Action, dir string) Action {

    if len(args) == 1 {
      //get the expstr of the first argument
      name := cast(interpreter(args[0], cli_params, vars, true, this_vals, dir).Exp, "string").ExpStr

      _, e := os.Stat(name)

      if e != nil { //if there is an error, then it does not exist
        return undef
      }

      dirstats, err := ioutil.ReadDir(name)

      if err != nil {
        return undef //if there was an error, return undef
      }

      var ommarray = arr

      idx := zero

      for _, v :=  range dirstats {
        statStr := emptyString
        statStr.ExpStr = v.Name()

        ommarray.Hash_Values[cast(idx, "string").ExpStr] = statStr
        idx = number__plus__number(idx, one, cli_params)
      }

      return ommarray
    }

    return undef
  },
  "files.write": func(args [][]Action, cli_params CliParams, vars map[string]Variable, this_vals []Action, dir string) Action {

    if len(args) == 2 {
      //get the filename and contents
      name := cast(interpreter(args[0], cli_params, vars, true, this_vals, dir).Exp, "string").ExpStr //filename
      contents := cast(interpreter(args[1], cli_params, vars, true, this_vals, dir).Exp, "string") //content

      err := ioutil.WriteFile(name, []byte(contents.ExpStr), 0644)

      if err != nil {
        return undef
      }
      return contents
    }

    return undef
  },
  ///////////////////

}
