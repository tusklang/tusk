package compiler

import (
	"encoding/json"
	"regexp"
	"strings"
	"unicode"
)

//Lex represents a single a token
type Lex struct {
	Name  string
	Line  uint64
	Type  string
	OName string
	Dir   string
}

//see tokens JSON
//this represents a token in tokensJSON
type compiletoken struct {
	Name    string `json:"name"`
	Remove  string `json:"remove"`
	Pattern string `json:"pattern"`
	Type    string `json:"type"`
}

var tokens []compiletoken
var _ = json.Unmarshal(tokensJSON, &tokens)

func testkey(file string, key compiletoken, lex []Lex) bool {
	r, _ := regexp.MatchString("^("+key.Pattern+")", file)

	//- and + can also be a positive or negative sign
	if r && (key.Name == "+" || key.Name == "-") && len(lex) != 0 && (lex[len(lex)-1].Type == "id" || lex[len(lex)-1].Type == "operation" || lex[len(lex)-1].Type == "?open_brace") {
		return false
	}

	return r
}

func lexer(file, filename string) ([]Lex, error) {

	var lex []Lex //lex to return
	var curline uint64 = 1

	for file != "" {
		if file[0] == '\n' { //if its a newline, add to the current line
			file = file[1:]
			curline++
			continue
		}

		if unicode.IsSpace(rune(file[0])) { //if its whitespace, ignore it
			file = file[1:]
			continue
		}

		//detect a comment
		//single line comments are written as ;comment
		//like in assembly and lisp
		if file[0] == ';' {
			var end = strings.Index(file, "\n")

			//if there is no newline after the comment, just break the loop
			if end == -1 {
				break
			}
			file = file[end:]

			continue
		}

		var cont bool

		for _, v := range tokens { //see if the current index is a key

			if testkey(file, v, lex) {
				file = strings.TrimPrefix(file, v.Remove)
				lex = append(lex, Lex{
					Name:  v.Name,
					Line:  curline,
					Type:  v.Type,
					OName: v.Name,
					Dir:   filename,
				})
				cont = true
				break
			}

		}

		if cont {
			continue
		}

		if file[0] == '"' || file[0] == '\'' || file[0] == '`' {
			var escaped bool
			qtype := file[0] //get the type of quote
			fullstr := ""

			for file = file[1:]; len(file) != 0; file = file[1:] {

				if escaped {
					switch file[0] {

					//acount for special characters (\n \r \t \v)
					case 'n':
						fullstr += "\n"
					case 'r':
						fullstr += "\r"
					case 't':
						fullstr += "\t"
					case 'v':
						fullstr += "\v"
					/////////////////////////////////////////////

					default:
						fullstr += string(file[0])
					}

					escaped = false
					continue
				}

				//backslash is escape
				if file[0] == '\\' {
					escaped = true
					continue
				}

				//if it is the same quote type, break
				if file[0] == qtype {
					break
				}

				fullstr += string(file[0])
			}

			fullstr = string(qtype) + fullstr + string(qtype) //surround it with the quotes
			file = file[1:]                                   //remove the ending quote from the file
			lex = append(lex, Lex{
				Name:  fullstr,
				Line:  curline,
				Type:  "expression value",
				OName: fullstr,
				Dir:   filename,
			})
			continue
		} else if unicode.IsDigit(rune(file[0])) || file[0] == '+' || file[0] == '-' || file[0] == '.' {
			var positive = true

			if file[0] == '-' {
				positive = false
				file = file[1:]
			} else if file[0] == '+' {
				file = file[1:]
			}

			numv := ""

			for len(file) != 0 && (unicode.IsDigit(rune(file[0])) || file[0] == '.') {
				numv += string(file[0])
				file = file[1:]
			}

			if !positive {
				numv = "-" + numv
			}

			lex = append(lex, Lex{
				Name:  numv,
				Line:  curline,
				Type:  "expression value",
				OName: numv,
				Dir:   filename,
			})
			continue
		} else {

			if file[0] == '\n' { //ends in newline
				file = file[1:]
				curline++
				continue
			}

			if unicode.IsSpace(rune(file[0])) { //it its a space (tab too)
				file = file[1:]
				continue
			}

			variable := ""

			for ; len(file) != 0; file = file[1:] {

				if file[0] == '\n' { //ends in newline, break
					file = file[1:]
					curline++
					break
				}

				if unicode.IsSpace(rune(file[0])) { //it its a space, stop
					file = file[1:]
					break
				}

				for _, v := range tokens {

					//only count operations
					if (v.Type == "operation" || v.Type == "?operation" || v.Type == "?open_brace" || v.Type == "?close_brace") && testkey(file, v, lex) || file[0] == ';' /* it is a comment */ {
						goto break_var_loop
					}

				}

				variable += string(file[0])
			}
		break_var_loop:

			//only if the variable is not nothing
			if variable != "" {
				lex = append(lex, Lex{
					Name:  "$" + variable,
					Line:  curline,
					Type:  "expression value",
					OName: variable,
					Dir:   filename,
				})
			}
		}

	}

	lex = term_inserter(tilde_inserter(insertArrows(lex)))

	return lex, nil
}
