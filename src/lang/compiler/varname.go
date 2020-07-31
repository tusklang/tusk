package compiler

import "strconv"

import . "lang/types"

/*
Explanation of why this exists:

  In omm (and java/c/c++ I think), there is one list of variables at compile time. For security, some variables cannot be accessed from different scopes, and it will
  give a **compile time** error. In a lot of dynamic languages, like javascript and python, "variable not existsing" errors happen at runtime. If you run the
  following python code:

    print("First")
    print(testvariable)

  it will print "First" even though testvariable is not defined. It will still give an error, but it will print "First" first. In c++ this code

    #include <iostream>

    int main() {
      std::cout << "First" << std::endl;
      std::cout << testvariable << std::endl;
    }

  a compiler error would come. Now I will explain why we need to change the variable names in the compiler. Imagine we have some omm code like this:

    var main: fn() => {
      var test: 1
      testf async []
      test: 3
    }

    var testf: fn() => {
      var test: 2
      var i: 0
      while (i < 100) => i: i + 1;
      log test
    }

  It would give 3, not 2 because testf and main have the same variable set. This can open a security leak, so that is why this file exists.

  Also this file kinda kills two birds with one stone because it also serves as the variable existance checker (which prevents a user from refrencing a variable that was not yet declared)
  most statically typed languages do this, but the dynamic languages dont. Omm is dynamic, but acts like a static language in this regard.

*/

var curvar uint64 = 0

func changevarnames(actions []Action, newnames_ map[string]string) (map[string]string, CompileErr) {

  var e CompileErr

  newnames_["$__dirname"] = "$__dirname" //__dirname is a global

  var newnames = make(map[string]string)

  //make newnames_ not get mutated
  for k, v := range newnames_ {
    newnames[k] = v
  }

  for k, v := range actions {

    if v.Type == "function" {

      var fn = v.Value.(OmmFunc)

      var params = make(map[string]string)

      //clone `newnames` into `params` which is a list of variables that are passed into the function
      for k, v := range newnames {
        params[k] = v
      }

      for i, p := range v.Value.(OmmFunc).Params { //add the params to the current variables
        pname := "v" + strconv.FormatUint(curvar, 10)
        fn.Params[i] = pname //also modify the parameters in the actual function
        params[p] = pname
        curvar++
      }
      _, e = changevarnames(fn.Body, params)
      actions[k].Value = fn
      if e != nil {
        return nil, e
      }

      continue
    }
    if v.Type == "each" { //if it is each, also give the key and value variables
      key := v.First[1].Name
      val := v.First[2].Name

      var keyandvalvars = make(map[string]string)

      //clone newnames into keyandvalvars
      for key, val := range newnames {
        keyandvalvars[key] = val
      }
      ///////////////////////////////////

      keyandvalvars[key] = "v" + strconv.FormatUint(curvar, 10)
      v.First[1].Name = "v" + strconv.FormatUint(curvar, 10)
      curvar++
      keyandvalvars[val] = "v" + strconv.FormatUint(curvar, 10)
      v.First[2].Name = "v" + strconv.FormatUint(curvar, 10)
      curvar++

      _, e = changevarnames([]Action{ v.First[0] }, keyandvalvars)
      if e != nil {
        return nil, e
      }
      _, e = changevarnames(v.ExpAct, keyandvalvars)
      if e != nil {
        return nil, e
      }
      continue
    }
    if v.Type == "proto" {

      for i := range v.Static {
        var val = v.Static[i][0]

        var passvals = make(map[string]string)

        for k, v := range newnames {
          passvals[k] = v
        }

        var passarr = []Action{ val }
        changevarnames(passarr, passvals)
        actions[k].Static[i] = passarr
      }

      var instanceproto = make(map[string]string)

      for _, val := range v.Instance {
        instanceproto[val[0].Name] = "v" + strconv.FormatUint(curvar, 10)
        curvar++
      }

      for i := range v.Instance {
        var val = v.Instance[i][0]

        var passvals = make(map[string]string)

        for k, v := range instanceproto {
          passvals[k] = v
        }
        for k, v := range newnames {
          passvals[k] = v
        }

        var passarr = []Action{ val }
        changevarnames(passarr, passvals)
        actions[k].Instance[i] = passarr
      }

      for k := range instanceproto {
        delete(newnames, k) //prevent outside of the proto from using proto variables
      }

      actions[k] = v

      continue
    }

    //perform checkvars on all of the sub actions
    _, e = changevarnames(actions[k].ExpAct, newnames)
    if e != nil {
      return nil, e
    }
    _, e = changevarnames(actions[k].First, newnames)
    if e != nil {
      return nil, e
    }
    _, e = changevarnames(actions[k].Second, newnames)
    if e != nil {
      return nil, e
    }

    //also do it for the (runtime) arrays and hashes
    for i := range v.Array {
      _, e = changevarnames(v.Array[i], newnames)
      if e != nil {
        return nil, e
      }
    }
    for i := range v.Hash {
      _, e = changevarnames(v.Hash[i], newnames)
      if e != nil {
        return nil, e
      }
    }
    ////////////////////////////////////////////////

    /////////////////////////////////////////////

    if v.Type == "var" || v.Type == "declare" {
      newnames[v.Name] = "v" + strconv.FormatUint(curvar, 10)
      actions[k].Name = "v" + strconv.FormatUint(curvar, 10)
      curvar++
    }

    if v.Type == "variable" {
      if _, exists := newnames[v.Name]; !exists {
        return nil, makeCompilerErr("Variable " + v.Name[1:] + " was not declared", v.File, v.Line)
      }
      actions[k].Name = newnames[v.Name]
    }

  }

  return newnames, e
}
