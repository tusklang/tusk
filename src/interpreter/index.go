package main

import "os/exec"
import "os"
import "strings"

type Funcs struct {
  Name string
  Line uint64
}

func index(fileName, dir string, calcParams paramCalcOpts) {

  file := read("./pre.omm", "", true) + read(dir + fileName, "File Not Found: " + dir + fileName, true)
  fileNQ, _ := NQReplace(file)

  lexCmd := exec.Command("perl", "./lexer/main.pl")

  lexCmd.Stdin = strings.NewReader(fileNQ + "\n")
  lexCmd.Stdout = os.Stdout
  lexCmd.Stderr = os.Stderr

  err := lexCmd.Run()

  if err != nil {
    panic(lexCmd.Stderr)
  }

  /*var actions = actionizer(lex)

  parser(actions, calcParams, dir, 0, []Funcs{ Funcs{ file, 0 } }, make(map[string]Variable), false)*/
}
