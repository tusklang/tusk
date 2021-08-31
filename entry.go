package main

import (
	"flag"

	"github.com/tusklang/tusk/initialize"
)

func main() {

	config := flag.String("config", "tusk.config.json", "supply configuration file for tusk")

	flag.Parse()

	initialize.Initialize(*config)

	// //all of this is temporary

	// a, _ := ioutil.ReadFile("./test.tusk")

	// lex := tokenizer.Tokenizer(string(a))
	// ast, e := ast.GenerateAST(lex)

	// if e != nil {
	// 	log.Fatal(e)
	// }

	// if ev := validator.Validate(ast); ev != nil {
	// 	//types and variables are invalid somewhere
	// 	fmt.Println(ev)
	// }

	// f := initialize.Initialize()
	// _ = f

	// j, _ := json.MarshalIndent(ast, "", "  ")
	// fmt.Println(string(j))
}
