package ast

import (
	"github.com/tusklang/tusk/tokenizer"
)

//util function to match braces and return everything in between
//this implementation is kinda jank soo uh, refactor soon:tm:
func braceMatcher(lex []tokenizer.Token, i *int, matchOpen []string, matchClose []string, removeTopBraces bool, stopAt string) []tokenizer.Token {

	var ret []tokenizer.Token

	//count of openers
	//this is increased if we locate an opener
	//this is decreased if we locate a closer
	var cnt = make([]int, len(matchOpen))

	//loop through the given lex
	for ; *i < len(lex); *i++ {

		ret = append(ret, lex[*i])

		for k, v := range matchOpen {
			if lex[*i].Type == v {
				cnt[k]++
			}
			if lex[*i].Type == matchClose[k] {
				cnt[k]--
			}
		}

		for _, v := range cnt {
			if v != 0 {
				goto skip
			}
		}

		if stopAt == "" {
			break
		}

		//include the `stopAt` token in the final result
		if lex[*i].Type == stopAt {
			ret = append(ret, lex[*i])
			*i++
			break
		}

	skip:
	}

	if removeTopBraces && len(ret) >= 2 {
		//remove the first and last braces if this option is given
		ret = ret[1 : len(ret)-1]
	}

	if stopAt != "" {
		//remove the last value that we stopped at
		*i--
		ret = ret[:len(ret)-1]
	}

	return ret
}
