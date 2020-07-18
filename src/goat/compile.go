package goat

import "io/ioutil"
import "lang/compiler"
import "oat/helper"
import . "lang/types"

//export CompileFile
func CompileFile(filename string) Oat {
  f, _ := ioutil.ReadFile(filename)
  actions, vars := compiler.Compile(string(f), "goat compile")
  return Oat{ actions, vars }
}

//export CompileString
func CompileString(script string) Oat {
  actions, vars := compiler.Compile(script, "goat compile")
  return Oat{ actions, vars }
}

//export GetOat
func GetOat(file string) Oat {
  return oatHelper.FromOat(file)
}
