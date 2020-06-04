package lang

import "fmt"
import "os"
import "strings"

// #cgo CFLAGS: -std=c99
// #include "bind.h"
import "C"

//this is just a function to actionize the "this" process
func this_calc(i *int, lex []Lex, PARAM_COUNT uint, name, dir, filename string, id int) Action {

  var curLex = lex[*i]

  (*i)++

  args, indexes, subcaller, _ := callCalcParams(i, lex, len(lex), dir, filename)

  //if the required params are not equal to the arguments given
  if PARAM_COUNT != uint(len(args)) {

    //throw an error
    C.colorprint(C.CString("Error while actionizing in " + curLex.Dir + "!\n"), C.int(12))
    fmt.Println(" Expected", PARAM_COUNT, "argument(s), but got", len(args), "instead to call", /* say the process */ name, "\n\nError occured on line", curLex.Line, "\nFound near:", strings.TrimSpace(curLex.Exp))

    //exit the process
    os.Exit(1)
  }

  return Action{ name, "", []string{}, []Action{}, []string{}, args, []Condition{}, id, []Action{}, []Action{}, []Action{}, [][]Action{}, indexes, make(map[string][]Action), false, "private", subcaller, []int{}, []int{} }
}
