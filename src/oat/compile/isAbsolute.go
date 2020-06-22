package oatCompile

import "regexp"

//export IsAbsolute
func IsAbsolute(dir string) bool {

  match, _ := regexp.MatchString("^[a-zA-Z]:", dir)

  return match
}
