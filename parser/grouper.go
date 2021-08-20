package parser

import (
	"github.com/tusklang/tusk/tokenizer"
)

type gItem struct {
	Sub   []gItem
	Token tokenizer.Token
}

var braceMatcher = map[string]string{
	"{": "}",
	"[": "]",
	"(": ")",

	"}": "{",
	"]": "[",
	")": "(",
}

func nonZero(m map[string]int) string {
	for k, v := range m {
		if v != 0 {
			return k
		}
	}
	return ""
}

func detectBrace(v tokenizer.Token, braceAmt *map[string]int) {
	for k := range *braceAmt {
		if v.Type == k {
			(*braceAmt)[k]++
		} else if v.Type == braceMatcher[k] {
			(*braceAmt)[k]--
		}

		if (*braceAmt)[k] < 0 {
			//compiler error
		}
	}
}

//grouper groups each parenthesis, brace, or curly brace block
func grouper(lex []tokenizer.Token) []gItem {

	braceAmt := map[string]int{
		"{": 0,
		"[": 0,
		"(": 0,
	}

	var fin []gItem

	for i := 0; i < len(lex); i++ {

		//if its a whitespace token, skip it
		if lex[i].Type == "whitespace" {
			continue
		}

		detectBrace(lex[i], &braceAmt)

		if b := nonZero(braceAmt); braceAmt[b] != 0 {
			//a brace has been detected

			var sub []tokenizer.Token

			v := lex[i] //save the brace token

			for i++; i < len(lex); i++ {
				detectBrace(lex[i], &braceAmt)

				//a brace has closed it
				if braceAmt[b] == 0 {
					break
				}

				sub = append(sub, lex[i])
			}

			if braceAmt[b] != 0 {
				//compiler error (no closing brace)
			}

			fin = append(fin, gItem{
				Sub:   grouper(sub),
				Token: v,
			})

		} else {
			//regular token
			fin = append(fin, gItem{
				Sub:   nil,
				Token: lex[i],
			})
		}
	}

	return fin
}
