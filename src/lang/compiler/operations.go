package compiler

type Operation struct {
  Type        string
  Line        uint64
  Left       *Operation
  Right      *Operation
  Degree     *Operation
  Item        Item //in case there is no operation
}

var operations []map[string]func(exp []Item, index int, opType string)Operation

func operationIncludes(group []Item) bool {
  for _, v := range operations {
    for k := range v {
      for _, i := range group {
        if i.Token.Name == k {
          return true
        }
      }
    }
  }

  return false
}

func indexOper(group []Item, oper string) int {

  if oper == "~" || oper == ":" { //because these ones start from the beginning
    for k, v := range group {
      if v.Token.Name == oper {
        return k
      }
    }

    return -1
  }

  for i := len(group) - 1; i >= 0; i-- {
    if group[i].Token.Name == oper {
      return i
    }
  }

  return -1
}

//operation function used for most operators (except assignment, not gate, similarity, function calls, etc...)
func normalOpFunc(exp []Item, index int, opType string) Operation {

  var (
    left = exp[:index]
    right = exp[index + 1:]
  )

  return Operation{
    Type: opType,
    Line: exp[index].Line,
    Left: &makeOperations([][]Item{ left })[0],
    Right: &makeOperations([][]Item{ right })[0],
  }
}

//operation function for both similarity operators (~~ and ~~~)
func similarityOpFunc(exp []Item, index int, opType string) Operation {

  hasColon := false //if it has a colon (to indicate a degree)
  var degExp []Item //expression of the degree
  var i int

  for i = index + 1; i < len(exp); i++ {
    if exp[i].Token.Name == ":" {
      hasColon = true
      break
    }
    degExp = append(degExp, exp[i])
  }

  if !hasColon {
    return normalOpFunc(exp, index, opType)
  }

  left := exp[:index]
  right := exp[index + i + 1:]

  return Operation{
    Type: opType,
    Line: exp[index].Line,
    Left: &makeOperations([][]Item{ left })[0],
    Right: &makeOperations([][]Item{ right })[0],
    Degree: &makeOperations([][]Item{ degExp })[0],
  }
}

//ODO is
// statement operator (~)
// assigner (:)
// boolean operations (except not gate)
// comparisons
// exponentiation
// mult, div, modulo
// add, subtract
// index operations and function calls
// not gate
// assignment operations

func makeOperations(groups [][]Item) []Operation {

  operations = []map[string]func(exp []Item, index int, opType string) Operation {
    map[string]func(exp []Item, index int, opType string) Operation {
      "~": normalOpFunc,
      "?": normalOpFunc,
    },
    map[string]func(exp []Item, index int, opType string) Operation {
      ":": normalOpFunc,
    },
    map[string]func(exp []Item, index int, opType string) Operation {
      "&": normalOpFunc,
      "|": normalOpFunc,
    },
    map[string]func(exp []Item, index int, opType string) Operation {
      "=": normalOpFunc,
      "!=": normalOpFunc,
      ">": normalOpFunc,
      ">=": normalOpFunc,
      "<": normalOpFunc,
      "<=": normalOpFunc,
      "~~": similarityOpFunc,
      "~~~": similarityOpFunc,
    },
    map[string]func(exp []Item, index int, opType string) Operation {
      "^": normalOpFunc,
    },
    map[string]func(exp []Item, index int, opType string) Operation {
      "*": normalOpFunc,
      "/": normalOpFunc,
      "%": normalOpFunc,
    },
    map[string]func(exp []Item, index int, opType string) Operation {
      "+": normalOpFunc,
      "-": normalOpFunc,
    },
    map[string]func(exp []Item, index int, opType string) Operation {
      "::": normalOpFunc,
      "sync": normalOpFunc,
      "async": normalOpFunc,
    },
    map[string]func(exp []Item, index int, opType string) Operation {
      "!": func(exp []Item, index int, opType string) Operation {
        return Operation{
          Type: opType,
          Line: exp[index].Line,
          Right: &makeOperations([][]Item{ exp[index + 1:] })[0],
        }
      },
    },
    map[string]func(exp []Item, index int, opType string) Operation {
      "++": func(exp []Item, index int, opType string) Operation {
        return Operation{
          Type: opType,
          Line: exp[index].Line,
          Left: &makeOperations([][]Item{ exp[:index] })[0],
        }
      },
      "--": func(exp []Item, index int, opType string) Operation {
        return Operation{
          Type: opType,
          Line: exp[index].Line,
          Left: &makeOperations([][]Item{ exp[:index] })[0],
        }
      },
      "+=": normalOpFunc,
      "-=": normalOpFunc,
      "*=": normalOpFunc,
      "/=": normalOpFunc,
      "%=": normalOpFunc,
      "^=": normalOpFunc,
    },
  }

  var newGroups []Operation

  for _, v := range groups {

    if !operationIncludes(v) {
      newGroups = append(newGroups, Operation{
        Type: "none",
        Line: v[0].Line,
        Item: v[0],
      })
    }

    for _, val := range operations {
      for oper, function := range val {
        indexOfOper := indexOper(v, oper)
        if indexOfOper != -1 {
          newGroups = append(newGroups, function(v, indexOfOper, oper))
          goto breakOuter
        }
      }
    }

    breakOuter:
  }

  return newGroups
}
