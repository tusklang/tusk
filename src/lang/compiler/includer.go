package compiler

import "io/ioutil"
import "path"
import "strings"
import "oat/helper"

import . "lang/types"

func includeSingle(filename string, line uint64, dir string) ([]Action, CompileErr) {
  if strings.HasSuffix(filename, ".oat") {
    var decoded = oatHelper.FromOat(filename)
    return decoded.Actions, nil
  }

  if strings.HasSuffix(filename, ".omm") {
    filename = strings.TrimSuffix(filename, ".omm")
  }

  filename+=".omm"

  for _, v := range included {
    if v == filename {
      return []Action{}, nil
    }
  }

  content, err := ioutil.ReadFile(filename)

  included = append(included, filename)

  if err != nil {
    return []Action{}, makeCompilerErr("Could not find file: " + filename, dir, line)
  }

  compiled, _, e := Compile(string(content), filename)

  if e != nil {
    return []Action{}, e
  }

  return compiled, nil
}

func includer(filename string, line uint64, dir string) ([][]Action, CompileErr) {
  if strings.HasSuffix(filename, "*") {

    files, e := ioutil.ReadDir(strings.TrimSuffix(filename, "*"))

    if e != nil {
      return [][]Action{}, makeCompilerErr("Could not find directory: " + filename, dir, line)
    }

    var actions [][]Action

    for _, v := range files {

      if !strings.HasSuffix(v.Name(), ".omm") || !strings.HasSuffix(v.Name(), ".oat") {
        continue
      }

      if v.IsDir() {
        inc, e := includer(path.Join(strings.TrimSuffix(filename, "*"), v.Name() + "/*"), line, dir)

        if e != nil {
          return [][]Action{}, e
        }

        actions = append(actions, inc...)
      } else {
        inc, e := includeSingle(path.Join(strings.TrimSuffix(filename, "*"), v.Name()), line, dir)

        if e != nil {
          return [][]Action{}, e
        }

        actions = append(actions, inc)
      }
    }

    return actions, nil
  } else {
    inc, e := includeSingle(filename, line, dir)

    if e != nil {
      return [][]Action{}, e
    }

    return [][]Action{ inc }, nil
  }
}
