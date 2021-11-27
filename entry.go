package main

import (
	"flag"

	"github.com/tusklang/tusk/compiler"
	"github.com/tusklang/tusk/parser"
)

func main() {
	config := flag.String("config", "tusk.config.json", "supply configuration file for tusk")
	flag.Parse()

	prog := parser.Initialize(*config)

	if prog == nil {
		return
	}

	compiler.Compile(prog, "test.ll")
}
