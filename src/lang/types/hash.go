package types

import "strings"

type OmmHash struct {
  hash map[string]OmmType
  Length  uint64
}

func (hash OmmHash) At(idx string) OmmType {
  return hash.hash[idx]
}

func (hash *OmmHash) Set(idx string, val OmmType) {

  if hash.hash == nil {
    hash.hash = map[string]OmmType{}
  }

  if _, exists := hash.hash[idx]; !exists {
    hash.Length++
  }

  hash.hash[idx] = val
}

func (hash OmmHash) Exists(idx string) bool {
  _, exists := hash.hash[idx]
  return exists
}

func (hash OmmHash) Format() string {

  return func() string {

    if len(hash.hash) == 0 {
      return "[::]"
    }

    var formatted = "[:"

    for k, v := range hash.hash {

      vFormatted := v.Format()

      switch v.(type) {
        case OmmHash: //if it is another hash, add the indents
          if vFormatted != "[::]" {
            newlineSplit := strings.Split(vFormatted, "\n")

            vFormatted = ""

            for _, v := range newlineSplit {
              vFormatted+=strings.Repeat(" ", 2) + v + "\n"
            }

            vFormatted = strings.TrimSpace(vFormatted) //remove the trailing \n (because an extra was added) and the leading two spaces (because it will be on the same line)
          }
      }

      formatted+="\n" + strings.Repeat(" ", 2) + k + ": " + vFormatted + ","
    }

    return formatted + "\n:]"
  }() //staring with 2
}
