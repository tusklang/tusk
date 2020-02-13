package main

import "strconv"

func glossaryIndex(expStr_ []string, indexes [][]string, functions []Funcs, line uint64, calc_params paramCalcOpts, vars map[string]Variable, dir string) ([]string, [][]string) {
  formatter := func(gl []string) [][][]string {
    formatted := [][][]string{[][]string{[]string{}}}

    cbCnt := 0
    glCnt := 0
    bCnt := 0
    pCnt := 0

    expStr := gl[1:len(gl) - 1]

    for o := 0; o < len(expStr); o++ {
      if expStr[o] == "{" {
        cbCnt++;
      }
      if expStr[o] == "}" {
        cbCnt--;
      }

      if expStr[o] == "[:" {
        glCnt++;
      }
      if expStr[o] == ":]" {
        glCnt--;
      }

      if expStr[o] == "[" {
        bCnt++;
      }
      if expStr[o] == "]" {
        bCnt--;
      }

      if expStr[o] == "(" {
        pCnt++;
      }
      if expStr[o] == ")" {
        pCnt--;
      }

      if cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 0 && expStr[o] == "," {
        formatted = append(formatted, [][]string{[]string{}})
      } else if cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 0 && expStr[o] == ":" {
        formatted[len(formatted) - 1] = append(formatted[len(formatted) - 1], []string{})
      } else {
        formatted[len(formatted) - 1][len(formatted[len(formatted) - 1]) - 1] = append(formatted[len(formatted) - 1][len(formatted[len(formatted) - 1]) - 1], expStr[o])
      }
    }

    return formatted
  }

  cur := expStr_

  for ;len(indexes) > 0; {
    if cur[0] == "[:" {
      formatted := formatter(cur)

      for j := 0; j < len(formatted); j++ {

        left := mathParse([][]string{ formatted[j][0] }, functions, line, calc_params, vars, dir)
        index := mathParse([][]string{ indexes[0] }, functions, line, calc_params, vars, dir)

        if left[0][0] == index[0][1] {
          cur = formatted[j][1]
        }
      }

      indexes = indexes[1:]
    } else if cur[0] == "[" {
      cur, indexes = arrayIndex(cur, indexes, functions, line, calc_params, vars, dir)
    } else {
      cur, indexes = stringIndex(cur, indexes, functions, line, calc_params, vars, dir)
    }
  }

  return cur, indexes
}

func arrayIndex(expStr_ []string, indexes [][]string, functions []Funcs, line uint64, calc_params paramCalcOpts, vars map[string]Variable, dir string) ([]string, [][]string) {
  formatter := func(gl []string) [][]string {
    formatted := [][]string{[]string{}}

    cbCnt := 0
    glCnt := 0
    bCnt := 0
    pCnt := 0

    expStr := gl[1:len(gl) - 1]

    for o := 0; o < len(expStr); o++ {
      if expStr[o] == "{" {
        cbCnt++;
      }
      if expStr[o] == "}" {
        cbCnt--;
      }

      if expStr[o] == "[:" {
        glCnt++;
      }
      if expStr[o] == ":]" {
        glCnt--;
      }

      if expStr[o] == "[" {
        bCnt++;
      }
      if expStr[o] == "]" {
        bCnt--;
      }

      if expStr[o] == "(" {
        pCnt++;
      }
      if expStr[o] == ")" {
        pCnt--;
      }

      if cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 0 && expStr[o] == "," {
        formatted = append(formatted, []string{})
      } else {
        formatted[len(formatted) - 1] = append(formatted[len(formatted) - 1], expStr[o])
      }
    }

    return formatted
  }

  cur := expStr_

  for ;len(indexes) > 0; {

    if cur[0] == "[:" {
      cur, indexes = glossaryIndex(cur, indexes, functions, line, calc_params, vars, dir)
    } else if cur[0] == "[" {
      formatted := formatter(cur)

      index, _ := strconv.Atoi(mathParse(indexes, functions, line, calc_params, vars, dir)[0][1])

      cur = formatted[index]
      indexes = indexes[1:]
    } else {
      cur, indexes = stringIndex(cur, indexes, functions, line, calc_params, vars, dir)
    }
  }

  return cur, indexes
}

func stringIndex(expStr_ []string, indexes [][]string, functions []Funcs, line uint64, calc_params paramCalcOpts, vars map[string]Variable, dir string) ([]string, [][]string) {

  expStr := mathParse([][]string{ expStr_ }, functions, line, calc_params, vars, dir)[0][0]

  expStr = expStr[1:len(expStr) - 1]

  if len(indexes) > 1 {
    return []string{"undefined"}, indexes
  } else {

    index := mathParse([][]string{ indexes[0][1:len(indexes[0]) - 1] }, functions, line, calc_params, vars, dir)[0][0]

    index_, _ := strconv.Atoi(index)

    return []string{ string(expStr[index_]) }, indexes
  }
}
