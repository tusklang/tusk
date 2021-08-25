package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/tusklang/tusk/grouper"
	"github.com/tusklang/tusk/operations"
	"github.com/tusklang/tusk/tokenizer"
)

func main() {

	cwd, _ := os.Getwd()
	wd := flag.String("wd", cwd, "Set the working directory of a Tusk program, defaults to the current working directory")

	flag.Parse()

	os.Chdir(*wd)

	//tmp
	a, _ := ioutil.ReadFile("./test.tusk")

	lex := tokenizer.Tokenizer(string(a))
	groups := grouper.Grouper(lex)
	ops, _ := operations.OperationsParser(groups)

	j, _ := json.MarshalIndent(ops, "", "  ")
	fmt.Println(string(j))
}
