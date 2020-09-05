package compiler

import (
	"strconv"

	. "github.com/omm-lang/omm/lang/types"
)

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

var curvar uint64

func changevarnames(actions []Action, newnames_ map[string]string) (map[string]string, error) {

	var e error

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

			for kk, vv := range fn.Overloads {
				for i, p := range vv.Params { //add the params to the current variables
					pname := "v " + strconv.FormatUint(curvar, 10)
					params[p] = pname
					fn.Overloads[kk].Params[i] = pname //also modify the parameters in the actual function
					curvar++
				}
				_, e = changevarnames(fn.Overloads[kk].Body, params)

				if e != nil {
					return nil, e
				}
			}

			actions[k].Value = fn
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

			keyandvalvars[key] = "v " + strconv.FormatUint(curvar, 10)
			v.First[1].Name = "v " + strconv.FormatUint(curvar, 10)
			curvar++
			keyandvalvars[val] = "v " + strconv.FormatUint(curvar, 10)
			v.First[2].Name = "v " + strconv.FormatUint(curvar, 10)
			curvar++

			tmp := []Action{v.First[0]}
			_, e = changevarnames(tmp, keyandvalvars)
			v.First[0] = tmp[0]
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

				var val = v.Value.(OmmProto).Static[i]

				var passvals = make(map[string]string)

				for k, v := range newnames {
					passvals[k] = v
				}

				var passarr = []Action{Action{
					Type:  (*val).Type(),
					Value: *val,
				}}
				_, e := changevarnames(passarr, passvals)
				*val = passarr[0].Value

				if e != nil {
					return nil, e
				}

				var tmp = actions[k].Value.(OmmProto)
				tmp.Static[i] = val
				actions[k].Value = tmp
			}

			var instanceproto = make(map[string]string)

			for k := range v.Value.(OmmProto).Instance {
				instanceproto[k] = k
				curvar++
			}

			for i := range v.Value.(OmmProto).Instance {

				var val = v.Value.(OmmProto).Instance[i]

				var passvals = make(map[string]string)

				for k, v := range instanceproto {
					passvals[k] = v
				}
				for k, v := range newnames {
					passvals[k] = v
				}

				var passarr = []Action{Action{
					Type:  (*val).Type(),
					Value: *val,
				}}
				_, e := changevarnames(passarr, passvals)
				*val = passarr[0].Value

				if e != nil {
					return nil, e
				}

				var tmp = actions[k].Value.(OmmProto)
				tmp.Instance[i] = val
				actions[k].Value = tmp
			}

			actions[k] = v

			continue
		}

		if v.Type == "var" || v.Type == "declare" {

			//check if it already exists
			if _, exists := newnames[v.Name]; exists {
				return nil, makeCompilerErr("Variable "+v.Name[1:]+" was already declared", v.File, v.Line)
			}

			newnames[v.Name] = "v " + strconv.FormatUint(curvar, 10)
			actions[k].Name = "v " + strconv.FormatUint(curvar, 10)
			curvar++
			_, e = changevarnames(actions[k].ExpAct, newnames)
			if e != nil {
				return nil, e
			}
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
			_, e = changevarnames(v.Hash[i][0], newnames)

			if e != nil {
				return nil, e
			}

			_, e = changevarnames(v.Hash[i][1], newnames)
			if e != nil {
				return nil, e
			}
		}
		////////////////////////////////////////////////

		/////////////////////////////////////////////

		if v.Type == "variable" || v.Type == "ovld" {
			if _, exists := newnames[v.Name]; !exists {
				return nil, makeCompilerErr("Variable "+v.Name[1:]+" was not declared", v.File, v.Line)
			}
			actions[k].Name = newnames[v.Name]
		}

	}

	return newnames, e
}
