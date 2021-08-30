package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/tusklang/tusk/ast"
	"github.com/tusklang/tusk/initialize"
	"github.com/tusklang/tusk/tokenizer"
	"github.com/tusklang/tusk/validator"
)

func main() {

	//all of this is temporary

	a, _ := ioutil.ReadFile("./test.tusk")

	lex := tokenizer.Tokenizer(string(a))
	ast, e := ast.GenerateAST(lex)

	if e != nil {
		log.Fatal(e)
	}

	if ev := validator.Validate(ast); ev != nil {
		//types and variables are invalid somewhere
		fmt.Println(ev)
	}

	f := initialize.Initialize(ast)
	_ = f

	j, _ := json.MarshalIndent(f, "", "  ")
	fmt.Println(string(j))
}
