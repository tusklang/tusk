package main

import "os"

var operators = []string{"^", "*", "/", "%", "+", "-", "&", "|", "!", "~", ";"}

type paramCalcOpts struct {
  PREC             int
  LONG_MULT_THRESH int
}

func defaults(params *paramCalcOpts) {
  (*params).PREC = 1000
  (*params).LONG_MULT_THRESH = 5
}

func main() {
  var args = os.Args
  var dir = args[1]
  var fileName = args[2]
  params := paramCalcOpts{}

  defaults(&params)

  index(fileName, dir, params)
}
