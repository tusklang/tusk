package mango_rm

import "os"
import "encoding/json"
import "io/ioutil"

import "lang/compiler"

func includes(value []interface{}, includer interface{}) bool {

  for _, v := range value {
    if v == includer {
      return true
    }
  }

  return false
}

func indexof(value []interface{}, sub interface{}) uint64 {

  var i uint64 = 0

  for _, v := range value {

    if v == sub {
      return i
    }

    i++
  }

  return 0
}

func removePackage(value []interface{}, sub interface{}) []interface{} {

  index := indexof(value, sub)

  //remove the index of the value
  return append(value[:index], value[index + 1:]...)
}

//export Remove
func Remove() {

  var args = os.Args[1:]
  var dir = args[0]
  args = args[1:]

  for _, v := range args {

    var gomap map[string]interface{}

    json_val := compiler.ReadFileJS(dir + ".mngo")[0]["Content"]
    e := json.Unmarshal([]byte(json_val), &gomap)

    if e != nil {
      continue
    }

    if gomap["dependencies"] == nil {
      continue
    }

    if !includes(gomap["dependencies"].([]interface{}), v) {
      continue
    }

    //remove the value from the deps
    deps := removePackage(gomap["dependencies"].([]interface{}), v)
    gomap["dependencies"] = deps

    encoded, _ := json.MarshalIndent(gomap, "", "  ")

    _ = ioutil.WriteFile(dir + ".mngo", encoded, 0644)

  }
}
