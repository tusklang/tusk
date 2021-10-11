package main

import (
	"flag"

	"github.com/tusklang/tusk/compiler"
	"github.com/tusklang/tusk/initialize"
)

func main() {
	config := flag.String("config", "tusk.config.json", "supply configuration file for tusk")
	flag.Parse()

	prog := initialize.Initialize(*config)
	compiler.Compile(prog, "test.ll")
}
