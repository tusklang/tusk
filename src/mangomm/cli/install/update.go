package mango_get

import "net/http"
import "fmt"
import "encoding/json"
import "io/ioutil"

import "lang"

func includes(value []interface{}, includer interface{}) bool {

  for _, v := range value {
    if v == includer {
      return true
    }
  }

  return false
}

func update(args Args, dir string) {

  for _, v := range args.packages {

    res, err := http.Get(uri + v)

    if err != nil {
      fmt.Println("Warning: cannot perform mango get on", v)
      continue
    }
    if res.StatusCode != 200 {
      fmt.Println("Warning: cannot perform mango get on", v)
      continue
    }

    //read the json
    var gomap map[string]interface{}

    json_val := lang.ReadFileJS(dir + ".mngo")[0]["Content"]
    e := json.Unmarshal([]byte(json_val), &gomap)

    if e != nil {
      gomap = map[string]interface{}{}
    }

    if gomap["dependencies"] == nil {
      gomap["dependencies"] = []interface{}{}
    }

    if !includes(gomap["dependencies"].([]interface{}), v) {
      gomap["dependencies"] = append(gomap["dependencies"].([]interface{}), v)
    }

    encoded, _ := json.MarshalIndent(gomap, "", "  ")

    _ = ioutil.WriteFile(dir + ".mngo", encoded, 0644)
  }
}
