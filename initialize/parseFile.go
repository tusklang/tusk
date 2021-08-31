package initialize

import (
	"io/ioutil"
	"strings"

	"github.com/tusklang/tusk/ast"
	"github.com/tusklang/tusk/tokenizer"
	"github.com/tusklang/tusk/validator"
)

func parseFile(name string) (*File, error) {

	f, e := ioutil.ReadFile(name)

	if e != nil {
		return nil, e
	}

	lex := tokenizer.Tokenizer(string(f))
	ast, e := ast.GenerateAST(lex)

	if e != nil {
		return nil, e
	}

	if e = validator.Validate(ast); e != nil {
		return nil, e
	}

	retf := File{
		Name: strings.TrimSuffix(name, ".tusk"),
	}

	fetchGlobals(ast, &retf, &retf.Private)

	return &retf, nil
}
