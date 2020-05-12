package ocs

//omm cementing script (ocs pronounced ox)
//used to call other languages from omm
//meant to "cement" all of the languages together

//work on this in the future

import "lang"

type OCSAction struct {
  lang.Action
}

//export Compile
func Compile(file, dir, filename string) []lang.Action {

  lex := lang.Lexer(file, dir, filename)
  acts := lang.Actionizer(lex)

  return acts
}

func getOCS(acts []lang.Action) []OCSAction {

  for _, v := range acts {


  }
}

//export Run
func Run(filename, dir string) {

  file := lang.ReadFileJS(dir + filename)[0]
  actions := Compile(file, dir, filename)


}
