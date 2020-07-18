package goat

import "io/ioutil"
import "lang/compiler"
import "oat/helper"
import . "lang/types"

//export CompileFile
func CompileFile(filename string) (Oat, CompileErr) {
  f, _ := ioutil.ReadFile(filename)
  actions, vars, e := compiler.Compile(string(f), "goat compile")
  return Oat{ actions, vars }, e
}

//export CompileString
func CompileString(script string) (Oat, CompileErr) {
  actions, vars, e := compiler.Compile(script, "goat compile")
  return Oat{ actions, vars }, e
}

//export GetOat
func GetOat(file string) Oat {
  return oatHelper.FromOat(file)
}
