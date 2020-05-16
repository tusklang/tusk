package mango_get

import "net/http"
import "fmt"
import "encoding/gob"
import "io/ioutil"
import "os"

import "lang"

func save(args Args, dir string) {

  for _, v := range args.packages {
    res, err := http.Get(uri + v)

    if err != nil {
      fmt.Println("Warning: cannot perform mango get on", v)
      continue
    }

    defer res.Body.Close()
    file, err := ioutil.ReadAll(res.Body)

    if res.StatusCode != 200 {
      fmt.Println("Warning: cannot perform mango get on", v)
      continue
    }

    lex := lang.Lexer(string(file), dir, v)
    acts := lang.Actionizer(lex, false, dir, v)

    writefile, _ := os.Create(dir + "/mango/" + v)

    defer writefile.Close()

    encoder := gob.NewEncoder(writefile)
    encoder.Encode(acts)
  }
}
