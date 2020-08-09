package interpreter

import "oat/helper"
import "path/filepath"
import "path"
import "os"
import . "lang/types"

//this file is for goat (go to omm/oat binding)

type Instance struct {
  Params   CliParams
  vars     map[string]*OmmVar
  globals  map[string]*OmmVar
}

func (ins *Instance) HasGlobal(name string) bool {
  _, exists := ins.globals["$" + name]
  return exists
}

func (ins *Instance) GetGlobal(name string) *OmmType {
  variable, exists := ins.globals["$" + name]

  if !exists {
    panic("Given global does not exists: " + name)
  }

  return variable.Value
}

func (ins *Instance) SetLcl(name string, value *OmmType) {
  ins.vars[name] = &OmmVar{
    Name: name,
    Value: value,
  }
}

func (ins *Instance) CallFunc(name string, args... *OmmType) *OmmType {

  _fn, exists := ins.globals["$" + name]

  if !exists {
    panic("Given global does not exists: " + name)
  }

  fn := *_fn.Value

  if fn.Type() != "function" {
    panic("Given global is not a function: " + name)
  }

  var argarr = OmmArray{
    Array: args,
    Length: uint64(len(args)),
  }

  return Operations["function <- array"](fn, argarr, ins, []string{ "at goat caller" }, 0, path.Join(ins.Params.Directory, ins.Params.Name))
}

func (ins *Instance) FromOat(filename string) {
  var decoded = oatHelper.FromOat(filename)

  var interpreted = make(map[string]*OmmVar)

  for k, v := range decoded.Variables {
    interpreted[k] = &OmmVar{
      Name: k,
      Value: ins.interpreter(v, []string{ "at goat usage" }).Exp,
    }
  }

  ins.globals = interpreted
  ins.vars = interpreted

  var dirnameOmmStr OmmString
  __d, _ := filepath.Abs(filepath.Dir(os.Args[0]))
  dirnameOmmStr.FromGoType(__d)
  var dirnameOmmType OmmType = dirnameOmmStr

  ins.vars["$__dirname"] = &OmmVar{ //put the dirname
    Name: "$__dirname",
    Value: &dirnameOmmType,
  }

  //also account for the GoFuncs
  for k, v := range GoFuncs {
    var gofunc OmmType = OmmGoFunc{
      Function: v,
    }
    ins.vars["$" + k] = &OmmVar{
      Name: "$" + k,
      Value: &gofunc,
    }
  }
}

func (ins *Instance) FromOatStruct(oat Oat) {
  var interpreted = make(map[string]*OmmVar)

  for k, v := range oat.Variables {
    interpreted[k] = &OmmVar{
      Name: k,
      Value: ins.interpreter(v, []string{ "at goat usage" }).Exp,
    }
  }

  ins.globals = interpreted
  ins.vars = interpreted

  var dirnameOmmStr OmmString
  __d, _ := filepath.Abs(filepath.Dir(os.Args[0]))
  dirnameOmmStr.FromGoType(__d)
  var dirnameOmmType OmmType = dirnameOmmStr

  ins.vars["$__dirname"] = &OmmVar{ //put the dirname
    Name: "$__dirname",
    Value: &dirnameOmmType,
  }

  //also account for the GoFuncs
  for k, v := range GoFuncs {
    var gofunc OmmType = OmmGoFunc{
      Function: v,
    }
    ins.vars["$" + k] = &OmmVar{
      Name: "$" + k,
      Value: &gofunc,
    }
  }
}
