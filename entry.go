package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/tusklang/tusk/ast"
	"github.com/tusklang/tusk/tokenizer"
)

func main() {

	//all of this is temporary

	a, _ := ioutil.ReadFile("./test.tusk")

	lex := tokenizer.Tokenizer(string(a))
	groups := ast.Grouper(lex)

	j, _ := json.MarshalIndent(groups, "", "  ")
	fmt.Println(string(j))
}
