package parser

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/tusklang/tusk/tokenizer"
)

//Parse returns the abstract syntax tree of a Tusk project
func Parse(dir string) {
	files, e := GetFiles(dir)

	_, _ = files, e

	absp, _ := filepath.Abs(dir)
	pkgname := filepath.Base(absp)

	_ = pkgname

	//read all the package's files into a single string
	var pkgfull string

	for _, v := range files {
		d, _ := ioutil.ReadFile(v) //read file
		pkgfull += string(d)       //append data to the full package
	}

	tokens := tokenizer.Tokenizer(pkgfull)
	groups := grouper(tokens)
	astv := genAST(groups)

	if asterr != nil {
		//error
	}

	j, _ := json.MarshalIndent(astv, "", "  ")
	fmt.Println(string(j))
}
