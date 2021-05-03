package ast

import "github.com/tusklang/tusk/parser/tokenizer"

type gItem struct {
	sub   []gItem
	token tokenizer.Token
}

var braceMatcher = map[string]string{
	"{": "}",
	"[": "]",
	"(": ")",

	"}": "{",
	"]": "[",
	")": "(",
}

func allZero(m map[string]int) bool {
	for _, v := range m {
		if v != 0 {
			return false
		}
	}
	return true
}

//grouper groups each parenthesis, brace, or curly brace block
func grouper(lex []tokenizer.Token) []gItem {

	braceAmt := map[string]int{
		"{": 0,
		"[": 0,
		"(": 0,
	}

	for _, v := range lex {

		for k := range braceAmt {
			if v.Name == k {
				braceAmt[k]++
			} else if v.Name == braceMatcher[k] {
				braceAmt[k]--
			}

			if braceAmt[k] < 0 {
				//compiler error
			}
		}

		if !allZero(braceAmt) {
			//a brace has been detected

		}
	}
}
