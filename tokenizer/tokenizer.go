package tokenizer

import (
	"strings"

	"github.com/dlclark/regexp2"
)

func Tokenizer(data, filenam string) (tokens []Token) {

	c := 0 //current index in the file string

	var row, col = 1, 1 //current row and column in the file

	for c < len(data) {

		for _, v := range tokenlist {

			//find all matches (we have lookbehinds so we can't just chop off the string's first half and search after that)

			var vreg = regexp2.MustCompile(v.regexp, 0)
			m1, _ := vreg.FindStringMatch(data)
			var matches = []*regexp2.Match{m1}

			for matches[len(matches)-1] != nil {
				m, _ := vreg.FindNextMatch(matches[len(matches)-1])
				matches = append(matches, m)
			}
			var match *regexp2.Match

			for _, m := range matches {
				if m != nil && m.Group.Capture.Index == c {
					match = m
					break
				}
			}

			if match != nil {
				//the current token is matched
				matched := match.Group.Capture.String()

				tokens = append(tokens, Token{
					Name:    matched,
					Type:    v.tokentype,
					File:    filenam,
					Snippet: strings.Split(data, "\n")[row-1],
					Row:     row,
					Col:     col,
				})

				if matched == "\n" {
					row++
					col = 1
				}

				mlen := len(matched)

				c += mlen
				col += mlen
				continue
			}
		}
	}

	//remove the whitespace tokens and replace varnames with keywords if they are
	var wsRem []Token

	for _, v := range tokens {

		if v.Type == "whitespace" || v.Type == "newline" {
			continue
		}

		if v.Type == "varname" && (func() bool {
			for _, vv := range keywords {
				if vv == v.Name {
					return true
				}
			}
			return false
		}()) {
			//it's a keyword
			//set the type to the name, as a keyword's type is it's name
			//e.g. var's type is var
			v.Type = v.Name
		}

		wsRem = append(wsRem, v)
	}

	tokens = wsRem

	return
}
