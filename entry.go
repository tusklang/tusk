package main

import (
	"flag"
	"os"

	"github.com/tusklang/tusk/parser"
)

func main() {
	flag.Parse()

	cwd, _ := os.Getwd()
	wd := flag.String("wd", cwd, "Set the working directory of a Tusk program, defaults to the current working directory")

	os.Chdir(*wd)

	parser.Parse("./")
}
