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
      v.Value = fn
      _, e = changevarnames(v.Value.(OmmFunc).Body, params)
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

      for i := range v.Value.(OmmProto).Static {

        var val = *v.Value.(OmmProto).Static[i]

        if val.Type() == "function" {

          var fn = val.(OmmFunc)

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
          if e != nil {
            return nil, e
          }
        }

      }

      var instanceproto map[string]string
      for i := range v.Value.(OmmProto).Instance {

        var val = *v.Value.(OmmProto).Static[i]

        if val.Type() == "function" {

          var fn = val.(OmmFunc)

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

          protonames, e := changevarnames(fn.Body, newnames)

          /* Image we have this proto

              proto {
                instance var self
                var init: fn(self) {
                  self::self: self
                }
                var testf: fn() {
                  log self ; without the following, this would cause an error
                }
              }
          */
          newnames = protonames
          instanceproto = protonames
          if e != nil {
            return nil, e
          }
        }

      }

      for k := range instanceproto {
        delete(newnames, k) //prevent outside of the proto from using proto variables
      }

      continue
    }

    //perform checkvars on all of the sub actions
    _, e = changevarnames(v.ExpAct, newnames)
    if e != nil {
      return nil, e
    }
    _, e = changevarnames(v.First, newnames)
    if e != nil {
      return nil, e
    }
    _, e = changevarnames(v.Second, newnames)
    if e != nil {
      return nil, e
    }
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
