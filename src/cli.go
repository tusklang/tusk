package main

import "os"
import "strings"
import "fmt"
import "strconv"

//omm addons
import "lang" //omm language

import oatcompile "oat/compile" //compile omm to oat
import oatrun "oat/run" //run an oat file

//mango
import "mango/cli/install"
///////
////////////

func defaults(params *map[string]map[string]interface{}, name string) {

  (*params)["Calc"]["PREC"] = 1000
  (*params)["Calc"]["LONG_MULT_THRESH"] = 7

  if strings.LastIndex(name, ".") == -1 {
    (*params)["Calc"]["O"] = name + ".oat"
  } else {
    (*params)["Calc"]["O"] = name[:strings.LastIndex(name, ".")] + ".oat"
  }


  (*params)["Package"]["PACKAGE"] = "lang"
  (*params)["Files"]["NAME"] = ""
  (*params)["Files"]["DIR"] = "C:"
}

func main() {
  args := os.Args;

  var params = make(map[string]map[string]interface{})

  params["Files"] = make(map[string]interface{})
  params["Package"] = make(map[string]interface{})
  params["Calc"] = make(map[string]interface{})

  if len(args) <= 2 {
    fmt.Println("Error, no input was given")
    os.Exit(1)
  }

  defaults(&params, args[2])

  params["Files"]["DIR"] = args[1]
  params["Files"]["NAME"] = args[2]

  for i := 2; i < len(args); i++ {

    v := args[i]

    if strings.HasPrefix(v, "--") {

      switch strings.ToUpper(v) {

        default:
          params["Package"]["PACKAGE"] = v[2:]
      }

    } else if strings.HasPrefix(v, "-") {

      switch strings.ToUpper(v[1:]) {

        case "C":
          params["Package"]["PACKAGE"] = "compile"
        case "R":
          params["Package"]["PACKAGE"] = "run"
        case "PREC":
          params["Calc"]["PREC"], _ = strconv.Atoi(args[i + 1])
          i++
        case "LONG_MULT_THRESH":
          params["Calc"]["LONG_MULT_THRESH"], _ = strconv.Atoi(args[i + 1])
          i++
        case "O":
          params["Calc"]["O"] = args[i + 1]
          i++
        default:
          fmt.Println("Caution, there is no cli parameter named", v)
          i++
      }
    }
  }

  switch strings.ToLower(params["Package"]["PACKAGE"].(string)) {

    case "lang":
      lang.Run(params)
    case "compile":
      oatcompile.Compile(params)
    case "run":
      oatrun.Run(params)
    case "mango-get":
      mango_get.Get()
    case "mango-rm":
    case "mango-wipe":
    default:
      fmt.Println("Error: cannot use omm addon", params["Package"]["PACKAGE"])
  }
}
