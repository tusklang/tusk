package main

import "fmt"
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

func NQReplace(s string) (string, error) {
    rs := make([]rune, 0, len(s))
    const out = rune(0)
    var quote rune = out
    var escape = false
    for _, r := range s {
        if !escape {
            if r == '`' || r == '"' || r == '\'' {
                if quote == out {
                    quote = r
                } else if quote == r {
                    quote = out
                }
            }
        }
        escape = !escape && r == '\\'
        if quote != out || !(r == ' ' || r == '\t') {
            rs = append(rs, r)
        }
    }
    if quote != out {
        err := fmt.Errorf("unmatched unescaped quote: %q", quote)
        return "", err
    }
    return string(rs), nil
}
