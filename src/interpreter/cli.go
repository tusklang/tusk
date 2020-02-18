package main

import "os"
import "strconv"

var operators = []string{"^", "*", "/", "%", "+", "-", "&", "|", "!", "~", ";"}

type paramCalcOpts struct {
  logger      bool
  mult_thresh string
  precision   int
}

func defaults(params paramCalcOpts) paramCalcOpts {
  params.logger = false
  params.mult_thresh = "1000"
  params.precision = 100

  return params
}

func main() {
  var args = os.Args
  var dir = args[1]
  var fileName = args[2]
  params := paramCalcOpts{}

  params = defaults(params)

  if len(args) > 2 {
    var extraParams = args[3:]

    for i := 0; i < len(extraParams) - 1; i++ {

      switch extraParams[i] {
        case "-logger":
          params.logger = extraParams[i + 1] == "true"
        case "-mult_thresh":
          params.mult_thresh = extraParams[i + 1]
        case "-precision":
          params.precision, _ = strconv.Atoi(extraParams[i + 1])
      }
    }
  }

  index(fileName, dir, params)
}
