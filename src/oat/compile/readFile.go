package main

import "fmt"
import "os"
import "bufio"

func read(fileName string, err string, newline bool) string {
  filePointer, error := os.Open(fileName)

  if error != nil {
    fmt.Println(err)
    os.Exit(1);
  }

  var scanner = bufio.NewScanner(filePointer)

  var file string

  for scanner.Scan() {

    if newline == false {
      file+=(scanner.Text())
    } else {
      file+=("\n" + scanner.Text())
    }
  }

  return file
}
