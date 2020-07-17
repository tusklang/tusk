package interpreter

import "strconv"
import "oat/helper"
import . "lang/types"

//this file is for goat (go to omm/oat binding)

type Instance struct {
  Params   CliParams
  vars     map[string]*OmmVar
  globals  map[string]*OmmVar
}

func (ins *Instance) CallFunc(name string, args... *OmmType) *OmmType {

  _fn, exists := ins.globals["$" + name]
  fn := *_fn.Value

  if !exists {
    panic("Given global does not exists: " + name)
  }

  if fn.Type() != "function" {
    panic("Given global is not a function: " + name)
  }

  var body = fn.(OmmFunc).Body
  var params = fn.(OmmFunc).Params

  if len(params) != len(args) {
    panic("Function " + name + " requires " + strconv.Itoa(len(params)) + " arguments, but was given " + strconv.Itoa(len(args)) + " instead")
  }

  for k, v := range params {
    ins.vars[v] = &OmmVar{
      Name: v,
      Value: args[k],
    }
  }

  return ins.interpreter(body, []string{ "at goat caller" }).Exp
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
}
