package main

import "os"

var operators = []string{"^", "*", "/", "%", "+", "-", "&", "|", "!", "~", ";"}

type paramCalcOpts struct {}

func defaults(params paramCalcOpts) paramCalcOpts {

  return params
}

func main() {
  var args = os.Args
  var dir = args[1]
  var fileName = args[2]
  params := paramCalcOpts{}

  params = defaults(params)

  index(fileName, dir, params)
}
