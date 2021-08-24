package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/tusklang/tusk/grouper"
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

	j, _ := json.MarshalIndent(groups, "", "  ")
	fmt.Println(string(j))
}
