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
import "mangomm/cli/install"
import "mangomm/cli/remove"
import "mangomm/cli/wipe"
///////

////////////

import "lang/interpreter/bind"

func defaults(calc *bind.Calc, package_ *bind.Package, files *bind.Files, name string) {

  (*calc).SetPREC(10)

  if strings.LastIndex(name, ".") == -1 {
    (*calc).SetO(name + ".oat")
  } else {
    (*calc).SetO(name[:strings.LastIndex(name, ".")] + ".oat")
  }

  (*package_).SetADDON("lang")
  (*files).SetNAME("")
  (*files).SetDIR("C:/")
}

func main() {
  args := os.Args;

  var params = bind.NewCliParams()

  var (
    calc = bind.NewCalc()
    package_ = bind.NewPackage()
    files = bind.NewFiles()
  )

  if len(args) <= 2 {
    fmt.Println("Error, no input was given")
    os.Exit(1)
  }

  defaults(&calc, &package_, &files, args[2])

  files.SetDIR(args[1])
  files.SetNAME(args[2])

  for i := 2; i < len(args); i++ {

    v := args[i]

    if strings.HasPrefix(v, "--") {

      switch strings.ToUpper(v) {

        default:
          package_.SetADDON(v[2:])
      }

    } else if strings.HasPrefix(v, "-") {

      switch strings.ToUpper(v[1:]) {

        case "C":
          package_.SetADDON("compile")
        case "R":
          package_.SetADDON("run")
        case "PREC":
          temp_prec, _ := strconv.Atoi(args[i + 1])
          calc.SetPREC(temp_prec)
          i++
          i++
        case "O":
          calc.SetO(args[i + 1])
          i++
        default:
          fmt.Println("Caution, there is no cli parameter named", v)
          i++
      }
    }
  }

  params.SetCalc(calc)
  params.SetPackage(package_)
  params.SetFiles(files)

  switch strings.ToLower(package_.GetADDON()) {

    case "lang":
      lang.Run(params)
    case "compile":
      oatcompile.Compile(params)
    case "run":
      oatrun.Run(params)
    case "mango-get":
      mango_get.Get()
    case "mango-rm":
      mango_rm.Remove()
    case "mango-wipe":
      mango_wipe.Wipe()
    default:
      fmt.Println("Error: cannot use omm addon", package_.GetADDON())
  }
}
