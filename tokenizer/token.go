package tokenizer

type Token struct {
	Name    string
	Type    string
	File    string
	Snippet string
	Row     int
	Col     int
}
