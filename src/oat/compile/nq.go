package main

import "strings"

func Chunk(val string, by int) []string {
  nVal := []string{}
  last := 0

  for i := 0; i < len(val); i++ {

    if i % by == 0 {
      nVal = append(nVal, val[last:i])
      last = i
    }
  }

  if (last != len(val)) {
    nVal = append(nVal, val[last:len(val)])
  }

  return nVal;
}

func NQindexOf(val string, sub string) int {
  var chunks = Chunk(val, len(sub));

  if len(val) < len(sub) {
    return -1
  }

  var quote = ""

  for i := 0; i < len(chunks); i++ {

    if strings.Contains(chunks[i], "'") {

      if chunks[i] == sub {
        return i * len(sub) - len(sub)
      }

      if quote == "" {

        quote = "'";
        continue;
      } else if quote == "'" {
        quote = "";
      }
    }

    if strings.Contains(chunks[i], "\"") {

      if chunks[i] == sub {
        return i * len(sub) - len(sub)
      }

      if quote == "" {

        quote = "\"";
        continue;
      } else if quote == "\"" {
        quote = "";
      }
    }

    if strings.Contains(chunks[i], "`") {

      if chunks[i] == sub {
        return i * len(sub) - len(sub)
      }

      if quote == "" {

        quote = "`";
        continue;
      } else if quote == "`" {
        quote = "";
      }
    }

    if quote != "" {
      continue;
    }

    if (chunks[i] == sub) {
      return i * len(sub) - len(sub)
    }
  }

  return -1
}

func getAllIndexes(exp string, phrase string) []uint64 {
  var indexes = []uint64{}

  typeOfQ := ""
  isEscaped := false

  for i := uint64(0); len(exp) > 0; i++ {

    if !isEscaped {
      switch (exp[0:1]) {
        case "'":
          if typeOfQ == "" {
            typeOfQ = "'"
          } else if typeOfQ == "'" {
            typeOfQ = ""
          }
        case "\"":
          if typeOfQ == "" {
            typeOfQ = "\""
          } else if typeOfQ == "\"" {
            typeOfQ = ""
          }
        case "`":
          if typeOfQ == "" {
            typeOfQ = "`"
          } else if typeOfQ == "`" {
            typeOfQ = ""
          }
      }
    } else {
      isEscaped = true
    }

    if typeOfQ == "" {
      if strings.HasPrefix(exp, phrase) {
        indexes = append(indexes, i)

        exp = exp[len(phrase):]
      } else {

        exp = exp[1:]
      }
    }
  }

  return indexes
}

func ReplaceNQ(exp string, phrase string, replace string) string {
  locations := getAllIndexes(exp, phrase)

  for i := 0; i < len(locations); i++ {
    exp_ := exp[0:locations[i]]
    _exp := exp[locations[i]:]

    _exp = strings.Replace(_exp, phrase, replace, 1)

    exp = exp_ + _exp
  }

  return exp
}
