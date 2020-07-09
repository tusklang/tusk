package compiler

import "path"

// import . "lang/interpreter"

var included = []string{} //list of the imported files from omm

//export Run
func Run(params map[string]map[string]interface{}) {

  dir := params["Files"]["DIR"]
  fileName := params["Files"]["NAME"]

  included = append(included, path.Join(dir.(string), fileName.(string)))

  file := ReadFileJS(path.Join(dir.(string), fileName.(string)))[0]["Content"]

  _, _ = Compile(file, dir.(string), fileName.(string))

  // _, variables := Actionizer(lex, false, dir.(string), fileName.(string))

  //RunInterpreter(variables, params, dir.(string))
}
