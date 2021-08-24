package grouper

import "github.com/tusklang/tusk/tokenizer"

//util function to match braces and return everything in between
func braceMatcher(lex []tokenizer.Token, i *int, matchOpen string, matchClose string, removeTopBraces bool) []tokenizer.Token {

	var ret []tokenizer.Token

	//count of openers
	//this is increased if we locate an opener
	//this is decreased if we locate a closer
	cnt := 0

	//loop through the given lex
	for ; *i < len(lex); *i++ {

		ret = append(ret, lex[*i])

		if lex[*i].Type == matchOpen {
			cnt++
		}
		if lex[*i].Type == matchClose {
			cnt--
		}

		if cnt == 0 {
			break
		}
	}

	if removeTopBraces && len(ret) >= 2 {
		//remove the first and last braces if this option is given
		ret = ret[1 : len(ret)-1]
	}

	return ret
}
