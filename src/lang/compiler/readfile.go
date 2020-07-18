package compiler

import "io/ioutil"
import "path"
import "strings"
import "oat/helper"

import . "lang/types"

func includeSingle(filename string, line uint64, dir string) []Action {
  if strings.HasSuffix(filename, ".oat") {
    var decoded = oatHelper.FromOat(filename)
    return decoded.Actions
  } else if strings.HasSuffix(filename, ".omm") {
    filename = strings.TrimSuffix(filename, ".omm")
  }

  filename+=".omm"

  for _, v := range included {
    if v == filename {
      return []Action{}
    }
  }

  content, err := ioutil.ReadFile(filename)

  included = append(included, filename)

  if err != nil {
    compilerErr("Could not find file: " + filename, dir, line)
  }

  compiled, _ := Compile(string(content), filename)
  return compiled
}

func includer(filename string, line uint64, dir string) [][]Action {
  if strings.HasSuffix(filename, "*") {

    files, e := ioutil.ReadDir(strings.TrimSuffix(filename, "*"))

    if e != nil {
      compilerErr("Could not find directory: " + filename, dir, line)
    }

    var actions [][]Action

    for _, v := range files {

      if !strings.HasSuffix(v.Name(), ".omm") || !strings.HasSuffix(v.Name(), ".oat") {
        continue
      }

      if v.IsDir() {
        actions = append(actions, includer(path.Join(strings.TrimSuffix(filename, "*"), v.Name() + "/*"), line, dir)...)
      } else {
        actions = append(actions, includeSingle(path.Join(strings.TrimSuffix(filename, "*"), v.Name()), line, dir))
      }
    }

    return actions
  } else {
    return [][]Action{ includeSingle(filename, line, dir) }
  }
}
