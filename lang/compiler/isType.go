package compiler

import "regexp"

func isType(val string) bool {

  match, _ := regexp.MatchString("string|hash|number|boolean|falsey|array", val) //test if it is one of the type values

  return match
}
