package main

import "os"

func main() {

  var dir = os.Args[1]
  var fileName = os.Args[2]
  var oatFileName = os.Args[3]

  file := read(dir + fileName, "File Not Found: " + dir + fileName, true)

  file = ReplaceNQ(file, " ", "")
  file = ReplaceNQ(file, "\t", "")

  var lex = lexer(file)

  var actions = actionizer(lex)

  writeOat(actions, dir, oatFileName)
}
