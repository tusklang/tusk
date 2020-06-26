package main

import "os"
import "strings"
import "fmt"
import "strconv"

//omm addons
import "lang/compiler" //omm language (compile into go slices and structs)

import "oat/compile" //compile omm to oat
import "oat/run" //run an oat file

//mango
import "mangomm/cli/install"
import "mangomm/cli/remove"
import "mangomm/cli/wipe"
///////

////////////

func defaults(cli_params *map[string]map[string]interface{}, name string) {

  (*cli_params)["Calc"]["PREC"] = 50

  if strings.LastIndex(name, ".") == -1 {
    (*cli_params)["Calc"]["O"] = name + ".oat"
  } else {
    (*cli_params)["Calc"]["O"] = name[:strings.LastIndex(name, ".")] + ".oat"
  }

  (*cli_params)["Package"]["ADDON"] = "lang"
  (*cli_params)["Files"]["NAME"] = ""
  (*cli_params)["Files"]["DIR"] = "C:/"
}

func main() {
  args := os.Args;

  var cli_params = make(map[string]map[string]interface{})

  cli_params["Calc"] = map[string]interface{}{}
  cli_params["Package"] = map[string]interface{}{}
  cli_params["Files"] = map[string]interface{}{}

  if len(args) <= 2 {
    fmt.Println("Error, no input was given")
    os.Exit(1)
  }

  defaults(&cli_params, args[2])

  cli_params["Files"]["DIR"] = args[1]
  cli_params["Files"]["NAME"] = args[2]

  for i := 2; i < len(args); i++ {

    v := args[i]

    if strings.HasPrefix(v, "--") {

      switch strings.ToUpper(v) {

        default:
          cli_params["Package"]["ADDON"] = v[2:]
      }

    } else if strings.HasPrefix(v, "-") {

      switch strings.ToUpper(v[1:]) {

        case "C":
          cli_params["Package"]["ADDON"] = "compile"
        case "R":
          cli_params["Package"]["ADDON"] = "run"
        case "PREC":
          temp_prec, _ := strconv.Atoi(args[i + 1])
          cli_params["Calc"]["PREC"] = temp_prec
          i++
          i++
        case "O":
          cli_params["Calc"]["O"] = args[i + 1]
          i++
        default:
          fmt.Println("Caution, there is no cli parameter named", v)
          i++
      }
    }
  }

  switch strings.ToLower(cli_params["Package"]["ADDON"].(string)) {

    case "lang":
      compiler.Run(cli_params)
    case "compile":
      oatCompile.Compile(cli_params)
    case "run":
      oatRun.Run(cli_params)
    case "mango-get":
      mango_get.Get()
    case "mango-rm":
      mango_rm.Remove()
    case "mango-wipe":
      mango_wipe.Wipe()
    default:
      fmt.Println("Error: cannot use omm addon", cli_params["Package"]["ADDON"].(string))
  }
}
