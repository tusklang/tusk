package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/tusklang/tusk/parser"
)

func main() {

	cwd, _ := os.Getwd()
	wd := flag.String("wd", cwd, "Set the working directory of a Tusk program, defaults to the current working directory")

	flag.Parse()

	os.Chdir(*wd)

	fmt.Println(*wd)

	parser.Parse("./")
}
