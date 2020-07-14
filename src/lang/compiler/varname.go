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

func changevarnames(actions []Action, newnames_ map[string]string) {

  var newnames = make(map[string]string)

  //make newnames_ not get mutated
  for k, v := range newnames_ {
    newnames[k] = v
  }

  for k, v := range actions {

    if v.Type == "function" {

      var fn = v.Value.(OmmFunc)
      var params = newnames

      for i, p := range v.Value.(OmmFunc).Params { //add the params to the current variables
        pname := "v" + strconv.FormatUint(curvar, 10)
        fn.Params[i] = pname //also modify the parameters in the actual function
        params[p] = pname
        curvar++
      }
      v.Value = fn
      changevarnames(v.Value.(OmmFunc).Body, params)
      continue
    }
    if v.Type == "each" { //if it is each, also give the key and value variables
      key := v.First[1].Name
      val := v.First[2].Name

      var keyandvalvars = newnames
      keyandvalvars[key] = "v" + strconv.FormatUint(curvar, 10)
      curvar++
      keyandvalvars[val] = "v" + strconv.FormatUint(curvar, 10)
      curvar++

      changevarnames(v.ExpAct, keyandvalvars)
      continue
    }
    if v.Type == "proto" {

      for i := range v.Static {
        changevarnames(v.Static[i], newnames)
      }
      for i := range v.Instance {
        changevarnames(v.Instance[i], newnames)
      }

      continue
    }

    //perform checkvars on all of the sub actions
    changevarnames(v.ExpAct, newnames)
    changevarnames(v.First, newnames)
    changevarnames(v.Second, newnames)
    changevarnames(v.Degree, newnames)

    //also do it for the arrays and hashes
    for i := range v.Array {
      changevarnames(v.Array[i], newnames)
    }
    for i := range v.Hash {
      changevarnames(v.Hash[i], newnames)
    }
    //////////////////////////////////////

    /////////////////////////////////////////////

    if v.Type == "var" {
      newnames[v.Name] = "v" + strconv.FormatUint(curvar, 10)
      actions[k].Name = "v" + strconv.FormatUint(curvar, 10)
      curvar++
    }

    if v.Type == "variable" {
      if _, exists := newnames[v.Name]; !exists {
        compilerErr("Variable " + v.Name[1:] + " was not declared", v.File, v.Line)
      }
      actions[k].Name = newnames[v.Name]
    }

  }
}
