package parser

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/tusklang/tusk/ast"
	"github.com/tusklang/tusk/errhandle"
	"github.com/tusklang/tusk/tokenizer"
)

func parseFile(name string) (*File, *errhandle.TuskError) {

	f, readerr := ioutil.ReadFile(name)

	if readerr != nil {
		return nil, nil
	}

	lex := tokenizer.Tokenizer(string(f), name)
	generatedAST, e := ast.GenerateAST(lex)

	if e != nil {
		return nil, e
	}

	retf := File{
		Name: strings.TrimSuffix(filepath.Base(name), ".tusk"), //get the classname of the file
	}

	ferr := fetchGlobals(generatedAST, &retf, 2, 0)

	return &retf, ferr
}
