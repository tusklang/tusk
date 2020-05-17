package mango_get

import "os"
import "strings"

const uri = "https://ommscript.herokuapp.com/mango/"

type Args struct {

  packages  []string

  isSaved     bool
}

//export Get
func Get() {

  var args = os.Args[1:]
  var dir = args[0]
  args = args[1:]

  var argStruct Args

  for _, v := range args {

    if (v[0] == '-' || strings.HasPrefix(v, "--")) {

      var substr = 2

      if v[0] == '-' {
        substr = 1
      }

      cmd := v[substr:]

      switch cmd {
        case "s": fallthrough
        case "save":
          argStruct.isSaved = true
      }
    } else {

      argStruct.packages = append(argStruct.packages, v)
    }
  }

  if argStruct.isSaved {
    save(argStruct, dir)
  } else {
    update(argStruct, dir)
  }
}
