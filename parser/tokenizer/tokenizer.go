package tokenizer

func Tokenizer(data string) (tokens []Token) {

	c := 0 //current index in the file string

	var row, col = 1, 1 //current row and column in the file

	for c < len(data) {
		for _, v := range tokenlist {
			m, _ := v.FindStringMatch(data[c:])

			if m != nil && m.Group.Capture.Index == 0 {
				//the current token is matched
				matched := m.Group.Capture.String()

				tokens = append(tokens, Token{
					Name: matched,
					Row:  row,
					Col:  col,
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

	return
}
