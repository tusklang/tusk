package main

import "encoding/json"
import "strings"
import "regexp"

type Keyword struct {
	Name        string
	Remove      string
	Pattern     string
	Type        string
}

func lexer(file string) []string {

  var keywords_ = read("./keywords.json", "Error Keywords Cannot Be Read.\nHave You Modified The Source?", false)

  var keywordList []Keyword
  json.Unmarshal([]byte(keywords_), &keywordList)

  keywords := []string{}

  for ;len(file) > 0; {

		if file[0:1] == " " || file[0:1] == "	" {
			file = file[1:]
		}

    cont := false

    for i := 0; i < len(keywordList); i++ {

      if strings.HasPrefix(file, string(keywordList[i].Pattern)) {
        keywords = append(keywords, keywordList[i].Name)

        file = strings.TrimPrefix(file, keywordList[i].Remove)
        cont = true
        break
      }
    }

    if cont {
      continue
    }

    numMatch, _ := regexp.MatchString(`^(\d|\.)`, file)
    strMatch, _ := regexp.MatchString("^('|\"|`)", file)

    if numMatch {

      var num string

      for ;numMatch; {
        num+=file[0:1]

        file = file[1:]
        numMatch, _ = regexp.MatchString(`^(\d|\.)`, file)
      }

      keywords = append(keywords, num)
    } else if strMatch {

      var str string

			typeOfQ := ""
			isEscaped := false

			for i := 0; i < len(file); i++ {

				if (!isEscaped) {

					if []rune(file)[i] == '\\' {

						str+=string([]rune(file)[i])
						isEscaped = true
						continue;
					}

					if string([]rune(file)[i]) == "'" {
						if typeOfQ == "" {
							typeOfQ = "'"
						} else if typeOfQ == "'" {
							typeOfQ = "'"
						}
					}

					if string([]rune(file)[i]) == "\"" {
						if typeOfQ == "" {
							typeOfQ = "\""
						} else if typeOfQ == "\"" {
							typeOfQ = ""
						}
					}

					if string([]rune(file)[i]) == "`" {
						if typeOfQ == "" {
							typeOfQ = "`"
						} else if typeOfQ == "`" {
							typeOfQ = ""
						}
					}

					str+=string([]rune(file)[i])

					if typeOfQ == "" {
						break
					}
				} else {

					str+=string([]rune(file)[i])
					isEscaped = false
				}
			}

			file = file[len(str):]
			keywords = append(keywords, str)
    } else {
      varName := ""
      operators := []string{"+", "-", "*", "/", "~", "&", "|", "(", ")", "^", "%", "<", ">", "=", "!", "'", "\"", "`", "$", "#", ":", ";", ",", "[", "]", "{", "}"}

      for ;len(file) != 0 && !arrayContain(operators, file[:1]); {
        varName+=file[:1]
        file = file[1:]
      }

			keywords = append(keywords, varName)
    }
  }

  return keywords;
}
