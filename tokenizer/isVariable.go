package tokenizer

import "github.com/dlclark/regexp2"

//IsVariable tests if a token is a variable or not
func IsVariable(t Token) bool {

	//go through all other tokens
	for _, v := range tokenlist {

		//skip the pattern for variable itself (otherwise it'd match and return false)
		if v == VarPat {
			continue
		}

		re := regexp2.MustCompile(v, 0)

		//if it matches with a pattern for a non-variable, then its not a variable
		if b, e := re.MatchString(t.Name); b && e == nil {
			return false
		}
	}

	return true
}
