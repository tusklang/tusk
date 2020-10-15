package compiler

//Item represents either a token or a brace enclosed group
type Item struct {

	//if type is a brace type, the Group field will be used
	//otherwise, the Token field will be used

	Type  string
	Name  string //for definite arrays
	File  string
	Line  uint64
	Group [][]Item
	Token Lex
}

var braceMatcher = map[string]string{
	"}":  "{",
	":]": "[:",
	")":  "(",
	"]":  "[",
}

func makeGroups(lex []Lex) ([][]Item, error) {

	var groups = [][]Item{[]Item{}}

	for i := 0; i < len(lex); i++ {

		if lex[i].Type == "?close_brace" {
			//error
			//unexpected closing brace
			return nil, makeCompilerErr("Unexpected closing brace", lex[i].Dir, lex[i].Line)
		}

		if lex[i].Type == "?open_brace" {

			var exp []Lex

			//brace types
			var braceTypes = map[string]int{
				"{":  0,
				"[:": 0,
				"(":  0,
				"[":  0,
			}
			/////////////

			var openBrace = lex[i]
			braceType := openBrace.Name

			braceTypes[braceType]++

			var line uint64 = 0

			for i++; i < len(lex); i++ {

				if line == 0 {
					line = lex[i].Line
				}

				//account for opening braces
				if _, exists := braceTypes[lex[i].Name]; exists {
					braceTypes[lex[i].Name]++
				}
				//account for closing braces
				if _, exists := braceTypes[braceMatcher[lex[i].Name]]; exists {
					braceTypes[braceMatcher[lex[i].Name]]--
				}

				if braceTypes["{"] == 0 && braceTypes["[:"] == 0 && braceTypes["("] == 0 && braceTypes["["] == 0 {
					break
				}

				exp = append(exp, lex[i])
			}

			if func() bool { //if any braces were not matched, return an error
				for _, v := range braceTypes {
					if v != 0 {
						return true
					}
				}

				return false
			}() {
				return [][]Item{}, makeCompilerErr("No brace found to match: "+braceType, openBrace.Dir, openBrace.Line)
			}

			groupedExp, e := makeGroups(exp)
			if e != nil {
				return groupedExp, e
			}

			var definitearray string

			if braceType == "[" { //definite array
				braceType = "("
				definitearray = "definite-array"
			}

			groups[len(groups)-1] = append(groups[len(groups)-1], Item{
				Type:  braceType,
				Name:  definitearray,
				File:  lex[i].Dir,
				Line:  line,
				Group: groupedExp,
			})
		} else {

			if lex[i].Name == "$term" || lex[i].Name == "," {
				groups = append(groups, []Item{})
				continue
			}

			groups[len(groups)-1] = append(groups[len(groups)-1], Item{
				Type:  lex[i].Type,
				File:  lex[i].Dir,
				Line:  lex[i].Line,
				Token: lex[i],
			})
		}

	}

	//filter empty groups out
	var filteredGroups = [][]Item{}

	for _, v := range groups {
		if len(v) != 0 {
			filteredGroups = append(filteredGroups, v)
		}
	}

	return filteredGroups /* <-- remove the last trailing group */, nil
}
