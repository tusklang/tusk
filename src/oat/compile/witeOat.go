package main

import "os"
import "encoding/json"
import "io/ioutil"
import "fmt"

func writeOat(actions []Action, dir, fileName string) {
  jsonOat, err := json.Marshal(actions)

  if err != nil {
    fmt.Println("There Was An Error While Writing Your File")
    os.Exit(1)
  }

  ioutil.WriteFile(dir + fileName, jsonOat, 0644)
}
